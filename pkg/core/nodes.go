package core

import "github.com/supergiant/supergiant/pkg/model"

type NodesInterface interface {
	Create(*model.Node) error
	Provision(*int64, *model.Node) ActionInterface
	Get(*int64, model.Model) error
	GetWithIncludes(*int64, model.Model, []string) error
	Update(*int64, model.Model, model.Model) error
	Delete(*int64, *model.Node) ActionInterface
}

type Nodes struct {
	Collection
}

func (c *Nodes) Create(m *model.Node) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	return c.Core.Nodes.Provision(m.ID, m).Async()
}

func (c *Nodes) Provision(id *int64, m *model.Node) ActionInterface {
	return &Action{
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
			// We currently do not have a great solution in place for this problem.
			// In the meantime, MaxRetries is set low to prevent creating several
			// duplicate, billable assets in the user's cloud account. If there is an
			// error, the user will know about it quickly, instead of after 20 retries.
			MaxRetries: 0,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount").Preload("Kube.Entrypoints.Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateNode(m, a)
		},
	}
}

func (c *Nodes) Delete(id *int64, m *model.Node) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			if m.ProviderID == "" {
				c.Core.Log.Warnf("Deleting Node %d which has no provider_id", *m.ID)
			} else {
				if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteNode(m, a); err != nil {
					return err
				}
			}
			return c.Collection.Delete(id, m)
		},
	}
}
