package core

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/supergiant/guber"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// TODO
var globalK8SHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func (c *Core) K8S(m *model.Kube) guber.Client {
	return guber.NewClient(m.MasterPublicIP, m.Username, m.Password, globalK8SHTTPClient)
}

//------------------------------------------------------------------------------

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

	provision := &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  20,
		},
		core:       c.core,
		resourceID: m.UUID,
		model:      m,
		fn: func(a *Action) error {
			if err := c.core.CloudAccounts.provider(m.CloudAccount).CreateKube(m, a); err != nil {
				return err
			}
			return c.core.DB.Model(m).Update("ready", true).Error
		},
	}
	return provision.Async()
}

func (c *Kubes) Delete(id *int64, m *model.Kube) *Action {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		core:           c.core,
		scope:          c.core.DB.Preload("CloudAccount").Preload("Entrypoints").Preload("Volumes.Kube.CloudAccount").Preload("Apps.Components.Instances").Preload("Apps.Components.Releases").Preload("Nodes.Kube.CloudAccount"),
		model:          m,
		id:             id,
		cancelExisting: true,
		fn: func(_ *Action) error {
			for _, entrypoint := range m.Entrypoints {
				if err := c.core.Entrypoints.Delete(entrypoint.ID, entrypoint).Now(); err != nil {
					return err
				}
			}
			// Delete nodes first to get rid of any potential hanging volumes
			for _, node := range m.Nodes {
				if err := c.core.Nodes.Delete(node.ID, node).Now(); err != nil {
					return err
				}
			}

			// Delete App records -- no need to delete assets
			// TODO... might be better to have Kubernetes-related operations first
			// check to see if Kube is flagged for delete?
			for _, app := range m.Apps {
				if err := c.core.DB.Delete(app); err != nil {
					return err
				}
				for _, component := range app.Components {
					if err := c.core.DB.Delete(component); err != nil {
						return err
					}
					for _, release := range component.Releases {
						if err := c.core.DB.Delete(release); err != nil {
							return err
						}
					}
					for _, instance := range component.Instances {
						if err := c.core.DB.Delete(instance); err != nil {
							return err
						}
					}
				}
			}

			// Delete Volumes
			for _, volume := range m.Volumes {
				if err := c.core.Volumes.Delete(volume.ID, volume).Now(); err != nil {
					return err
				}
			}
			if err := c.core.CloudAccounts.provider(m.CloudAccount).DeleteKube(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
