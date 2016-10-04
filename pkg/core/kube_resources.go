package core

import (
	"fmt"
	"time"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

type KubeResourcesInterface interface {
	Create(*model.KubeResource) error
	Get(*int64, model.Model) error
	GetWithIncludes(*int64, model.Model, []string) error
	Update(*int64, *model.KubeResource, *model.KubeResource) error
	Delete(*int64, *model.KubeResource) ActionInterface
	Start(*int64, *model.KubeResource) ActionInterface
	Stop(*int64, *model.KubeResource) ActionInterface
	Refresh(*model.KubeResource) error
}

type KubeResources struct {
	Collection
}

func (c *KubeResources) Create(m *model.KubeResource) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	// NOTE we call this from core to get the interface
	return c.Core.KubeResources.Start(m.ID, m).Async()
}

// TODO
func (c *KubeResources) Update(id *int64, oldM *model.KubeResource, m *model.KubeResource) error {
	return c.Collection.Update(id, oldM, m)
}

func (c *KubeResources) Delete(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("Kube.CloudAccount"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(_ *Action) error {
			if err := c.provisioner(m).Teardown(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}

func (c *KubeResources) Start(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "starting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			k8s := c.Core.K8S(m.Kube)
			if err := k8s.EnsureNamespace(m.Namespace); err != nil {
				return err
			}
			if err := c.provisioner(m).Provision(m); err != nil {
				return err
			}
			// Wait for Resource to be ready
			desc := fmt.Sprintf("%s '%s' in Namespace '%s' to start", m.Kind, m.Name, m.Namespace)
			waitErr := util.WaitFor(desc, c.Core.KubeResourceStartTimeout, 3*time.Second, func() (bool, error) {
				return c.provisioner(m).IsRunning(m)
			})
			if waitErr != nil {
				return waitErr
			}
			return c.Core.DB.Model(m).Update("started", true)
		},
	}
}

func (c *KubeResources) Stop(id *int64, m *model.KubeResource) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "stopping",
			MaxRetries:  5,
		},
		Core:           c.Core,
		Scope:          c.Core.DB.Preload("Kube.CloudAccount"),
		Model:          m,
		ID:             id,
		CancelExisting: true,
		Fn: func(a *Action) error {
			// NOTE we use default provisioner teardown here, because it works as a
			// stop() for all resources... it does not clean up any assets, which is
			// just what we need here.
			provisioner := &DefaultProvisioner{c.Core}
			if err := provisioner.Teardown(m); err != nil {
				return err
			}
			return c.Core.DB.Model(m).Update("started", false)
		},
	}
}

func (c *KubeResources) Refresh(m *model.KubeResource) (err error) {
	m.Started, err = c.provisioner(m).IsRunning(m)
	if err != nil {
		return err
	}
	// Save Artifact and Started
	return c.Core.DB.Save(m)
}

// Private

func (c *KubeResources) provisioner(m *model.KubeResource) Provisioner {
	switch m.Kind {
	case "Pod":
		return c.Core.PodProvisioner
	case "Service":
		return c.Core.ServiceProvisioner
	}
	return c.Core.DefaultProvisioner
}
