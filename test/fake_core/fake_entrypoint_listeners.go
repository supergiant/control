package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type EntrypointListeners struct {
	CreateFn          func(*model.EntrypointListener) error
	ProvisionFn       func(*int64, *model.EntrypointListener) core.ActionInterface
	GetFn             func(*int64, model.Model) error
	GetWithIncludesFn func(*int64, model.Model, []string) error
	UpdateFn          func(*int64, model.Model, model.Model) error
	DeleteFn          func(*int64, *model.EntrypointListener) core.ActionInterface
}

func (c *EntrypointListeners) Create(m *model.EntrypointListener) error {
	if c.CreateFn == nil {
		return nil
	}
	return c.CreateFn(m)
}

func (c *EntrypointListeners) Provision(id *int64, m *model.EntrypointListener) core.ActionInterface {
	if c.ProvisionFn == nil {
		return nil
	}
	return c.ProvisionFn(id, m)
}

func (c *EntrypointListeners) Get(id *int64, m model.Model) error {
	if c.GetFn == nil {
		return nil
	}
	return c.GetFn(id, m)
}

func (c *EntrypointListeners) GetWithIncludes(id *int64, m model.Model, includes []string) error {
	if c.GetWithIncludesFn == nil {
		return nil
	}
	return c.GetWithIncludesFn(id, m, includes)
}

func (c *EntrypointListeners) Update(id *int64, oldM model.Model, m model.Model) error {
	if c.UpdateFn == nil {
		return nil
	}
	return c.UpdateFn(id, oldM, m)
}

func (c *EntrypointListeners) Delete(id *int64, m *model.EntrypointListener) core.ActionInterface {
	if c.DeleteFn == nil {
		return nil
	}
	return c.DeleteFn(id, m)
}
