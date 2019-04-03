package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"runtime/debug"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type TaskType string

const (
	MasterTask       = "master"
	NodeTask         = "node"
	ClusterTask      = "cluster"
	PreProvisionTask = "preprovision"
)

// Task is an entity that has it own state that can be tracked
// and written to persistent storage through repository, it executes
// particular workflow of steps.
type Task struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Config       *steps.Config   `json:"config"`
	Status       statuses.Status `json:"status"`
	StepStatuses []StepStatus    `json:"stepsStatuses"`

	workflow   Workflow
	repository storage.Interface
}

func NewTask(taskType string, repository storage.Interface) (*Task, error) {
	w := GetWorkflow(taskType)

	if w == nil {
		return nil, sgerrors.ErrNotFound
	}

	t := newTask(taskType, w, repository)

	// This must be done in NewTask
	// Create list of statuses to track
	for _, step := range t.workflow {
		t.StepStatuses = append(t.StepStatuses, StepStatus{
			Status:   statuses.Todo,
			StepName: step.Name(),
			ErrMsg:   "",
		})
	}

	// Set status for task
	t.Status = statuses.Todo

	// Try to sync the task at first time
	err := t.sync(context.Background())

	return t, err

}

func newTask(workflowType string, workflow Workflow, repository storage.Interface) *Task {
	return &Task{
		ID:           uuid.New(),
		Type:         workflowType,
		Status:       statuses.Todo,
		StepStatuses: make([]StepStatus, 0, 0),

		workflow:   workflow,
		repository: repository,
	}
}

// Run executes all steps of workflow and tracks the progress in persistent storage
func (t *Task) Run(ctx context.Context, config steps.Config, out io.WriteCloser) chan error {
	errChan := make(chan error, 1)

	if t.Status == statuses.Success {
		errChan <- nil
		return errChan
	}

	if len(t.workflow) == 0 {
		return errChan
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Status = statuses.Error
				if err := t.sync(ctx); err != nil {
					logrus.Errorf("sync error %v for task %s", err, t.ID)
				}
				debug.PrintStack()
				errChan <- errors.Errorf("provisioning failed, unexpected panic: %v ", r)
			}
		}()

		t.Config = &config

		// Save task state before first step
		if err := t.sync(ctx); err != nil {
			logrus.Errorf("Error saving task state %v", err)
		}

		startIndex := 0
		// Skip successfully finished steps in case of restart
		for index, stepStatus := range t.StepStatuses {
			if stepStatus.Status != statuses.Success {
				startIndex = index
				break
			}
		}

		logrus.Debugf("start task from step #%d startIndex %s",
			startIndex, t.StepStatuses[startIndex].StepName)

		// Start from the first step
		err := t.startFrom(ctx, t.ID, out, startIndex)

		if err != nil {
			if ctx.Err() == context.Canceled {
				t.Status = statuses.Cancelled
				// Save task in cancelled state
				if err := t.sync(context.Background()); err != nil {
					logrus.Errorf("failed to sync task %s to db: %v", t.ID, err)
				}
				errChan <- ctx.Err()
			} else {
				t.Status = statuses.Error
				if err := t.sync(ctx); err != nil {
					logrus.Errorf("failed to sync task %s to db: %v", t.ID, err)
				}
				errChan <- err
			}

			return
		}

		// Set task state to success and save this state
		t.Status = statuses.Success

		if err := t.sync(ctx); err != nil {
			logrus.Errorf("failed to sync task %s to db: %v", t.ID, err)
		}

		logrus.Infof("Task %s has finished successfully", t.ID)
		// Notify provisioner that task output closed with error
		if err := out.Close(); err != nil {
			errChan <- err
		}
		close(errChan)
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

		// sync to storage with task in executing state
		w.Status = statuses.Executing
		w.StepStatuses[index].Status = statuses.Executing

		if err := w.sync(ctx); err != nil {
			logrus.Errorf("sync error %v", err)
		}

		if err := step.Run(ctx, out, w.Config); err != nil {
			// Mark step status as error
			w.StepStatuses[index].Status = statuses.Error
			w.Status = statuses.Error
			w.StepStatuses[index].ErrMsg = err.Error()
			w.sync(ctx)

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
			w.StepStatuses[index].Status = statuses.Success
			w.StepStatuses[index].ErrMsg = ""
			w.Status = statuses.Success
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
