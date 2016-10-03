package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type Nodes struct {
	CreateFn                       func(*model.Node) error
	GetFn                          func(*int64, model.Model) error
	GetWithIncludesFn              func(*int64, model.Model, []string) error
	UpdateFn                       func(*int64, model.Model, model.Model) error
	DeleteFn                       func(*int64, *model.Node) core.ActionInterface
	HasPodsWithReservedResourcesFn func(*model.Node) (bool, error)
}

func (c *Nodes) Create(m *model.Node) error {
	if c.CreateFn == nil {
		return nil
	}
	return c.CreateFn(m)
}

func (c *Nodes) Get(id *int64, m model.Model) error {
	if c.GetFn == nil {
		return nil
	}
	return c.GetFn(id, m)
}

func (c *Nodes) GetWithIncludes(id *int64, m model.Model, includes []string) error {
	if c.GetWithIncludesFn == nil {
		return nil
	}
	return c.GetWithIncludesFn(id, m, includes)
}

func (c *Nodes) Update(id *int64, oldM model.Model, m model.Model) error {
	if c.UpdateFn == nil {
		return nil
	}
	return c.UpdateFn(id, oldM, m)
}

func (c *Nodes) Delete(id *int64, m *model.Node) core.ActionInterface {
	if c.DeleteFn == nil {
		return nil
	}
	return c.DeleteFn(id, m)
}

func (c *Nodes) HasPodsWithReservedResources(m *model.Node) (bool, error) {
	if c.HasPodsWithReservedResourcesFn == nil {
		return false, nil
	}
	return c.HasPodsWithReservedResourcesFn(m)
}
