package provisioner

import (
	"context"
	"os"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

// Provisioner gets kube profile and returns list of task ids of provision tasks
type Provisioner interface {
	Provision(ctx context.Context, nodes []node.Node)
	Cancel()
}

type TaskProvisioner struct {
	repository storage.Interface

	tasksIds    []string
	cancelFuncs []func()
}

// Provision runs provision process among nodes that have been provided for provision
func (r *TaskProvisioner) Provision(ctx context.Context, nodes []node.Node) ([]string, error) {
	r.cancelFuncs = make([]func(), 0, len(nodes))
	r.tasksIds = make([]string, 0, len(nodes))

	for _, n := range nodes {
		c, cancel := context.WithCancel(ctx)
		r.cancelFuncs = append(r.cancelFuncs, cancel)
		config := steps.Config{}
		t, err := workflows.NewTask(n.Role, r.repository)

		if err != nil {
			return nil, err
		}

		// TODO(stgleb): pass buffer here
		t.Run(c, config,  os.Stdout)
		r.tasksIds = append(r.tasksIds, t.Id)
	}

	return r.tasksIds, nil
}

// Cancel call cancel functions of all context of all tasks
func (r *TaskProvisioner) Cancel() {
	for _, cancel := range r.cancelFuncs {
		cancel()
	}
}
