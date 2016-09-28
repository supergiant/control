package core

import "github.com/supergiant/supergiant/pkg/model"

type VolumesInterface interface {
	Create(*model.Volume) error
	Get(*int64, model.Model) error
	GetWithIncludes(*int64, model.Model, []string) error
	Update(*int64, *model.Volume, *model.Volume) error
	Delete(*int64, *model.Volume) ActionInterface
	Resize(*int64, *model.Volume) ActionInterface
	WaitForAvailable(*int64, *model.Volume) error
}

type Volumes struct {
	Collection
}

func (c *Volumes) Create(m *model.Volume) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	if err := c.Core.DB.Preload("CloudAccount").Where("name = ?", m.KubeName).First(m.Kube); err != nil {
		return err
	}
	action := &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			// TODO
			// This resource has an issue with retryable provisioning -- which in this
			// context means creating an remote asset from the local record.
			//
			// Apps, for example, which use their user-set Name field as the actual
			// identifier for the provisioned Kubernetes Namespace. That makes the
			// creation of the Namespace retryable, because it is IDEMPOTENT.
			//
			// The problem here, is that WE CANNOT SET AN IDENTIFIER UP FRONT. The ID
			// is given to us upon successful creation of the remote asset.
			//
			// You can tag volumes after creation, but that means it is a 2-step
			// process, which means it fails to be atomic -- if tag creation fails,
			// retrying would re-create a volume, since our identifer (which is used
			// to load and check existence of the asset) was never set.
			//
			// We currently do not have a great solution in place for this problem.
			// In the meantime, MaxRetries is set low to prevent creating several
			// duplicate, billable assets in the user's cloud account. If there is an
			// error, the user will know about it quickly, instead of after 20 retries.
			MaxRetries: 1,
		},
		Core:       c.Core,
		ResourceID: m.UUID,
		Model:      m,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateVolume(m, a)
		},
	}
	return action.Now()
}

func (c *Volumes) Update(id *int64, oldM *model.Volume, m *model.Volume) error {
	if err := c.Collection.Update(id, oldM, m); err != nil {
		return err
	}
	if oldM.Size != m.Size {
		// Resize expects the model arg to be the new size, and will save the record
		// to update. (NOTE this may need a little work. Need to make sure all
		// provider implementations are saving the model on resize.)
		return c.Resize(id, m).Async()
	}
	return nil
}

func (c *Volumes) Delete(id *int64, m *model.Volume) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(_ *Action) error {
			if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteVolume(m); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}

// Resize the Volume
func (c *Volumes) Resize(id *int64, m *model.Volume) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "resizing",
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).ResizeVolume(m, a)
		},
	}
}

func (c *Volumes) WaitForAvailable(id *int64, m *model.Volume) error {
	action := &Action{
		Status: &model.ActionStatus{
			Description: "waiting for available",
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).WaitForVolumeAvailable(m, a)
		},
	}
	return action.Now()
}
