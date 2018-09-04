package provisioner

import (
	"context"
	"io"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

// Provisioner gets kube profile and returns list of task ids of provision masterTasks
type Provisioner interface {
	Provision(context.Context, *profile.Profile, *steps.Config) (map[string][]*workflows.Task, error)
}

type TaskProvisioner struct {
	repository   storage.Interface
	getWriter    func(string) (io.WriteCloser, error)
	provisionMap map[clouds.Name][]string
}

func NewProvisioner(repository storage.Interface) *TaskProvisioner {
	return &TaskProvisioner{
		repository: repository,
		provisionMap: map[clouds.Name][]string{
			clouds.DigitalOcean: {workflows.DigitalOceanMaster, workflows.DigitalOceanNode},
		},
		getWriter: func(name string) (io.WriteCloser, error) {
			// TODO(stgleb): Add log directory to params of supergiant
			return os.OpenFile(path.Join("/tmp", name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		},
	}
}

// prepare creates all tasks for provisioning according to cloud provider
func (r *TaskProvisioner) prepare(name clouds.Name, masterCount, nodeCount int) ([]*workflows.Task, []*workflows.Task, *workflows.Task) {
	masterTasks := make([]*workflows.Task, 0, masterCount)
	nodeTasks := make([]*workflows.Task, 0, nodeCount)

	for i := 0; i < masterCount; i++ {
		t, err := workflows.NewTask(r.provisionMap[name][0], r.repository)

		if err != nil {
			logrus.Errorf("Task type %s not found", r.provisionMap[name][0])
		}
		masterTasks = append(masterTasks, t)
	}

	for i := 0; i < nodeCount; i++ {
		t, err := workflows.NewTask(r.provisionMap[name][1], r.repository)

		if err != nil {
			logrus.Errorf("Task type %s not found", r.provisionMap[name][1])
		}
		nodeTasks = append(nodeTasks, t)
	}

	clusterTask, _ := workflows.NewTask(workflows.Cluster, r.repository)

	return masterTasks, nodeTasks, clusterTask
}

// Provision runs provision process among nodes that have been provided for provision
func (r *TaskProvisioner) Provision(ctx context.Context, profile *profile.Profile, config *steps.Config) (map[string][]*workflows.Task, error) {
	masterTasks, nodeTasks, clusterTask := r.prepare(config.Provider, len(profile.MasterProfiles),
		len(profile.NodesProfiles))

	go func() {
		// Provision masters and wait until n/2 + 1 of masters with etcd are up and running
		doneChan, failChan, err := r.provisionMasters(ctx, profile, config, masterTasks)

		if err != nil {
			logrus.Errorf("Provision master %v", err)
		}

		select {
		case <-ctx.Done():
			logrus.Errorf("Master cluster has not been created %v", ctx.Err())
			return
		case <-doneChan:
		case <-failChan:
			logrus.Errorf("Master cluster deployment has been failed")
			return
		}

		logrus.Infof("Master provisioning for cluster %s has finished successfully", config.ClusterName)

		// Provision nodes
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

func (p *TaskProvisioner) provisionMasters(ctx context.Context, profile *profile.Profile, config *steps.Config, tasks []*workflows.Task) (chan struct{}, chan struct{}, error) {
	config.IsMaster = true

	// master wait group controls when the majority of masters with etcd are up and running
	// so etcd is available for writes of flannel that starts on each machine
	masterWg := sync.WaitGroup{}
	masterWg.Add(len(profile.MasterProfiles)/2 + 1)

	// If we fail n /2 of master deploy jobs - all cluster deployment is failed
	failWg := sync.WaitGroup{}
	failWg.Add(len(profile.MasterProfiles)/2 + 1)

	// Provision master nodes
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
			result := t.Run(ctx, *config, out)
			err = <-result

			if err != nil {
				// Keep track of failed master deployment tasks
				failWg.Done()
				logrus.Errorf("master task %s has finished with error %v", t.ID, err)
			} else {
				masterWg.Done()
				logrus.Infof("master-task %s has finished", t.ID)
			}
		}(masterTask)
	}

	doneChan := make(chan struct{})
	failChan := make(chan struct{})

	go func() {
		masterWg.Wait()
		close(doneChan)
	}()

	go func() {
		failWg.Wait()
		close(failChan)
	}()

	return doneChan, failChan, nil
}

func (p *TaskProvisioner) provisionNodes(ctx context.Context, profile *profile.Profile, config *steps.Config, tasks []*workflows.Task) {
	config.IsMaster = false
	config.ManifestConfig.IsMaster = false
	// Do internal communication inside private network
	config.FlannelConfig.EtcdHost = config.GetMaster().PrivateIp

	// Provision nodes
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
	// Waitgroup controls entire cluster deployment, waits until all final checks are done
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
		// Run
		cfg := *config

		if master := cfg.GetMaster(); master != nil {
			cfg.Node = *master
		}

		result := t.Run(ctx, cfg, out)
		err = <-result

		if err != nil {
			logrus.Errorf("cluster task %s has finished with error %v", t.ID, err)
		} else {
			logrus.Infof("cluster-task %s has finished", t.ID)
		}
	}(clusterTask)

	// Wait for all task to be finished
	clusterWg.Wait()
}
