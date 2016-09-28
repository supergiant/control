package core

import "github.com/supergiant/supergiant/pkg/model"

type Entrypoints struct {
	Collection
}

func (c *Entrypoints) Create(m *model.Entrypoint) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}

	// Load Kube and CloudAccount
	if err := c.Core.DB.Preload("Nodes").Preload("CloudAccount").Where("name = ?", m.KubeName).First(m.Kube); err != nil {
		return err
	}

	provision := &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  5,
		},
		Core:       c.Core,
		ResourceID: m.UUID,
		Model:      m,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateEntrypoint(m, a)
		},
	}
	return provision.Async()
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
		Fn: func(_ *Action) error {
			// Delete listener records directly
			for _, listener := range m.EntrypointListeners {
				if err := c.Core.DB.Delete(listener); err != nil {
					return err
				}
			}
			if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteEntrypoint(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
