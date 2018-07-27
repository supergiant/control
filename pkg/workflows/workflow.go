package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
)

type StepStatus struct {
	Status   steps.Status `json:"status"`
	StepName string       `json:"step_name"`
	ErrMsg   string       `json:"error_message"`
}

type WorkFlow struct {
	Type         string       `json:"type"`
	Config       steps.Config `json:"config"`
	StepStatuses []StepStatus `json:"steps"`

	workflowSteps []steps.Step
	synchronizer  Synchronizer
}

const (
	MasterWorkFlow = "master"
	NodeWorkflow   = "node"
)

var (
	masterSteps = []steps.Step{}
	nodeSteps   = []steps.Step{}

	ErrUnknownWorkflowType = errors.New("unknown workflow type")
)

type Synchronizer interface {
	Sync(context.Context, string, string) error
}

func New(workflowType string, config steps.Config, syncer Synchronizer) (*WorkFlow, error) {
	if workflowType == MasterWorkFlow {
		return &WorkFlow{
			Config:        config,
			workflowSteps: masterSteps,
		}, nil
	} else if workflowType == NodeWorkflow {
		return &WorkFlow{
			Config:        config,
			workflowSteps: nodeSteps,
		}, nil
	}

	return nil, ErrUnknownWorkflowType
}

func (w *WorkFlow) Run(ctx context.Context, out io.Writer) (string, chan error) {
	errChan := make(chan error)
	v4, _ := uuid.NewV4()
	id := v4.String()

	go func() {
		// Create list of statuses to track
		for _, step := range w.workflowSteps {
			w.StepStatuses = append(w.StepStatuses, StepStatus{
				Status:   steps.StatusTodo,
				StepName: step.Name(),
				ErrMsg:   "",
			})
		}

		for index, step := range w.workflowSteps {
			if err := step.Run(ctx, out, w.Config); err != nil {
				// Mark step status as error
				w.StepStatuses[index].Status = steps.StatusError
				w.StepStatuses[index].ErrMsg = err.Error()
				w.sync(ctx, id)

				errChan <- err
			} else {
				// Mark step as success
				w.StepStatuses[index].Status = steps.StatusSuccess
				w.sync(ctx, id)
			}
		}

		close(errChan)
	}()

	return id, errChan
}

func Restart(ctx context.Context, id string) chan error {
	errChan := make(chan error)
	// TODO(stgleb): implement reading stuff about this particular workflow run and start from last failed step.
	return errChan
}

func (w *WorkFlow) sync(ctx context.Context, id string) error {
	data, err := json.Marshal(w)

	if err != nil {
		return err
	}

	return w.synchronizer.Sync(ctx, fmt.Sprintf("workflows/%s", id), string(data))
}
