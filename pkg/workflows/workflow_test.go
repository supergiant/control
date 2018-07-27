package workflows

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
	"testing"
)

type fakeSynchronizer struct {
	storage map[string]string
}

func (f *fakeSynchronizer) Sync(ctx context.Context, key, data string) error {
	f.storage[key] = data

	return nil
}

type fakeStep struct {
	name        string
	description string
	err         error
}

func (f *fakeStep) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	out.Write([]byte(f.err.Error()))
	return f.err
}

func (f *fakeStep) Name() string {
	return f.name
}

func (f *fakeStep) Description() string {
	return f.description
}

func TestWorkFlowRun(t *testing.T) {
	workflow := WorkFlow{
		synchronizer: &fakeSynchronizer{
			storage: make(map[string]string),
		},
		workflowSteps: []steps.Step{
			&fakeStep{name: "step1", err: nil},
			&fakeStep{name: "step2", err: errors.New("shit happens")},
			&fakeStep{name: "step3", err: nil}},
	}

	buffer := &bytes.Buffer{}
	workflow.Run(context.Background(), buffer)
}

func TestWorkflowRestart(t *testing.T) {

}
