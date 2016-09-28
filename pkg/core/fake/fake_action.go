package fake

import (
	"time"

	"github.com/supergiant/supergiant/pkg/model"
)

type Action struct {
	NowFn                func() error
	AsyncFn              func() error
	CancellableWaitForFn func(string, time.Duration, time.Duration, func() (bool, error)) error
	GetStatusFn          func() *model.ActionStatus
}

func (a *Action) Now() error {
	if a.NowFn == nil {
		return nil
	}
	return a.NowFn()
}

func (a *Action) Async() error {
	if a.AsyncFn == nil {
		return nil
	}
	return a.AsyncFn()
}

func (a *Action) CancellableWaitFor(desc string, d time.Duration, i time.Duration, fn func() (bool, error)) error {
	if a.CancellableWaitForFn == nil {
		return nil
	}
	return a.CancellableWaitForFn(desc, d, i, fn)
}

func (a *Action) GetStatus() *model.ActionStatus {
	if a.GetStatusFn == nil {
		return nil
	}
	return a.GetStatusFn()
}
