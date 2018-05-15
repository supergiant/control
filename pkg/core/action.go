package core

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"context"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

type ActionInterface interface {
	Now() error
	Async() error
	CancellableWaitFor(string, time.Duration, time.Duration, func() (bool, error)) error
	GetStatus() *model.ActionStatus
}

type Action struct {
	Status         *model.ActionStatus
	Core           *Core
	Scope          DBInterface
	Model          model.Model
	ID             *int64
	ResourceID     string
	Fn             func(*Action) error
	CancelExisting bool
}

//------------------------------------------------------------------------------

type RepeatedActionError struct {
	ResourceID string
}

func (err *RepeatedActionError) Error() string {
	return "Already perform action for " + err.ResourceID
}

//------------------------------------------------------------------------------

func (a *Action) Now() error {
	if err := a.prepare(); err != nil {
		return err
	}

	// TODO we may want some means of communicating with the existing action, to
	// know that it has stopped its goroutines before continuing.

	a.Core.Actions.Put("Begin  : "+a.description(), a.ResourceID, a)

	// Remove Action from map regardless of success or failure
	defer a.stopUnlessCancelled()

	return a.Fn(a)
}

func (a *Action) Async() error {
	if err := a.prepare(); err != nil {
		return err
	}

	a.Core.Actions.Put("Begin  : "+a.description(), a.ResourceID, a)

	go func() {
		for {
			if a.Status.Cancelled {
				break // Goto Remove from Actions
			}

			err := a.Fn(a)
			if err == nil {
				break // Goto Remove from Actions
			}

			a.Status.Error = err.Error()

			a.Core.Log.Error(err)

			if a.Status.Retries >= a.Status.MaxRetries {
				return // Don't goto Remove from Actions
			}

			// TODO this should be configurable exponential backoff
			time.Sleep(time.Second)

			a.Status.Retries++
		}

		// Remove from Actions
		a.stopUnlessCancelled()
	}()

	return nil
}

func (a *Action) CancellableWaitFor(desc string, d time.Duration, i time.Duration, fn func() (bool, error)) error {
	// TODO(stgleb): Context should be inherited from higher level context
	ctx, _ := context.WithTimeout(context.Background(), d)

	return util.WaitFor(desc, ctx, i, func() (bool, error) {
		if a.Status.Cancelled {
			return false, fmt.Errorf("action cancelled while waiting for %s", desc)
		}
		return fn()
	})
}

func (a *Action) GetStatus() *model.ActionStatus {
	return a.Status
}

// Private

func (a *Action) description() string {
	modelType := strings.Split(reflect.TypeOf(a.Model).String(), ".")[1]
	return fmt.Sprintf("%s %s %s", a.Status.Description, modelType, a.ResourceID)
}

func (a *Action) prepare() error {
	// Load model
	if a.ResourceID != "" {
		return nil
	}
	if err := a.Scope.First(a.Model, *a.ID); err != nil {
		return err
	}
	a.ResourceID = a.Model.GetUUID()
	if a.ResourceID == "" {
		return fmt.Errorf("Missing UUID for %#v", a.Model)
	}

	// Prevent concurrent actions (unless existing has failed)
	if ei := a.Core.Actions.Get(a.ResourceID); ei != nil {
		existing := ei.(*Action)
		if a.CancelExisting {
			existing.Status.Cancelled = true
			a.Core.Actions.Delete("Cancel : "+a.description(), a.ResourceID)
		} else if existing.Status.Retries < existing.Status.MaxRetries {
			return &RepeatedActionError{a.ResourceID}
		}
	}

	return nil
}

func (a *Action) stopUnlessCancelled() {
	if !a.Status.Cancelled {
		a.Core.Actions.Delete("End    : "+a.description(), a.ResourceID)
	}
}

////////////////////////////////////////////////////////////////////////////////
//\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
////////////////////////////////////////////////////////////////////////////////
//\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

func (c *Core) SetResourceActionStatus(m model.Model) {
	if ai := c.Actions.Get(m.GetUUID()); ai != nil {
		m.SetActionStatus(ai.(ActionInterface).GetStatus())
	}
}
