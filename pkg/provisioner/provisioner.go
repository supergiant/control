package provisioner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/pki"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner/dry"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/storage/memory"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/configmap"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/copycerts"
)

type KubeService interface {
	Create(ctx context.Context, k *model.Kube) error
	Get(ctx context.Context, name string) (*model.Kube, error)
}

type TaskProvisioner struct {
	kubeService KubeService
	repository  storage.Interface
	getWriter   func(string) (io.WriteCloser, error)
	// NOTE(stgleb): Since provisioner is shared object among all users of SG
	// this rate limiter will affect all users not allowing them to spin-up
	// to many instances at once, probably we may split rate limiter per user
	// in future to avoid interference between them.
	rateLimiter *RateLimiter

	// Cancel map - map of KubeID -> cancel function
	// that cancels
	cancelMap map[string]func()
}

func NewProvisioner(repository storage.Interface, kubeService KubeService,
	spawnInterval time.Duration, logDir string) *TaskProvisioner {
	return &TaskProvisioner{
		kubeService: kubeService,
		repository:  repository,
		getWriter:   util.GetWriterFunc(logDir),
		rateLimiter: NewRateLimiter(spawnInterval),
		cancelMap:   make(map[string]func()),
	}
}

type bufferCloser struct {
	io.Writer
	err error
}

func (b *bufferCloser) Close() error {
	return b.err
}

// ProvisionCluster runs provisionCluster process among nodes
// that have been provided for provisionCluster
func (tp *TaskProvisioner) ProvisionCluster(parentContext context.Context,
	clusterProfile *profile.Profile, config *steps.Config) (map[string][]*workflows.Task, error) {
	if err := util.BootstrapKeys(config); err != nil {
		return nil, errors.Wrap(err, "bootstrap keys")
	}

	if err := bootstrapCerts(config); err != nil {
		return nil, errors.Wrap(err, "bootstrap certs")
	}

	taskMap := tp.prepare(config, len(clusterProfile.MasterProfiles), len(clusterProfile.NodesProfiles))
	clusterTask := taskMap[workflows.ClusterTask][0]

	// Get clusterID from taskID
	if clusterTask != nil && len(clusterTask.ID) >= 8 {
		config.ClusterID = clusterTask.ID[:8]
	} else {
		return nil, errors.New(fmt.Sprintf("Wrong value of "+
			"cluster task %v", clusterTask))
	}

	// Save cancel that cancel cluster provisioning to cancelMap
	ctx, cancel := context.WithCancel(parentContext)
	tp.cancelMap[config.ClusterID] = cancel

	// TODO(stgleb): Make node names from task id before provisioning starts
	masters, nodes := nodesFromProfile(config.ClusterName,
		taskMap[workflows.MasterTask], taskMap[workflows.NodeTask],
		clusterProfile)

	// Gather all task ids
	taskIds := grabTaskIds(taskMap)
	// Save cluster before provisioning
	err := tp.buildInitialCluster(ctx, clusterProfile, masters, nodes,
		config, taskIds)

	if err != nil {
		return nil, errors.Wrap(err, "build initial cluster")
	}

	// monitor cluster state in separate goroutine
	go tp.monitorClusterState(ctx, config.ClusterID, config.NodeChan(),
		config.KubeStateChan(), config.ConfigChan())
	go tp.provision(ctx, taskMap, clusterProfile)
	// Move cluster to provisioning state
	config.KubeStateChan() <- model.StateProvisioning

	return taskMap, nil
}

func (tp *TaskProvisioner) ProvisionNodes(parentContext context.Context, nodeProfiles []profile.NodeProfile, kube *model.Kube, config *steps.Config) ([]string, error) {
	if len(kube.Masters) != 0 {
		for key := range kube.Masters {
			config.AddMaster(kube.Masters[key])
		}
	} else {
		return nil, errors.Wrap(sgerrors.ErrNotFound, "master node")
	}

	// Save cancel function that cancels node provisioning to cancelMap
	ctx, cancel := context.WithCancel(parentContext)
	tp.cancelMap[config.ClusterID] = cancel

	if err := tp.loadCloudSpecificData(ctx, config); err != nil {
		return nil, errors.Wrap(err, "load cloud specific config")
	}

	// monitor cluster state in separate goroutine
	go tp.monitorClusterState(ctx, config.ClusterID,
		config.NodeChan(), config.KubeStateChan(), config.ConfigChan())

	tasks := make([]string, 0, len(nodeProfiles))

	// TODO(stgleb): do this in async to avoid blocking the UI
	for _, nodeProfile := range nodeProfiles {
		// Protect cloud API with rate limiter
		tp.rateLimiter.Take()

		// Take node workflow for the provider
		t, err := workflows.NewTask(config, workflows.ProvisionNode, tp.repository)
		if err != nil {
			return nil, errors.Wrap(sgerrors.ErrNotFound, "workflow")
		}

		tasks = append(tasks, t.ID)

		fileName := util.MakeFileName(t.ID)
		writer, err := tp.getWriter(fileName)

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

		go func(task *workflows.Task, cfg *steps.Config, errChan chan error) {
			err = <-errChan

			if err != nil {
				task.Config.Node.State = model.MachineStateError
				task.Config.AddNode(&task.Config.Node)
				task.Config.NodeChan() <- task.Config.Node
				logrus.Errorf("add node to cluster %s caused an error %v", kube.ID, err)
				return
			}
		}(t, config, errChan)
	}

	return tasks, nil
}

func (tp *TaskProvisioner) Cancel(clusterID string) error {
	if cancelFunc := tp.cancelMap[clusterID]; cancelFunc != nil {
		cancelFunc()
	} else {
		return sgerrors.ErrNotFound
	}

	return nil
}

func (tp *TaskProvisioner) RestartClusterProvisioning(parentCtx context.Context,
	clusterProfile *profile.Profile,
	config *steps.Config, taskIdMap map[string][]string) error {

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Minute*30)
	tp.cancelMap[config.ClusterID] = cancel
	logrus.Debugf("Deserialize tasks")

	// Deserialize tasks and put them to map
	taskMap, err := tp.deserializeClusterTasks(ctx, config, taskIdMap)

	if err != nil {
		logrus.Errorf("Restart cluster provisioning %v", err)
		return errors.Wrapf(err, "Restart cluster provisioning")
	}

	go tp.monitorClusterState(ctx, config.ClusterID,
		config.NodeChan(), config.KubeStateChan(), config.ConfigChan())
	go tp.provision(ctx, taskMap, clusterProfile)

	if config != nil && config.Kube.State != model.StateOperational {
		// Restore provisioning state for failed cluster
		config.KubeStateChan() <- model.StateProvisioning
	}

	return nil
}

func (tp *TaskProvisioner) UpgradeCluster(parentCtx context.Context, nextVersion string, k *model.Kube,
	tasks map[string][]*workflows.Task, config *steps.Config) {
	bootstrapTask := tasks[workflows.MasterTask][0]
	masterTasks := tasks[workflows.MasterTask][1:]
	nodeTasks := tasks[workflows.NodeTask]

	go tp.monitorClusterState(parentCtx, k.ID, config.NodeChan(),
		config.KubeStateChan(), config.ConfigChan())

	// TODO(stgleb): uncomment this once UI handle Upgrading state of the cluster
	//config.KubeStateChan() <- model.StateUpgrading
	logrus.Infof("Upgrade from %s to %s", k.K8SVersion, nextVersion)
	bootstrapTask.Config.IsBootstrap = true
	fileName := util.MakeFileName(bootstrapTask.ID)
	writer, err := tp.getWriter(fileName)

	if err != nil {
		logrus.Errorf("error creating writer %v", err)
		return
	}

	logrus.Infof("upgrade bootstrap node %v", bootstrapTask.Config.Node)
	tp.upgradeMachine(bootstrapTask, writer)

	for i := 0; i < len(masterTasks); i++ {
		masterTask := masterTasks[i]
		fileName := util.MakeFileName(masterTask.ID)
		writer, err := tp.getWriter(fileName)

		if err != nil {
			logrus.Errorf("error creating writer %v", err)
			return
		}

		logrus.Infof("Upgrade master node %v", masterTask.Config.Node)
		go tp.upgradeMachine(masterTask, writer)
	}

	for i := 0; i < len(nodeTasks); i++ {
		nodeTask := nodeTasks[i]

		fileName := util.MakeFileName(nodeTask.ID)
		writer, err := tp.getWriter(fileName)

		if err != nil {
			logrus.Errorf("error creating writer %v", err)
			return
		}

		logrus.Infof("Upgrade worker node %v", nodeTask.Config.Node)
		tp.upgradeMachine(nodeTask, writer)
		config.KubeStateChan() <- model.StateOperational
		config.ConfigChan() <- config
	}
}

// provision do actual provisioning of master and worker nodes
func (tp *TaskProvisioner) provision(ctx context.Context,
	taskMap map[string][]*workflows.Task, clusterProfile *profile.Profile) {
	preProvisionTask := taskMap[workflows.PreProvisionTask]

	config := preProvisionTask[0].Config
	if preProvisionTask != nil && len(preProvisionTask) > 0 {
		logrus.Debugf("preprovision task %s",
			preProvisionTask[0].ID)

		if preProvisionErr := tp.preProvision(ctx, preProvisionTask[0], config); preProvisionErr != nil {
			config.KubeStateChan() <- model.StateFailed
			logrus.Errorf("Pre provisioning cluster %v", preProvisionErr)
			return
		}

		kubeChan, nodeChan, configChan := config.KubeStateChan(), config.NodeChan(), config.ConfigChan()
		config = preProvisionTask[0].Config
		config.SetKubeStateChan(kubeChan)
		config.SetNodeChan(nodeChan)
		config.SetConfigChan(configChan)
	}

	if len(taskMap[workflows.MasterTask]) == 0 {
		return
	}

	config.IsMaster = true
	config.IsBootstrap = true
	// Get bootstrap task as a first master task
	bootstrapTask := taskMap[workflows.MasterTask][0]
	taskMap[workflows.MasterTask] = taskMap[workflows.MasterTask][1:]

	logrus.Debug("Provision bootstrap node")
	err := tp.bootstrapMaster(ctx, clusterProfile, config, bootstrapTask)

	if err != nil {
		config.KubeStateChan() <- model.StateFailed
		logrus.Errorf("provisioning bootstrap node has been failed")
		return
	}

	config = bootstrapTask.Config
	config.IsBootstrap = false
	logrus.Debug("Provision masters")
	err = tp.provisionMasters(ctx, clusterProfile,
		config, taskMap[workflows.MasterTask])

	if err != nil {
		config.KubeStateChan() <- model.StateFailed
		logrus.Errorf("master cluster deployment has been failed")
		return
	}

	// Save cluster state when masters are provisioned
	logrus.Infof("master provisioning for cluster"+
		"%s has finished successfully",
		config.ClusterID)

	config.IsBootstrap = false
	config.IsMaster = false
	err = tp.provisionNodes(ctx, clusterProfile, config,
		taskMap[workflows.NodeTask])

	if err != nil {
		config.KubeStateChan() <- model.StateFailed
		logrus.Errorf("Node provision has failed with %v", err)
		return
	}

	if len(taskMap[workflows.ClusterTask]) != 0 {
		// Wait for cluster checks are finished
		clusterTask := taskMap[workflows.ClusterTask][0]
		err = tp.waitCluster(ctx, clusterTask, config)

		if err != nil {
			config.KubeStateChan() <- model.StateFailed
			logrus.Errorf("cluster task %s has finished with error %v", clusterTask.ID, err)
		} else {
			config.KubeStateChan() <- model.StateOperational
			logrus.Infof("cluster-task %s has finished", clusterTask.ID)
		}
	}
	logrus.Infof("cluster %s deployment has finished",
		config.ClusterID)
}

// prepare creates all tasks for provisioning according to cloud provider
func (tp *TaskProvisioner) prepare(config *steps.Config, masterCount, nodeCount int) map[string][]*workflows.Task {
	var (
		infraTask   *workflows.Task
		clusterTask *workflows.Task
		err         error
	)

	masterTasks := make([]*workflows.Task, 0, masterCount)
	nodeTasks := make([]*workflows.Task, 0, nodeCount)
	//some clouds (e.g. AWS) requires running tasks before provisioning nodes (creating a VPC, Subnets, SecGroups, etc)
	infraTask, err = workflows.NewTask(config, fmt.Sprintf("%s%s", config.Provider, workflows.Infra), tp.repository)
	if err != nil {
		// We can't go further without pre provision task
		logrus.Errorf("create pre provision task has finished with %v", err)
		return nil
	}

	infraTask.Config = config
	for i := 0; i < masterCount; i++ {
		t, err := workflows.NewTask(config, workflows.ProvisionMaster, tp.repository)
		if err != nil {
			logrus.Errorf("Failed to set up task for %s workflow", workflows.ProvisionMaster)
			continue
		}
		masterTasks = append(masterTasks, t)
	}

	for i := 0; i < nodeCount; i++ {
		t, err := workflows.NewTask(config, workflows.ProvisionNode, tp.repository)
		if err != nil {
			logrus.Errorf("Failed to set up task for %s workflow", workflows.ProvisionNode)
			continue
		}
		t.Config = config
		nodeTasks = append(nodeTasks, t)
	}

	clusterTask, err = workflows.NewTask(config, workflows.PostProvision, tp.repository)
	if err != nil {
		logrus.Errorf("Failed to set up task for %s workflow", workflows.PostProvision)
		return nil
	}

	taskMap := map[string][]*workflows.Task{
		workflows.MasterTask:  masterTasks,
		workflows.NodeTask:    nodeTasks,
		workflows.ClusterTask: {clusterTask},
	}

	if infraTask != nil {
		taskMap[workflows.PreProvisionTask] = []*workflows.Task{infraTask}
	}

	return taskMap
}

// preProvision is for preparing activities before instances can be creates like
// creation of VPC, key pairs, security groups, subnets etc.
func (tp *TaskProvisioner) preProvision(ctx context.Context, preProvisionTask *workflows.Task, config *steps.Config) error {
	fileName := util.MakeFileName(preProvisionTask.ID)
	out, err := tp.getWriter(fileName)

	if err != nil {
		logrus.Errorf("Error getting writer for %s", fileName)
		return err
	}

	result := preProvisionTask.Run(ctx, *config, out)
	err = <-result
	config.ConfigChan() <- preProvisionTask.Config

	if err != nil {
		logrus.Errorf("pre provision task %s has finished with error %v",
			preProvisionTask.ID, err)
		config.KubeStateChan() <- model.StateFailed
	}

	logrus.Infof("pre provision task %s has finished", preProvisionTask.ID)

	return err
}

func (tp *TaskProvisioner) bootstrapMaster(ctx context.Context,
	profile *profile.Profile, rootConfig *steps.Config,
	bootstrapTask *workflows.Task) error {

	fileName := util.MakeFileName(bootstrapTask.ID)
	out, err := tp.getWriter(fileName)

	if err != nil {
		logrus.Errorf("Error getting writer for %s", fileName)
		return errors.Wrapf(err, "Error getting writer for %s", fileName)
	}

	if err := MergeConfig(rootConfig, bootstrapTask.Config); err != nil {
		return errors.Wrapf(err, "merge pre provision config to bootstrap task config")
	}

	// Fulfill task config with data about provider specific node configuration
	p := profile.MasterProfiles[0]

	err = FillNodeCloudSpecificData(profile.Provider, p, bootstrapTask.Config)

	if err != nil {
		return errors.Wrap(err, "fill master profile data to config")
	}

	bootstrapTask.Config.TaskID = bootstrapTask.ID
	bootstrapTask.Config.IsBootstrap = true
	bootstrapTask.Config.IsMaster = true

	err = <-bootstrapTask.Run(ctx, *bootstrapTask.Config, out)
	rootConfig.ConfigChan() <- bootstrapTask.Config

	if err != nil {
		logrus.Errorf("bootstrap task %s has finished with error %v", bootstrapTask.ID, err)
		return errors.Wrapf(err, "master bootstrap task %s has finished with error %v", bootstrapTask.ID, err)
	} else {
		logrus.Infof("bootstrap %s has finished", bootstrapTask.ID)
	}

	logrus.Infof("bootstrap task %s has finished", bootstrapTask.ID)

	return nil
}

func (tp *TaskProvisioner) provisionMasters(ctx context.Context,
	profile *profile.Profile, rootConfig *steps.Config,
	tasks []*workflows.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error)

	for index := range tasks {
		task := tasks[index]

		// We need to wait only for not finished tasks
		if task.Status != statuses.Success {
			wg.Add(1)
		}
	}

	rootConfig.IsMaster = true
	// ProvisionCluster rest of master nodes master nodes
	for index, masterTask := range tasks {
		// When this is restart code - don't process task that succeed
		if masterTask.Status == statuses.Success || masterTask.Status == statuses.Executing {
			continue
		}

		// Take token that allows perform action with Cloud Provider API
		tp.rateLimiter.Take()

		fileName := util.MakeFileName(masterTask.ID)
		out, err := tp.getWriter(fileName)

		if err != nil {
			logrus.Errorf("Error getting writer for %s", fileName)
		}

		if err := MergeConfig(rootConfig, masterTask.Config); err != nil {
			return errors.Wrapf(err, "merge pre provision config to bootstrap task config")
		}

		// Fulfill task config with data about provider specific node configuration
		p := profile.MasterProfiles[index]

		err = FillNodeCloudSpecificData(profile.Provider, p, masterTask.Config)

		if err != nil {
			return errors.Wrap(err, "fill master profile data to config")
		}

		go func(t *workflows.Task) {
			defer wg.Done()

			// Put task id to config so that create instance step can use this id when generate node name
			t.Config.TaskID = t.ID
			t.Config.IsMaster = true
			t.Config.IsBootstrap = false

			result := t.Run(ctx, *t.Config, out)
			err = <-result
			errChan <- err

			if err != nil {
				logrus.Errorf("master task %s has finished with error %v", t.ID, err)
			} else {
				logrus.Infof("master-task %s has finished", t.ID)
			}
		}(masterTask)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return <-errChan
}

func (tp *TaskProvisioner) provisionNodes(ctx context.Context, profile *profile.Profile, rootConfig *steps.Config, tasks []*workflows.Task) error {
	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	// ProvisionCluster nodes
	for index, nodeTask := range tasks {
		// When this is restart code - don't process task that succeed
		if nodeTask.Status == statuses.Success || nodeTask.Status == statuses.Executing {
			continue
		}

		// Take token that allows perform action with Cloud Provider API
		tp.rateLimiter.Take()

		fileName := util.MakeFileName(nodeTask.ID)
		out, err := tp.getWriter(fileName)

		if err != nil {
			logrus.Errorf("Error getting writer for %s", fileName)
			return errors.Wrapf(err, "Error getting writer for %s", fileName)
		}

		// Fulfill task config with data about provider specific node configuration
		p := profile.NodesProfiles[index]
		if err := MergeConfig(rootConfig, nodeTask.Config); err != nil {
			logrus.Errorf("merge pre provision config to bootstrap task config caused %v", err)
		}

		if err := FillNodeCloudSpecificData(profile.Provider, p, nodeTask.Config); err != nil {
			return errors.Wrapf(err, "fill nodes profile caused")
		}

		// Put task id to config so that create instance step can use this id when generate node name
		nodeTask.Config.TaskID = nodeTask.ID

		go func(t *workflows.Task) {
			t.Config.IsMaster = false
			t.Config.IsBootstrap = false
			result := t.Run(ctx, *t.Config, out)
			err = <-result

			if err != nil {
				// Put node to error state
				t.Config.Node.State = model.MachineStateError
				t.Config.AddNode(&t.Config.Node)
				t.Config.NodeChan() <- t.Config.Node

				logrus.Errorf("node task %s has finished with error %v", t.ID, err)
			} else {
				logrus.Infof("node-task %s has finished", t.ID)
			}

			wg.Done()
		}(nodeTask)
	}

	wg.Wait()
	return nil
}

func (tp *TaskProvisioner) waitCluster(ctx context.Context, clusterTask *workflows.Task, config *steps.Config) error {
	fileName := util.MakeFileName(clusterTask.ID)
	out, err := tp.getWriter(fileName)

	if err != nil {
		return errors.Wrapf(err, "Error getting writer for %s", fileName)
	}

	// Build node provision script and save it to ConfigMap
	buildNodeProvisionScript(ctx, config)

	cfg := *config

	// Get master node to run cluster task
	if master := config.GetMaster(); master != nil {
		logrus.Printf("Change one master %v to another %v", *master, cfg.Node)
		cfg.Node = *master
	} else {
		return errors.New("No master found, cluster deployment failed")
	}
	clusterTask.Config = &cfg
	result := clusterTask.Run(ctx, *clusterTask.Config, out)
	err = <-result

	if err != nil {
		return errors.Wrapf(err, "cluster task %s has finished with error %v", clusterTask.ID, err)
	}
	clusterTask.Config.ConfigChan() <- clusterTask.Config

	return nil
}

func buildNodeProvisionScript(ctx context.Context, config *steps.Config) {
	// Prepare config for confimap step
	dryRunner := dry.NewDryRunner()
	repository := memory.NewInMemoryRepository()

	dryConfig := *config
	// These values must still be templated
	dryConfig.Node.PublicIp = "{{ .PublicIp }}"
	dryConfig.Node.PrivateIp = "{{ .PrivateIp }}"
	dryConfig.ExternalDNSName = fmt.Sprintf("{{ .%s }}", configmap.KubeAddr)
	dryConfig.Runner = dryRunner
	dryConfig.DryRun = true

	if config.Provider == clouds.AWS {
		// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-instance-addressing.html#using-instance-addressing-common
		dryConfig.Node.PublicIp = "$(curl http://169.254.169.254/latest/meta-data/public-ipv4)"
		dryConfig.Node.PrivateIp = "$(curl http://169.254.169.254/latest/meta-data/local-ipv4)"
	}

	task, err := workflows.NewTask(&dryConfig, workflows.ProvisionNode, repository)

	if err != nil {
		return
	}

	task.Config = &dryConfig
	resultChan := task.Run(ctx, dryConfig, &bufferCloser{
		Writer: &bytes.Buffer{},
	})

	if err := <-resultChan; err != nil {
		return
	}

	config.ConfigMap.Data = dryRunner.GetOutput()
}

func (tp *TaskProvisioner) buildInitialCluster(ctx context.Context,
	profile *profile.Profile, masters, nodes map[string]*model.Machine,
	config *steps.Config, taskIds map[string][]string) error {

	cluster := &model.Kube{
		ID:           config.ClusterID,
		State:        model.StateProvisioning,
		Name:         config.ClusterName,
		Provider:     profile.Provider,
		ProfileID:    profile.ID,
		AccountName:  config.CloudAccountName,
		RBACEnabled:  profile.RBACEnabled,
		ServicesCIDR: profile.K8SServicesCIDR,
		Region:       profile.Region,
		Zone:         profile.Zone,
		User:         profile.User,
		Password:     profile.Password,

		Auth: config.Kube.Auth,

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
		CloudSpec: profile.CloudSpecificSettings,
		Masters:   masters,
		Nodes:     nodes,
		Tasks:     taskIds,

		SSHConfig: config.Kube.SSHConfig,

		ExposedAddresses: config.Kube.ExposedAddresses,
		APIServerPort:    config.Kube.APIServerPort,
	}

	return tp.kubeService.Create(ctx, cluster)
}

func (t *TaskProvisioner) loadCloudSpecificData(ctx context.Context, config *steps.Config) error {
	k, err := t.kubeService.Get(ctx, config.ClusterID)

	if err != nil {
		logrus.Errorf("get kube caused %v", err)
		return err
	}

	return util.LoadCloudSpecificDataFromKube(k, config)
}

func bootstrapCerts(config *steps.Config) error {
	ca, err := pki.NewCAPair([]byte(config.Kube.Auth.ParentCert))
	if err != nil {
		return errors.Wrap(err, "bootstrap CA for provisioning")
	}
	config.Kube.Auth.CACert = string(ca.Cert)
	config.Kube.Auth.CAKey = string(ca.Key)
	config.Kube.Auth.CACertHash = ca.CertHash

	if config.KubeadmConfig.CertificateKey, err = copycerts.CreateCertificateKey(); err != nil {
		return errors.Wrap(err, "create certificate key")
	}

	admin, err := pki.NewAdminPair(ca)
	if err != nil {
		return errors.Wrap(err, "create admin certificates")
	}
	config.Kube.Auth.AdminCert = string(admin.Cert)
	config.Kube.Auth.AdminKey = string(admin.Key)

	return nil
}

// All cluster state changes during provisioning must be made in this function
func (tp *TaskProvisioner) monitorClusterState(ctx context.Context,
	clusterID string, nodeChan chan model.Machine, kubeStateChan chan model.KubeState,
	configChan chan *steps.Config) {
	for {
		select {
		case n := <-nodeChan:
			k, err := tp.kubeService.Get(ctx, clusterID)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}

			if n.Role == model.RoleMaster {
				k.Masters[n.Name] = &n
			} else {
				k.Nodes[n.Name] = &n
			}

			err = tp.kubeService.Create(ctx, k)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}
		case state := <-kubeStateChan:
			logrus.Debugf("monitor: get kube %s", clusterID)
			k, err := tp.kubeService.Get(ctx, clusterID)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}

			k.State = state
			logrus.Debugf("monitor: update kube %s with state %s",
				k.ID, state)
			err = tp.kubeService.Create(ctx, k)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}
		case config := <-configChan:
			logrus.Debugf("update kube %s with config", clusterID)
			k, err := tp.kubeService.Get(ctx, clusterID)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}

			logrus.Debugf("update kube %s with config", k.ID)
			util.UpdateKubeWithCloudSpecificData(k, config)

			err = tp.kubeService.Create(ctx, k)

			if err != nil {
				logrus.Errorf("cluster monitor: update kube state caused %v", err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (tp *TaskProvisioner) deserializeClusterTasks(ctx context.Context, kubeConfig *steps.Config, taskIdMap map[string][]string) (map[string][]*workflows.Task, error) {
	taskMap := make(map[string][]*workflows.Task)

	for taskSet, tasks := range taskIdMap {
		for _, taskId := range tasks {
			data, err := tp.repository.Get(ctx, workflows.Prefix, taskId)

			if err != nil {
				logrus.Debugf("error getting task %s %v", taskId, err)
				return nil, errors.Wrapf(err, "task id %s not found %b", taskId, err)
			}

			task, err := workflows.DeserializeTask(data, tp.repository)

			if err != nil {
				logrus.Debugf("error deserializing task %s %v", taskId, err)
				return nil, errors.Wrapf(err, "error deserializing task %s %v", taskId, err)
			}

			err = MergeConfig(kubeConfig, task.Config)

			if err != nil {
				return nil, errors.Wrapf(err, "deserialize task")
			}

			logrus.Infof("deserialize task id %s", task.ID)
			taskMap[taskSet] = append(taskMap[taskSet], task)
		}
	}

	return taskMap, nil
}

func (tp *TaskProvisioner) upgradeMachine(task *workflows.Task, writer io.WriteCloser) {
	task.Config.Node.State = model.MachineStateUpgrading
	task.Config.NodeChan() <- task.Config.Node

	resultChan := task.Run(context.Background(), *task.Config, writer)
	err := <-resultChan

	if err != nil {
		task.Config.Node.State = model.MachineStateError
		task.Config.NodeChan() <- task.Config.Node
		logrus.Errorf("task %s has finished with error %v", task.ID, err)
	}

	task.Config.Node.State = model.MachineStateActive
	task.Config.NodeChan() <- task.Config.Node
}
