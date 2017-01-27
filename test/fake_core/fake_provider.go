package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type Provider struct {
	ValidateAccountFn    func(*model.CloudAccount) error
	CreateKubeFn         func(*model.Kube, *core.Action) error
	DeleteKubeFn         func(*model.Kube, *core.Action) error
	CreateNodeFn         func(*model.Node, *core.Action) error
	DeleteNodeFn         func(*model.Node, *core.Action) error
	CreateLoadBalancerFn func(*model.LoadBalancer, *core.Action) error
	UpdateLoadBalancerFn func(*model.LoadBalancer, *core.Action) error
	DeleteLoadBalancerFn func(*model.LoadBalancer, *core.Action) error
}

func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	if p.ValidateAccountFn == nil {
		return nil
	}
	return p.ValidateAccountFn(m)
}

func (p *Provider) CreateKube(m *model.Kube, a *core.Action) error {
	if p.CreateKubeFn == nil {
		return nil
	}
	return p.CreateKubeFn(m, a)
}

func (p *Provider) DeleteKube(m *model.Kube, a *core.Action) error {
	if p.DeleteKubeFn == nil {
		return nil
	}
	return p.DeleteKubeFn(m, a)
}

func (p *Provider) CreateNode(m *model.Node, a *core.Action) error {
	if p.CreateNodeFn == nil {
		return nil
	}
	return p.CreateNodeFn(m, a)
}

func (p *Provider) DeleteNode(m *model.Node, a *core.Action) error {
	if p.DeleteNodeFn == nil {
		return nil
	}
	return p.DeleteNodeFn(m, a)
}

func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, a *core.Action) error {
	if p.CreateLoadBalancerFn == nil {
		return nil
	}
	return p.CreateLoadBalancerFn(m, a)
}

func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, a *core.Action) error {
	if p.UpdateLoadBalancerFn == nil {
		return nil
	}
	return p.UpdateLoadBalancerFn(m, a)
}

func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, a *core.Action) error {
	if p.DeleteLoadBalancerFn == nil {
		return nil
	}
	return p.DeleteLoadBalancerFn(m, a)
}
