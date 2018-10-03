package provisioner

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/pki"
)

const keySize = 4096

type KubeService interface {
	Create(ctx context.Context, k *model.Kube) error
	Get(ctx context.Context, name string) (*model.Kube, error)
}

type TaskProvisioner struct {
	kubeService  KubeService
	repository   storage.Interface
	getWriter    func(string) (io.WriteCloser, error)
	provisionMap map[clouds.Name]workflows.WorkflowSet
}

func NewProvisioner(repository storage.Interface, kubeService KubeService) *TaskProvisioner {
	return &TaskProvisioner{
		kubeService: kubeService,
		repository:  repository,
		provisionMap: map[clouds.Name]workflows.WorkflowSet{
			clouds.DigitalOcean: {
				ProvisionMaster: workflows.DigitalOceanMaster,
				ProvisionNode:   workflows.DigitalOceanNode,
			},
		},
		getWriter: util.GetWriter,
	}
}

// ProvisionCluster runs provisionCluster process among nodes that have been provided for provisionCluster
func (r *TaskProvisioner) ProvisionCluster(ctx context.Context, profile *profile.Profile, config *steps.Config) (map[string][]*workflows.Task, error) {
	masterTasks, nodeTasks, clusterTask := r.prepare(config.Provider, len(profile.MasterProfiles),
		len(profile.NodesProfiles))

	// TODO(stgleb): Make node names from task id before provisioning starts
	masters, nodes := nodesFromProfile(config.ClusterName, masterTasks, nodeTasks, profile)
	// Save cluster before provisioning
	r.buildInitialCluster(ctx, profile, masters, nodes, config)

	// monitor cluster state in separate goroutine
	go r.monitorClusterState(ctx, config)

	if err := bootstrapKeys(config); err != nil {
		return nil, errors.Wrap(err, "bootstrap keys")
	}

	go func() {
		// ProvisionCluster masters and wait until n/2 + 1 of masters with etcd are up and running
		doneChan, failChan, err := r.provisionMasters(ctx, profile, config, masterTasks)

		if err != nil {
			logrus.Errorf("ProvisionCluster master %v", err)
		}

		select {
		case <-ctx.Done():
			logrus.Errorf("Master cluster has not been created %v", ctx.Err())
			return
		case <-doneChan:
		case <-failChan:
			config.KubeStateChan() <- model.StateFailed
			logrus.Errorf("Master cluster deployment has been failed")
			return
		}

		// Save cluster state when masters are provisioned
		logrus.Infof("Master provisioning for cluster %s has finished successfully", config.ClusterName)

		// ProvisionCluster nodes
		r.provisionNodes(ctx, profile, config, nodeTasks)

		// Wait for cluster checks are finished
		r.waitCluster(ctx, clusterTask, config)
		logrus.Infof("Cluster %s deployment has finished", config.ClusterName)
	}()

	return map[string][]*workflows.Task{
		"master":  masterTasks,
		"node":    nodeTasks,
		"cluster": {clusterTask},
	}, nil
}

func (p *TaskProvisioner) ProvisionNodes(ctx context.Context, nodeProfiles []profile.NodeProfile, kube *model.Kube, config *steps.Config) ([]string, error) {
	if len(kube.Masters) != 0 {
		for key := range kube.Masters {
			config.AddMaster(kube.Masters[key])
		}
	} else {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "master node")
	}

	if err := bootstrapKeys(config); err != nil {
		return nil, errors.Wrap(err, "bootstrap keys")
	}

	providerWorkflowSet, ok := p.provisionMap[config.Provider]

	if !ok {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "provider workflow")
	}

	tasks := make([]string, 0, len(nodeProfiles))

	for _, nodeProfile := range nodeProfiles {
		// Take node workflow for the provider
		t, err := workflows.NewTask(providerWorkflowSet.ProvisionNode, p.repository)
		tasks = append(tasks, t.ID)

		if err != nil {
			return nil, errors.Wrap(sgerrors.ErrNotFound, "workflow")
		}

		writer, err := p.getWriter(t.ID)

		if err != nil {
			return nil, errors.Wrap(err, "get writer")
		}

		err = FillNodeCloudSpecificData(config.Provider, nodeProfile, config)

		if err != nil {
			return nil, errors.Wrap(err, "fill node profile data to config")
		}

		// Put task id to config so that create instance step can use this id when generate node name
		config.TaskID = t.ID
		errChan := t.Run(ctx, *config, writer)

		go func(cfg *steps.Config, errChan chan error) {
			err = <-errChan

			if err != nil {
				logrus.Errorf("add node to cluster %s caused an error %v", kube.Name, err)
				return
			}

			if n := cfg.GetNode(); n != nil {
				kube.Nodes[n.ID] = n
				// TODO(stgleb): Use some other method like update or Patch instead of recreate
				p.kubeService.Create(context.Background(), kube)
			} else {
				logrus.Errorf("Add node to cluster %s node was not added", kube.Name)
			}
		}(config, errChan)
	}

	return tasks, nil
}

// prepare creates all tasks for provisioning according to cloud provider
func (r *TaskProvisioner) prepare(name clouds.Name, masterCount, nodeCount int) ([]*workflows.Task, []*workflows.Task, *workflows.Task) {
	masterTasks := make([]*workflows.Task, 0, masterCount)
	nodeTasks := make([]*workflows.Task, 0, nodeCount)

	for i := 0; i < masterCount; i++ {
		t, err := workflows.NewTask(r.provisionMap[name].ProvisionMaster, r.repository)

		if err != nil {
			logrus.Errorf("Task type %s not found", r.provisionMap[name].ProvisionMaster)
			continue
		}
		masterTasks = append(masterTasks, t)
	}

	for i := 0; i < nodeCount; i++ {
		t, err := workflows.NewTask(r.provisionMap[name].ProvisionNode, r.repository)

		if err != nil {
			logrus.Errorf("Task type %s not found", r.provisionMap[name].ProvisionNode)
			continue
		}
		nodeTasks = append(nodeTasks, t)
	}

	clusterTask, _ := workflows.NewTask(workflows.Cluster, r.repository)

	return masterTasks, nodeTasks, clusterTask
}

func (p *TaskProvisioner) provisionMasters(ctx context.Context, profile *profile.Profile, config *steps.Config, tasks []*workflows.Task) (chan struct{}, chan struct{}, error) {
	config.IsMaster = true
	doneChan := make(chan struct{})
	failChan := make(chan struct{})

	if len(profile.MasterProfiles) == 0 {
		close(doneChan)
		return doneChan, failChan, nil
	}
	// master latch controls when the majority of masters with etcd are up and running
	// so etcd is available for writes of flannel that starts on each machine
	masterLatch := util.NewCountdownLatch(ctx, len(profile.MasterProfiles)/2+1)

	// If we fail n /2 of master deploy jobs - all cluster deployment is failed
	failLatch := util.NewCountdownLatch(ctx, len(profile.MasterProfiles)/2+1)

	// ProvisionCluster master nodes
	for index, masterTask := range tasks {
		if masterTask == nil {
			logrus.Fatal(tasks)
		}
		fileName := util.MakeFileName(masterTask.ID)
		out, err := p.getWriter(fileName)

		if err != nil {
			logrus.Errorf("Error getting writer for %s", fileName)
			return nil, nil, err
		}

		// Fulfill task config with data about provider specific node configuration
		p := profile.MasterProfiles[index]
		FillNodeCloudSpecificData(profile.Provider, p, config)

		go func(t *workflows.Task) {
			// Put task id to config so that create instance step can use this id when generate node name
			config.TaskID = t.ID
			result := t.Run(ctx, *config, out)
			err = <-result

			if err != nil {
				failLatch.CountDown()
				logrus.Errorf("master task %s has finished with error %v", t.ID, err)
			} else {
				masterLatch.CountDown()
				logrus.Infof("master-task %s has finished", t.ID)
			}
		}(masterTask)
	}

	go func() {
		masterLatch.Wait()
		close(doneChan)
	}()

	go func() {
		failLatch.Wait()
		close(failChan)
	}()

	return doneChan, failChan, nil
}

func (p *TaskProvisioner) provisionNodes(ctx context.Context, profile *profile.Profile, config *steps.Config, tasks []*workflows.Task) {
	config.IsMaster = false
	config.ManifestConfig.IsMaster = false
	// Do internal communication inside private network
	if master := config.GetMaster(); master != nil {
		config.FlannelConfig.EtcdHost = master.PrivateIp
	} else {
		return
	}

	// ProvisionCluster nodes
	for index, nodeTask := range tasks {
		fileName := util.MakeFileName(nodeTask.ID)
		out, err := p.getWriter(fileName)

		if err != nil {
			logrus.Errorf("Error getting writer for %s", fileName)
			return
		}

		// Fulfill task config with data about provider specific node configuration
		p := profile.NodesProfiles[index]
		FillNodeCloudSpecificData(profile.Provider, p, config)

		go func(t *workflows.Task) {
			// Put task id to config so that create instance step can use this id when generate node name
			config.TaskID = t.ID
			result := t.Run(ctx, *config, out)
			err = <-result

			if err != nil {
				logrus.Errorf("node task %s has finished with error %v", t.ID, err)
			} else {
				logrus.Infof("node-task %s has finished", t.ID)
			}
		}(nodeTask)
	}
}

func (p *TaskProvisioner) waitCluster(ctx context.Context, clusterTask *workflows.Task, config *steps.Config) {
	// clusterWg controls entire cluster deployment, waits until all final checks are done
	clusterWg := sync.WaitGroup{}
	clusterWg.Add(1)

	fileName := util.MakeFileName(clusterTask.ID)
	out, err := p.getWriter(fileName)

	if err != nil {
		logrus.Errorf("Error getting writer for %s", fileName)
		return
	}

	go func(t *workflows.Task) {
		defer clusterWg.Done()
		cfg := *config

		if master := config.GetMaster(); master != nil {
			cfg.Node = *master
		} else {
			config.KubeStateChan() <- model.StateFailed
			logrus.Errorf("No master found, cluster deployment failed")
			return
		}

		result := t.Run(ctx, cfg, out)
		err = <-result

		if err != nil {
			config.KubeStateChan() <- model.StateFailed
			logrus.Errorf("cluster task %s has finished with error %v", t.ID, err)
		} else {
			config.KubeStateChan() <- model.StateOperational
			logrus.Infof("cluster-task %s has finished", t.ID)
		}
	}(clusterTask)

	// Wait for all task to be finished
	clusterWg.Wait()
}

func (p *TaskProvisioner) buildInitialCluster(ctx context.Context, profile *profile.Profile, masters, nodes map[string]*node.Node, config *steps.Config) error {
	cluster := &model.Kube{
		State:        model.StateProvisioning,
		Name:         config.ClusterName,
		AccountName:  config.CloudAccountName,
		RBACEnabled:  profile.RBACEnabled,
		Region:       profile.Region,
		SshUser:      config.SshConfig.User,
		SshPublicKey: []byte(config.SshConfig.PublicKey),

		Auth: model.Auth{},

		Arch:                   profile.Arch,
		OperatingSystem:        profile.OperatingSystem,
		OperatingSystemVersion: profile.UbuntuVersion,
		K8SVersion:             profile.K8SVersion,
		DockerVersion:          profile.DockerVersion,
		HelmVersion:            profile.HelmVersion,
		Networking: model.Networking{
			Manager: profile.FlannelVersion,
			Version: profile.FlannelVersion,
			Type:    profile.NetworkType,
			CIDR:    profile.CIDR,
		},
		Masters: masters,
		Nodes:   nodes,
	}

	return p.kubeService.Create(ctx, cluster)
}

// Create bootstrap key pair and save to config ssh section
func bootstrapKeys(config *steps.Config) error {
	private, public, err := generateKeyPair(keySize)

	if err != nil {
		return err
	}

	config.SshConfig.BootstrapPrivateKey = private
	config.SshConfig.BootstrapPublicKey = public

	return nil
}

func bootstrapCA(config *steps.Config) error {
	pki.NewPKI()
}

// All cluster state changes during provisioning are made in this function
func (p *TaskProvisioner) monitorClusterState(ctx context.Context, cfg *steps.Config) {
	for {
		select {
		case n := <-cfg.NodeChan():
			k, err := p.kubeService.Get(ctx, cfg.ClusterName)

			if err != nil {
				logrus.Errorf("update kube state caused %v", err)
				continue
			}

			if n.Role == node.RoleMaster {
				k.Masters[n.Name] = &n
			} else {
				k.Nodes[n.Name] = &n
			}

			err = p.kubeService.Create(ctx, k)

			if err != nil {
				logrus.Errorf("update kube state caused %v", err)
				continue
			}
		case state := <-cfg.KubeStateChan():
			k, err := p.kubeService.Get(ctx, cfg.ClusterName)

			if err != nil {
				logrus.Errorf("update kube state caused %v", err)
				continue
			}

			k.State = state
			err = p.kubeService.Create(ctx, k)

			if err != nil {
				logrus.Errorf("update kube state caused %v", err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
