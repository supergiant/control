package client

import "github.com/supergiant/supergiant/pkg/model"

type KubeResourcesInterface interface {
	CollectionInterface
	Start(*int64, *model.KubeResource) error
	Stop(*int64, *model.KubeResource) error
}

type KubeResources struct {
	Collection
}

func (c *KubeResources) Start(id *int64, m *model.KubeResource) error {
	return c.client.request("POST", c.memberPath(id)+"/start", nil, m, nil)
}

func (c *KubeResources) Stop(id *int64, m *model.KubeResource) error {
	return c.client.request("POST", c.memberPath(id)+"/stop", nil, m, nil)
}
