package fake_client

import "github.com/supergiant/supergiant/pkg/model"

type Kubes struct {
	Collection
	ProvisionFn func(*int64, *model.Kube) error
}

func (c *Kubes) Provision(id *int64, m *model.Kube) error {
	if c.ProvisionFn == nil {
		return nil
	}
	return c.ProvisionFn(id, m)
}
