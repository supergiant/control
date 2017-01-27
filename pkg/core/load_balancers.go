package core

import "github.com/supergiant/supergiant/pkg/model"

type LoadBalancers struct {
	Collection
}

func (c *LoadBalancers) Create(m *model.LoadBalancer) error {
	if err := c.Collection.Create(m); err != nil {
		return err
	}
	return c.Core.LoadBalancers.Provision(m.ID, m).Async()
}

func (c *LoadBalancers) Provision(id *int64, m *model.LoadBalancer) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "provisioning",
			MaxRetries:  5,
		},
		Core: c.Core,
		// Nodes are needed to register with ELB on AWS
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).CreateLoadBalancer(m, a)
		},
	}
}

func (c *LoadBalancers) Update(id *int64, oldM *model.LoadBalancer, m *model.LoadBalancer) error {
	if err := c.Collection.Update(id, oldM, m); err != nil {
		return err
	}
	action := &Action{
		Status: &model.ActionStatus{
			Description: "updating",
			MaxRetries:  5,
		},
		Core: c.Core,
		// Nodes are needed to register with ELB on AWS
		Scope: c.Core.DB.Preload("Kube.CloudAccount"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			return c.Core.CloudAccounts.provider(m.Kube.CloudAccount).UpdateLoadBalancer(m, a)
		},
	}
	return action.Async()
}

func (c *LoadBalancers) Delete(id *int64, m *model.LoadBalancer) ActionInterface {
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
			if err := c.Core.CloudAccounts.provider(m.Kube.CloudAccount).DeleteLoadBalancer(m, a); err != nil {
				return err
			}
			return c.Collection.Delete(id, m)
		},
	}
}
