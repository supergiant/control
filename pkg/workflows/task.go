package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

// Task is an entity that has it own state that can be tracked
// and written to persistent storage through repository, it executes
// particular workflow of steps.
type Task struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Config       *steps.Config `json:"config"`
	Status       steps.Status  `json:"status"`
	StepStatuses []StepStatus  `json:"stepsStatuses"`

	workflow   Workflow
	repository storage.Interface
}

func NewTask(taskType string, repository storage.Interface) (*Task, error) {
	w := GetWorkflow(taskType)

	if w == nil {
		return nil, ErrUnknownProviderWorkflowType
	}

	return newTask(taskType, w, repository), nil
}

func newTask(workflowType string, workflow Workflow, repository storage.Interface) *Task {
	return &Task{
		ID:     uuid.New(),
		Type:   workflowType,
		Status: steps.StatusTodo,

		workflow:   workflow,
		repository: repository,
	}
}

// Run executes all steps of workflow and tracks the progress in persistent storage
func (w *Task) Run(ctx context.Context, config steps.Config, out io.WriteCloser) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				w.Status = steps.StatusError
				if err := w.sync(ctx); err != nil {
					logrus.Errorf("sync error %v for task %s", err, w.ID)
				}
				errChan <- errors.Errorf("provisioning failed, unexpected panic: %v ", r)
			}
		}()
		if w == nil {
			return
		}

		// Create list of statuses to track
		for _, step := range w.workflow {
			w.StepStatuses = append(w.StepStatuses, StepStatus{
				Status:   steps.StatusTodo,
				StepName: step.Name(),
				ErrMsg:   "",
			})
		}

		// Set config to the task
		w.Config = &config
		// Save task state before first step
		if err := w.sync(ctx); err != nil {
			logrus.Errorf("Error saving task state %v", err)
		}
		// Start from the first step
		err := w.startFrom(ctx, w.ID, out, 0)

		if err != nil {
			errChan <- err
			return
		}

		logrus.Infof("Task %s has finished successfully", w.ID)
		// Notify provisioner that task output closed with error
		if err := out.Close(); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	return errChan
}

// Restart executes task from the last failed step
func (w *Task) Restart(ctx context.Context, id string, out io.Writer) chan error {
	errChan := make(chan error, 1)
	wsLog := util.GetLogger(out)

	wsLog.Infof("Restarting task %s", id)
	go func() {
		defer close(errChan)
		data, err := w.repository.Get(ctx, Prefix, id)

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
		err = w.startFrom(ctx, id, out, i)

		if err != nil {
			errChan <- err
		}
	}()
	return errChan
}

// start task execution from particular step
func (w *Task) startFrom(ctx context.Context, id string, out io.Writer, i int) error {
	// Start workflow from the last failed step
	wsLog := util.GetLogger(out)
	for index := i; index < len(w.StepStatuses); index++ {
		step := w.workflow[index]

		wsLog.Infof("[%s] - started", step.Name())
		logrus.Info(step.Name())

		// Sync to storage with task in executing state
		w.StepStatuses[index].Status = steps.StatusExecuting
		if err := w.sync(ctx); err != nil {
			logrus.Errorf("sync error %v", err)
		}

		if err := step.Run(ctx, out, w.Config); err != nil {
			// Mark step status as error
			w.StepStatuses[index].Status = steps.StatusError
			w.StepStatuses[index].ErrMsg = err.Error()

			wsLog.Infof("[%s] - failed: %s", step.Name(), err.Error())
			if err2 := w.sync(ctx); err2 != nil {
				logrus.Errorf("sync error %v for step %s", err2, step.Name())
			}

			if err3 := step.Rollback(ctx, out, w.Config); err3 != nil {
				logrus.Errorf("rollback: step %s : %v", step.Name(), err3)
			}

			return err
		} else {
			wsLog.Infof("[%s] - success", step.Name())
			// Mark step as success
			w.StepStatuses[index].Status = steps.StatusSuccess
			if err := w.sync(ctx); err != nil {
				logrus.Errorf("sync error %v for step %s", err, step.Name())
			}
		}
	}

	return nil
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

	return w.repository.Put(ctx, Prefix, w.ID, buf.Bytes())
}
