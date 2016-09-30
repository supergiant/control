package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type Volumes struct {
	CreateFn           func(*model.Volume) error
	GetFn              func(*int64, model.Model) error
	GetWithIncludesFn  func(*int64, model.Model, []string) error
	UpdateFn           func(*int64, *model.Volume, *model.Volume) error
	DeleteFn           func(*int64, *model.Volume) core.ActionInterface
	ResizeFn           func(*int64, *model.Volume) core.ActionInterface
	WaitForAvailableFn func(*int64, *model.Volume) error
}

func (c *Volumes) Create(m *model.Volume) error {
	if c.CreateFn == nil {
		return nil
	}
	return c.CreateFn(m)
}

func (c *Volumes) Get(id *int64, m model.Model) error {
	if c.GetFn == nil {
		return nil
	}
	return c.GetFn(id, m)
}

func (c *Volumes) GetWithIncludes(id *int64, m model.Model, includes []string) error {
	if c.GetWithIncludesFn == nil {
		return nil
	}
	return c.GetWithIncludesFn(id, m, includes)
}

func (c *Volumes) Update(id *int64, oldM *model.Volume, m *model.Volume) error {
	if c.UpdateFn == nil {
		return nil
	}
	return c.UpdateFn(id, oldM, m)
}

func (c *Volumes) Delete(id *int64, m *model.Volume) core.ActionInterface {
	if c.DeleteFn == nil {
		return nil
	}
	return c.DeleteFn(id, m)
}

func (c *Volumes) Resize(id *int64, m *model.Volume) core.ActionInterface {
	if c.ResizeFn == nil {
		return nil
	}
	return c.ResizeFn(id, m)
}

func (c *Volumes) WaitForAvailable(id *int64, m *model.Volume) error {
	if c.WaitForAvailableFn == nil {
		return nil
	}
	return c.WaitForAvailableFn(id, m)
}
