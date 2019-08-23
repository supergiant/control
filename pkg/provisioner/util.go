package provisioner

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type RateLimiter struct {
	bucket *time.Ticker
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		bucket: time.NewTicker(interval),
	}
}

// Take either returns giving calling code ability to execute or blocks until
// bucket is full again
func (r *RateLimiter) Take() {
	<-r.bucket.C
}

// Fill cloud account specific data gets data from the map and puts to particular cloud provider config
func FillNodeCloudSpecificData(provider clouds.Name, nodeProfile profile.NodeProfile, config *steps.Config) error {
	if nodeProfile["isMaster"] != "" {
		config.IsMaster, _ = strconv.ParseBool(nodeProfile["isMaster"])
	}

	switch provider {
	case clouds.AWS:
		return util.BindParams(nodeProfile, &config.AWSConfig)
	case clouds.GCE:
		return util.BindParams(nodeProfile, &config.GCEConfig)
	case clouds.DigitalOcean:
		return util.BindParams(nodeProfile, &config.DigitalOceanConfig)
	case clouds.Packet:
		return util.BindParams(nodeProfile, &config.PacketConfig)
	case clouds.OpenStack:
		return util.BindParams(nodeProfile, &config.OpenStackConfig)
	case clouds.Azure:
		return util.BindParams(nodeProfile, &config.AzureConfig)
	default:
		return sgerrors.ErrUnknownProvider
	}

	return nil
}

func MergeConfig(source *steps.Config, destination *steps.Config) error {
	switch source.Provider {
	case clouds.AWS:
		data, err := json.Marshal(&source.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.GCE:
		data, err := json.Marshal(&source.GCEConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.GCEConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.DigitalOcean:
		data, err := json.Marshal(&source.DigitalOceanConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.DigitalOceanConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.Packet:
		return nil
	case clouds.OpenStack:
		data, err := json.Marshal(&source.OpenStackConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.OpenStackConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.Azure:
		data, err := json.Marshal(&source.AzureConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.AzureConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	default:
		return sgerrors.ErrUnknownProvider
	}

	// These items must be shared among configs
	destination.Kube.ExternalDNSName = source.Kube.ExternalDNSName
	destination.Kube.InternalDNSName = source.Kube.InternalDNSName
	destination.SetKubeStateChan(source.KubeStateChan())
	destination.SetNodeChan(source.NodeChan())
	destination.SetConfigChan(source.ConfigChan())
	destination.Masters = source.Masters
	destination.Nodes = source.Nodes
	destination.Kube.ID = source.Kube.ID
	destination.Provider = source.Provider
	destination.Kube.Name = source.Kube.Name
	destination.Kube.BootstrapToken = source.Kube.BootstrapToken
	destination.IsBootstrap = source.IsBootstrap
	destination.Kube.K8SVersion = source.Kube.K8SVersion

	return nil
}

func nodesFromProfile(clusterName string, masterTasks, nodeTasks []*workflows.Task, profile *profile.Profile) (map[string]*model.Machine, map[string]*model.Machine) {
	masters := make(map[string]*model.Machine)
	nodes := make(map[string]*model.Machine)

	for index, p := range profile.MasterProfiles {
		taskId := masterTasks[index].ID
		name := util.MakeNodeName(clusterName, taskId, true)

		// TODO(stgleb): check if we can lowercase node names for all nodes
		if profile.Provider == clouds.GCE {
			name = strings.ToLower(name)
		}
		n := &model.Machine{
			TaskID:   taskId,
			Name:     name,
			Provider: profile.Provider,
			Region:   profile.Region,
			State:    model.MachineStatePlanned,
		}

		util.BindParams(p, n)
		masters[n.Name] = n
	}

	for index, p := range profile.NodesProfiles {
		taskId := nodeTasks[index].ID
		name := util.MakeNodeName(clusterName, taskId[:4], false)

		// TODO(stgleb): check if we can lowercase node names for all nodes
		if profile.Provider == clouds.GCE {
			name = strings.ToLower(name)
		}
		n := &model.Machine{
			TaskID:   taskId,
			Name:     name,
			Provider: profile.Provider,
			Region:   profile.Region,
			State:    model.MachineStatePlanned,
		}

		util.BindParams(p, n)
		nodes[n.Name] = n
	}

	return masters, nodes
}

func grabTaskIds(taskMap map[string][]*workflows.Task) map[string][]string {
	taskIds := make(map[string][]string, 0)

	for taskSet := range taskMap {
		tasks := make([]string, 0, len(taskMap[taskSet]))

		for _, task := range taskMap[taskSet] {
			tasks = append(tasks, task.ID)
		}

		taskIds[taskSet] = tasks
	}

	return taskIds
}
