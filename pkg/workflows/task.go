package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func NewTask(providerRole string, config steps.Config, repository storage.Interface) (*Task, error) {
	switch providerRole {
	case DigitalOceanMaster:
		return newTask(DigitalOceanMaster, GetWorkflow(DigitalOceanMaster), config, repository), nil
	case DigitalOceanNode:
		return newTask(DigitalOceanNode, GetWorkflow(DigitalOceanNode), config, repository), nil
	default:
		w := GetWorkflow(providerRole)

		if w != nil {
			return newTask(providerRole, w, config, repository), nil
		}
	}

	return nil, ErrUnknownProviderWorkflowType
}

func newTask(workflowType string, workflow Workflow, config steps.Config, repository storage.Interface) *Task {
	id := uuid.New()

	return &Task{
		Id:     id,
		Config: config,
		Type:   workflowType,

		workflow:   workflow,
		repository: repository,
	}
}

// Run executes all steps of workflow and tracks the progress in persistent storage
func (w *Task) Run(ctx context.Context, out io.Writer) chan error {
	errChan := make(chan error)

	go func() {
		// Create list of statuses to track
		for _, step := range w.workflow {
			w.StepStatuses = append(w.StepStatuses, StepStatus{
				Status:   steps.StatusTodo,
				StepName: step.Name(),
				ErrMsg:   "",
			})
		}

		// Save task state before first step
		w.sync(ctx)
		// Start from the first step
		w.startFrom(ctx, w.Id, out, 0, errChan)
		logrus.Infof("Task %s has finished successfully", w.Id)
		close(errChan)
	}()

	return errChan
}

// Restart executes task from the last failed step
func (w *Task) Restart(ctx context.Context, id string, out io.Writer) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		data, err := w.repository.Get(ctx, prefix, id)

		if err != nil {
			errChan <- err
			return
		}

		err = json.Unmarshal(data, w)

		if err != nil {
			errChan <- err
			return
		}

		i := 0
		// Skip successfully finished steps
		for index, stepStatus := range w.StepStatuses {
			if stepStatus.Status == steps.StatusError {
				i = index
				break
			}
		}
		// Start from the last failed one
		w.startFrom(ctx, id, out, i, errChan)
	}()
	return errChan
}

// start task execution from particular step
func (w *Task) startFrom(ctx context.Context, id string, out io.Writer, i int, errChan chan error) {
	// Start workflow from the last failed step
	for index := i; index < len(w.StepStatuses); index++ {
		step := w.workflow[index]
		logrus.Println(step.Name())
		// Sync to storage with task in executing state
		w.StepStatuses[index].Status = steps.StatusExecuting
		w.sync(ctx)
		if err := step.Run(ctx, out, &w.Config); err != nil {
			// Mark step status as error
			w.StepStatuses[index].Status = steps.StatusError
			w.StepStatuses[index].ErrMsg = err.Error()
			w.sync(ctx)

			logrus.Error(err)
			errChan <- err
			return
		} else {
			// Mark step as success
			w.StepStatuses[index].Status = steps.StatusSuccess
			w.sync(ctx)
		}
	}
}

// synchronize state of workflow to storage
func (w *Task) sync(ctx context.Context) error {
	data, err := json.Marshal(w)
	buf := &bytes.Buffer{}

	if err != nil {
		return err
	}

	err = json.Indent(buf, data, "", "\t")

	if err != nil {
		return err
	}

	return w.repository.Put(ctx, prefix, w.Id, buf.Bytes())
}
