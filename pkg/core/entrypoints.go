package core

import "github.com/supergiant/supergiant/pkg/model"

type Entrypoints struct {
	Collection
}

func (c *Entrypoints) Create(m *model.Entrypoint) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	return c.Core.Entrypoints.Provision(m.ID, m).Async()
}

func (c *Entrypoints) Provision(id *int64, m *model.Entrypoint) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  5,
		},
		Core: c.Core,
		// Nodes are needed to register with ELB on AWS
		Scope: c.Core.DB.Preload("Kube.CloudAccount").Preload("Kube.Nodes"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateEntrypoint(m, a)
		},
	}
}

func (c *Entrypoints) Delete(id *int64, m *model.Entrypoint) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount").Preload("EntrypointListeners"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			// Delete listener records directly
			for _, listener := range m.EntrypointListeners {
				if err := c.Core.DB.Delete(listener); err != nil {
					return err
				}
			}
			if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteEntrypoint(m, a); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
