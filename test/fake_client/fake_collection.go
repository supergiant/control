package fake_client

import "github.com/supergiant/supergiant/pkg/model"

type Collection struct {
	ListFn            func(model.List) error
	CreateFn          func(model.Model) error
	GetFn             func(interface{}, model.Model) error
	GetWithIncludesFn func(interface{}, model.Model, []string) error
	UpdateFn          func(interface{}, model.Model) error
	DeleteFn          func(interface{}, model.Model) error
}

func (c *Collection) List(list model.List) error {
	if c.ListFn == nil {
		return nil
	}
	return c.ListFn(list)
}

func (c *Collection) Create(m model.Model) error {
	if c.CreateFn == nil {
		return nil
	}
	return c.CreateFn(m)
}

func (c *Collection) Get(id interface{}, m model.Model) error {
	if c.GetFn == nil {
		return nil
	}
	return c.GetFn(id, m)
}

func (c *Collection) GetWithIncludes(id interface{}, m model.Model, includes []string) error {
	if c.GetWithIncludesFn == nil {
		return nil
	}
	return c.GetWithIncludesFn(id, m, includes)
}

func (c *Collection) Update(id interface{}, m model.Model) error {
	if c.UpdateFn == nil {
		return nil
	}
	return c.UpdateFn(id, m)
}

func (c *Collection) Delete(id interface{}, m model.Model) error {
	if c.DeleteFn == nil {
		return nil
	}
	return c.DeleteFn(id, m)
}
