package core

import "github.com/supergiant/supergiant/pkg/model"

type EntrypointListenersInterface interface {
	Create(*model.EntrypointListener) error
	Provision(*int64, *model.EntrypointListener) ActionInterface
	Get(*int64, model.Model) error
	GetWithIncludes(*int64, model.Model, []string) error
	Update(*int64, model.Model, model.Model) error
	Delete(*int64, *model.EntrypointListener) ActionInterface
}

type EntrypointListeners struct {
	Collection
}

func (c *EntrypointListeners) Create(m *model.EntrypointListener) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	return c.Core.EntrypointListeners.Provision(m.ID, m).Now()
}

func (c *EntrypointListeners) Provision(id *int64, m *model.EntrypointListener) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Entrypoint.Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Entrypoint.Kube.CloudAccount).CreateEntrypointListener(m, a)
		},
	}
}

func (c *EntrypointListeners) Delete(id *int64, m *model.EntrypointListener) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Entrypoint.Kube.CloudAccount"),
		Model: m,
		ID:    id,
		// ResourceID: m.UUID,
		Fn: func(a *Action) error {
			if err := c.Core.CloudAccounts.provider(m.Entrypoint.Kube.CloudAccount).DeleteEntrypointListener(m, a); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
