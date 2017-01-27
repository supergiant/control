package core

import "github.com/supergiant/supergiant/pkg/model"

type Provider interface {
	ValidateAccount(*model.CloudAccount) error

	CreateKube(*model.Kube, *Action) error
	DeleteKube(*model.Kube, *Action) error

	CreateNode(*model.Node, *Action) error
	DeleteNode(*model.Node, *Action) error

	CreateLoadBalancer(*model.LoadBalancer, *Action) error
	UpdateLoadBalancer(*model.LoadBalancer, *Action) error
	DeleteLoadBalancer(*model.LoadBalancer, *Action) error
}
