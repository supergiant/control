package provisioner

import (
	"testing"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestNodesFromProfile(t *testing.T) {
	region := "fra1"

	p := &profile.Profile{
		Provider: clouds.DigitalOcean,
		Region:   region,
		MasterProfiles: []profile.NodeProfile{
			{
				"image": "ubuntu-16-04-x64",
				"size":  "s-1vcpu-2gb",
			},
		},
		NodesProfiles: []profile.NodeProfile{
			{
				"image": "ubuntu-16-04-x64",
				"size":  "s-2vcpu-4gb",
			},
			{
				"image": "ubuntu-16-04-x64",
				"size":  "s-2vcpu-4gb",
			},
		},
	}

	cfg := &steps.Config{
		ClusterID: "test",
	}

	masterTasks, nodeTasks := []*workflows.Task{{ID: "1234"}}, []*workflows.Task{{ID: "5678"}, {ID: "4321"}}
	masters, nodes := nodesFromProfile(cfg.ClusterID, masterTasks, nodeTasks, p)

	if len(masters) != len(p.MasterProfiles) {
		t.Errorf("Wrong master node count expected %d actual %d",
			len(p.MasterProfiles), len(masters))
	}

	if len(nodes) != len(p.NodesProfiles) {
		t.Errorf("Wrong node count expected %d actual %d",
			len(p.NodesProfiles), len(nodes))
	}
}

func TestGrabTaskIds(t *testing.T) {
	clusterTsk := &workflows.Task{
		ID: "1234",
	}

	masterTasks := []*workflows.Task{
		{
			ID: "abcd",
		},
		{
			ID: "1sgsg",
		},
		{
			ID: "szrhhrs",
		},
	}

	preProvisionTask := &workflows.Task{}
	nodeTasks := []*workflows.Task{}

	taskMap := map[string][]*workflows.Task{
		workflows.ClusterTask:      {clusterTsk},
		workflows.MasterTask:       masterTasks,
		workflows.NodeTask:         nodeTasks,
		workflows.PreProvisionTask: {preProvisionTask},
	}

	taskIds := grabTaskIds(taskMap)

	if len(taskIds) != 4 {
		t.Errorf("Wrong task id count expected %d actual %d",
			len(masterTasks)+len(nodeTasks)+1, len(taskIds))
	}
}
