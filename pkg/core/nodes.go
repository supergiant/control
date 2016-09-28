package core

import (
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type Nodes struct {
	Collection
}

func (c *Nodes) Create(m *model.Node) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}

	provision := &Action{
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
		ID:    m.ID,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateNode(m, a)
		},
	}
	return provision.Async()
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
		Fn: func(_ *Action) error {
			if m.ProviderID == "" {
				c.Core.Log.Warnf("Deleting Node %d which has no provider_id", *m.ID)
			} else {
				if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteNode(m); err != nil {
					return err
				}
			}
			return c.Collection.Delete(id, m)
		},
	}
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (c *Nodes) hasPodsWithReservedResources(m *model.Node) (bool, error) {
	k8s := c.Core.K8S(m.Kube)
	pods, err := k8s.ListPods("fieldSelector=spec.nodeName=" + m.Name + ",status.phase=Running")
	if err != nil {
		return false, err
	}

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {

			// TODO
			//
			// These should be moved to json Unmarshal method, so these floats can be parsed once
			//
			gib, err := kubernetes.GiBFromMemString(container.Resources.Requests.Memory)
			if err != nil {
				return false, err
			}
			cores, err := kubernetes.CoresFromCPUString(container.Resources.Requests.CPU)
			if err != nil {
				return false, err
			}

			if gib > 0 || cores > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}
