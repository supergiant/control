package provisioner

import (
	"context"
	"os"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

// Provisioner gets kube profile and returns list of task ids of provision masterTasks
type Provisioner interface {
	Prepare(int, int) []string
	Provision(ctx context.Context, nodes []node.Node)
	Cancel()
}

type TaskProvisioner struct {
	repository storage.Interface

	masterTasks []*workflows.Task
	nodeTasks   []*workflows.Task

	cancelFuncs []func()
}

func NewProvisioner(repository storage.Interface) *TaskProvisioner {
	return &TaskProvisioner{
		repository: repository,
	}
}

func (r *TaskProvisioner) Prepare(masterCount, nodeCount int) []string {
	tasksIds := make([]string, 0, nodeCount+masterCount)
	r.masterTasks = make([]*workflows.Task, nodeCount+masterCount)

	for i := 0; i < nodeCount; i++ {
		t, _ := workflows.NewTask(workflows.Nodetask, r.repository)
		r.nodeTasks = append(r.nodeTasks, t)
		tasksIds = append(tasksIds, t.ID)
	}

	for i := 0; i < masterCount; i++ {
		t, _ := workflows.NewTask(workflows.MasterTask, r.repository)
		r.masterTasks = append(r.masterTasks, t)
		tasksIds = append(tasksIds, t.ID)
	}

	return tasksIds
}

// Provision runs provision process among nodes that have been provided for provision
func (r *TaskProvisioner) Provision(ctx context.Context, nodes []node.Node) error {
	r.cancelFuncs = make([]func(), 0, len(nodes))

	i, j := 0, 0
	for _, n := range nodes {
		c, cancel := context.WithCancel(ctx)
		r.cancelFuncs = append(r.cancelFuncs, cancel)
		config := steps.Config{
			Node: n,
		}

		if n.Role == workflows.MasterTask {
			t := r.masterTasks[i]
			t.Run(c, config, os.Stdout)

			i += 1
		} else {
			t := r.masterTasks[j]
			t.Run(c, config, os.Stdout)

			j += 1
		}
	}

	return nil
}

// Cancel call cancel functions of all context of all masterTasks
func (r *TaskProvisioner) Cancel() {
	for _, cancel := range r.cancelFuncs {
		cancel()
	}
}
