package provisioner

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"sync"
)

// Provisioner gets kube profile and returns list of task ids of provision masterTasks
type Provisioner interface {
	Provision(context.Context, *profile.KubeProfile, *steps.Config) ([]*workflows.Task, error)
}

type TaskProvisioner struct {
	repository storage.Interface
}

var provisionMap map[clouds.Name][]string

func init() {
	provisionMap = map[clouds.Name][]string{
		clouds.DigitalOcean: {workflows.DigitalOceanMaster, workflows.DigitalOceanNode},
	}
}

func NewProvisioner(repository storage.Interface) *TaskProvisioner {
	return &TaskProvisioner{
		repository: repository,
	}
}

// prepare creates all tasks for provisioning according to cloud provider
func (r *TaskProvisioner) prepare(name clouds.Name, masterCount, nodeCount int) ([]*workflows.Task, []*workflows.Task) {
	masterTasks := make([]*workflows.Task, 0, masterCount)
	nodeTasks := make([]*workflows.Task, 0, nodeCount)

	for i := 0; i < masterCount; i++ {
		t, _ := workflows.NewTask(provisionMap[name][0], r.repository)
		masterTasks = append(masterTasks, t)
	}

	for i := 0; i < nodeCount; i++ {
		t, _ := workflows.NewTask(provisionMap[name][1], r.repository)
		nodeTasks = append(nodeTasks, t)
	}

	return masterTasks, nodeTasks
}

// Provision runs provision process among nodes that have been provided for provision
func (r *TaskProvisioner) Provision(ctx context.Context, kubeProfile *profile.KubeProfile, config *steps.Config) ([]*workflows.Task, error) {
	masterTasks, nodeTasks := r.prepare(config.Provider, len(kubeProfile.MasterProfiles),
		len(kubeProfile.NodesProfiles))

	tasks := append(append(make([]*workflows.Task, 0), masterTasks...), nodeTasks...)

	go func() {
		config.Role = "master"
		config.ManifestConfig.IsMaster = true

		var wg sync.WaitGroup

		wg.Add(len(masterTasks))
		// Provision master nodes
		for _, masterTask := range masterTasks {
			go func() {
				result := masterTask.Run(ctx, config, os.Stdout)
				err := <-result

				if err != nil {
					logrus.Errorf("master task %s has finished with error %v", masterTask.ID, err)
				} else {
					logrus.Infof("master-task %s has finished", masterTask.ID)
				}
				wg.Done()
				}()
		}

		wg.Wait()
		// If we get no master node
		if config.GetMaster() == nil {
			logrus.Errorf("Cluster provisioning has failed")
			return
		}
		logrus.Info("Master provisioning has finished successfully")

		config.Role = "node"
		config.ManifestConfig.IsMaster = false
		// Let flannel communicate through private network
		config.FlannelConfig.EtcdHost = config.GetMaster().PrivateIp

		// Provision nodes
		for _, nodeTask := range nodeTasks {
			go func() {
				nodeTask.Run(ctx, config, os.Stdout)
			}()
		}
		logrus.Infof("Cluster %s has been deployed", config.ClusterName)
	}()

	return tasks, nil
}
