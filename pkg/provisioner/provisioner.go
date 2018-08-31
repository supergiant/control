package provisioner

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"sync"
)

// Provisioner gets kube profile and returns list of task ids of provision masterTasks
type Provisioner interface {
	Provision(context.Context, *profile.Profile, *steps.Config) ([]*workflows.Task, error)
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
func (r *TaskProvisioner) prepare(name clouds.Name, masterCount, nodeCount int) ([]*workflows.Task, []*workflows.Task) {
	masterTasks := make([]*workflows.Task, 0, masterCount)
	nodeTasks := make([]*workflows.Task, 0, nodeCount)

	for i := 0; i < masterCount; i++ {
		t, _ := workflows.NewTask(r.provisionMap[name][0], r.repository)
		masterTasks = append(masterTasks, t)
	}

	for i := 0; i < nodeCount; i++ {
		t, _ := workflows.NewTask(r.provisionMap[name][1], r.repository)
		nodeTasks = append(nodeTasks, t)
	}

	return masterTasks, nodeTasks
}

// Provision runs provision process among nodes that have been provided for provision
func (r *TaskProvisioner) Provision(ctx context.Context, kubeProfile *profile.Profile, config *steps.Config) ([]*workflows.Task, error) {
	masterTasks, nodeTasks := r.prepare(config.Provider, len(kubeProfile.MasterProfiles),
		len(kubeProfile.NodesProfiles))

	tasks := append(append(make([]*workflows.Task, 0), masterTasks...), nodeTasks...)

	go func() {
		config.IsMaster = true

		clusterWg := sync.WaitGroup{}
		masterWg := sync.WaitGroup{}

		clusterWg.Add(len(kubeProfile.MasterProfiles) + len(kubeProfile.NodesProfiles))
		masterWg.Add(len(kubeProfile.MasterProfiles))

		// Provision master nodes
		for index, masterTask := range masterTasks {
			fileName := util.MakeFileName(masterTask.ID)
			out, err := r.getWriter(fileName)

			if err != nil {
				logrus.Errorf("Error getting writer for %s", fileName)
				return
			}

			// Fulfill task config with data about provider specific node configuration
			p := kubeProfile.MasterProfiles[index]
			FillNodeCloudSpecificData(kubeProfile.Provider, p, config)

			go func(t *workflows.Task) {
				defer clusterWg.Done()
				defer masterWg.Done()

				result := t.Run(ctx, *config, out)
				err = <-result

				if err != nil {
					logrus.Errorf("master task %s has finished with error %v", t.ID, err)
				} else {
					logrus.Infof("master-task %s has finished", t.ID)
				}
			}(masterTask)
		}

		doneChan := make(chan struct{})
		go func() {
			config.Wait()
			close(doneChan)
		}()

		select {
		case <-ctx.Done():
			logrus.Errorf("Master cluster has not been created %v", ctx.Err())
			return
		case <-doneChan:
		}

		masterWg.Wait()
		logrus.Infof("Master provisioning for cluster %s has finished successfully", config.ClusterName)

		config.IsMaster = false
		config.ManifestConfig.IsMaster = false
		// Do internal communication inside private network
		config.FlannelConfig.EtcdHost = config.GetMaster().PrivateIp

		// Provision nodes
		for index, nodeTask := range nodeTasks {
			fileName := util.MakeFileName(nodeTask.ID)
			out, err := r.getWriter(fileName)

			if err != nil {
				logrus.Errorf("Error getting writer for %s", fileName)
				return
			}

			// Fulfill task config with data about provider specific node configuration
			p := kubeProfile.NodesProfiles[index]
			FillNodeCloudSpecificData(kubeProfile.Provider, p, config)

			go func(t *workflows.Task) {
				defer clusterWg.Done()

				result := t.Run(ctx, *config, out)
				err = <-result

				if err != nil {
					logrus.Errorf("node task %s has finished with error %v", t.ID, err)
				} else {
					logrus.Infof("node-task %s has finished", t.ID)
				}
			}(nodeTask)
		}

		// Wait for all task to be finished
		clusterWg.Wait()
		logrus.Infof("Cluster %s deployment has finished", config.ClusterName)
	}()

	return tasks, nil
}
