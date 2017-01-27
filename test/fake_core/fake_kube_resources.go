package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type KubeResources struct {
	PopulateFn        func() error
	CreateFn          func(*model.KubeResource) error
	GetFn             func(*int64, model.Model) error
	GetWithIncludesFn func(*int64, model.Model, []string) error
	UpdateFn          func(*int64, *model.KubeResource, *model.KubeResource) error
	DeleteFn          func(*int64, *model.KubeResource) core.ActionInterface
	StartFn           func(*int64, *model.KubeResource) core.ActionInterface
	StopFn            func(*int64, *model.KubeResource) core.ActionInterface
	RefreshFn         func(*model.KubeResource) error
}

func (c *KubeResources) Populate() error {
	return c.PopulateFn()
}

func (c *KubeResources) Create(m *model.KubeResource) error {
	return c.CreateFn(m)
}

func (c *KubeResources) Update(id *int64, oldM *model.KubeResource, m *model.KubeResource) error {
	return c.UpdateFn(id, oldM, m)
}

func (c *KubeResources) Get(id *int64, m model.Model) error {
	return c.GetFn(id, m)
}

func (c *KubeResources) GetWithIncludes(id *int64, m model.Model, includes []string) error {
	return c.GetWithIncludesFn(id, m, includes)
}

func (c *KubeResources) Delete(id *int64, m *model.KubeResource) core.ActionInterface {
	return c.DeleteFn(id, m)
}

func (c *KubeResources) Start(id *int64, m *model.KubeResource) core.ActionInterface {
	return c.StartFn(id, m)
}

func (c *KubeResources) Stop(id *int64, m *model.KubeResource) core.ActionInterface {
	return c.StopFn(id, m)
}

func (c *KubeResources) Refresh(m *model.KubeResource) error {
	return c.RefreshFn(m)
}
