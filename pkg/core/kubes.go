package core

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

type Kubes struct {
	Collection
}

func (c *Kubes) Create(m *model.Kube) error {
	// Defaults
	if m.Username == "" && m.Password == "" {
		m.Username = util.RandomString(16)
		m.Password = util.RandomString(8)
	}

	if err := c.Collection.Create(m); err != nil {
		return err
	}

	// TODO need a validation to make sure CloudAccount matches the provided config

	provision := &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  20,
		},
		Core:       c.Core,
		ResourceID: m.UUID,
		Model:      m,
		Fn: func(a *Action) error {
			if err := c.Core.CloudAccounts.provider(m.CloudAccount).CreateKube(m, a); err != nil {
				return err
			}
			return c.Core.DB.Model(m).Update("ready", true)
		},
	}
	return provision.Async()
}

func (c *Kubes) Delete(id *int64, m *model.Kube) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("CloudAccount").Preload("KubeResources").Preload("Entrypoints").Preload("Volumes").Preload("Nodes"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(_ *Action) error {
			// Delete Kube Resources directly (don't use provisioner Teardown)
			for _, kubeResource := range m.KubeResources {
				if err := c.Core.DB.Delete(kubeResource); err != nil {
					return err
				}
			}
			for _, entrypoint := range m.Entrypoints {
				if err := c.Core.Entrypoints.Delete(entrypoint.ID, entrypoint).Now(); err != nil {
					return err
				}
			}
			// Delete nodes first to get rid of any potential hanging volumes
			for _, node := range m.Nodes {
				if err := c.Core.Nodes.Delete(node.ID, node).Now(); err != nil {
					return err
				}
			}
			// Delete Volumes
			for _, volume := range m.Volumes {
				if err := c.Core.Volumes.Delete(volume.ID, volume).Now(); err != nil {
					return err
				}
			}
			if err := c.Core.CloudAccounts.provider(m.CloudAccount).DeleteKube(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
