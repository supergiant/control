package client

import "github.com/supergiant/supergiant/pkg/model"

type KubesInterface interface {
	CollectionInterface
	Provision(*int64, *model.Kube) error
}

type Kubes struct {
	Collection
}

func (c *Kubes) Provision(id *int64, m *model.Kube) error {
	return c.client.request("POST", c.memberPath(id)+"/provision", nil, m, nil)
}
