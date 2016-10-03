package fake_client

import "github.com/supergiant/supergiant/pkg/model"

type KubeResources struct {
	Collection
	StartFn func(*int64, *model.KubeResource) error
	StopFn  func(*int64, *model.KubeResource) error
}

func (c *KubeResources) Start(id *int64, m *model.KubeResource) error {
	if c.StartFn == nil {
		return nil
	}
	return c.StartFn(id, m)
}

func (c *KubeResources) Stop(id *int64, m *model.KubeResource) error {
	if c.StopFn == nil {
		return nil
	}
	return c.StopFn(id, m)
}
