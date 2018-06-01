package core

import (
	"strings"

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

	return c.Core.Kubes.Provision(m.ID, m).Async()
}

func (c *Kubes) Provision(id *int64, m *model.Kube) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  20,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			if err := c.Core.CloudAccounts.provider(m.CloudAccount).CreateKube(m, a); err != nil {
				return err
			}
			return c.Core.DB.Model(m).Update("ready", true)
		},
	}
}

func (c *Kubes) Delete(id *int64, m *model.Kube) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("CloudAccount").Preload("KubeResources").Preload("HelmReleases").Preload("LoadBalancers").Preload("Nodes"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(a *Action) error {
			// Delete Kube Resources directly (don't use provisioner Teardown)
			for _, kubeResource := range m.KubeResources {
				if err := c.Core.DB.Delete(kubeResource); err != nil {
					return err
				}
			}
			for _, loadBalancer := range m.LoadBalancers {
				if err := c.Core.LoadBalancers.Delete(loadBalancer.ID, loadBalancer).Now(); err != nil {
					return err
				}
			}
			
			// Get all kube nodes
			if err := c.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
				return err
			}

			// Delete nodes first to get rid of any potential hanging volumes
			for _, node := range m.Nodes {
				err := c.Core.Nodes.Delete(node.ID, node).Now()
				if err != nil && !strings.Contains("record not found", err.Error()) {
					return err
				}
			}

			// TODO -------------------------------------- and what about Volumes        (maybe we don't have to delete these?)
			// // Delete Volumes
			// for _, volume := range m.Volumes {
			// 	if err := c.Core.Volumes.Delete(volume.ID, volume).Now(); err != nil {
			// 		return err
			// 	}
			// }
			if err := c.Core.CloudAccounts.provider(m.CloudAccount).DeleteKube(m, a); err != nil {
				return err
			}

			// Delete HelmReleases directly (NOTE we do this after because of the periodic sync of these in bg)
			for _, release := range m.HelmReleases {
				if err := c.Core.DB.Delete(release); err != nil {
					return err
				}
			}
			return c.Collection.Delete(id, m)
		},
	}
}
