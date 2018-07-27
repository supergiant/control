package workflows

import (
	"bytes"
	"context"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"testing"
)

type fakeSynchronizer struct {
	storage map[string]string
}

func (f *fakeSynchronizer) Sync(ctx context.Context, key, data string) error {
	f.storage[key] = data

	return nil
}

func TestWorkFlowRun(t *testing.T) {
	workflow := WorkFlow{
		synchronizer: &fakeSynchronizer{
			storage: make(map[string]string),
		},
		workflowSteps: make([]steps.Step, 4),
	}

	buffer := &bytes.Buffer{}
	workflow.Run(context.Background(), buffer)
}

func TestWorkflowRestart(t *testing.T) {

}
