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
)

// Provisioner gets kube profile and returns list of task ids of provision masterTasks
type Provisioner interface {
	Provision(context.Context, *profile.KubeProfile, *steps.Config) ([]*workflows.Task, error)
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
func (r *TaskProvisioner) Provision(ctx context.Context, kubeProfile *profile.KubeProfile, config *steps.Config) ([]*workflows.Task, error) {
	masterTasks, nodeTasks := r.prepare(config.Provider, len(kubeProfile.MasterProfiles),
		len(kubeProfile.NodesProfiles))

	tasks := append(append(make([]*workflows.Task, 0), masterTasks...), nodeTasks...)

	go func() {
		config.IsMaster = true

		// TODO(stgleb): When we have concurrent provisioning use that to sync nodes and master provisioning
		// Provision master nodes
		for _, masterTask := range masterTasks {
			if masterTask == nil {
				continue
			}

			fileName := util.MakeFileName(masterTask.ID)
			out, err := r.getWriter(fileName)

			if err != nil {
				logrus.Errorf("Error getting writer for %s", fileName)
				return
			}

			result := masterTask.Run(ctx, *config, out)
			err = <-result

			if err != nil {
				logrus.Errorf("master task %s has finished with error %v", masterTask.ID, err)
			} else {
				logrus.Infof("master-task %s has finished", masterTask.ID)
			}
		}

		// TODO(stgleb): If master  provisioning has failed
		// on a step after build actual node handle this case
		// If we get no master node
		if config.GetMaster() == nil {
			logrus.Errorf("Cluster provisioning has failed, no master is up")
			return
		}
		logrus.Infof("Master provisioning for cluster %s has finished successfully", config.ClusterName)

		config.IsMaster = false
		config.ManifestConfig.IsMaster = false
		// Do internal communication inside private network
		config.FlannelConfig.EtcdHost = config.GetMaster().PrivateIp

		// Provision nodes
		for _, nodeTask := range nodeTasks {
			if nodeTask == nil {
				continue
			}

			fileName := util.MakeFileName(nodeTask.ID)
			out, err := r.getWriter(fileName)

			if err != nil {
				logrus.Errorf("Error getting writer for %s", fileName)
				return
			}
			result := nodeTask.Run(ctx, *config, out)
			err = <-result

			if err != nil {
				logrus.Errorf("node task %s has finished with error %v", nodeTask.ID, err)
			} else {
				logrus.Infof("node-task %s has finished", nodeTask.ID)
			}
		}
		logrus.Infof("Cluster %s has been deployed successfully", config.ClusterName)
	}()

	return tasks, nil
}
