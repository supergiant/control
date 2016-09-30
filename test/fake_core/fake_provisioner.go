package fake_core

import "github.com/supergiant/supergiant/pkg/model"

type Provisioner struct {
	ProvisionFn func(*model.KubeResource) error
	IsRunningFn func(*model.KubeResource) (bool, error)
	TeardownFn  func(*model.KubeResource) error
}

func (p *Provisioner) Provision(m *model.KubeResource) error {
	if p.ProvisionFn == nil {
		return nil
	}
	return p.ProvisionFn(m)
}

func (p *Provisioner) IsRunning(m *model.KubeResource) (bool, error) {
	if p.IsRunningFn == nil {
		return true, nil
	}
	return p.IsRunningFn(m)
}

func (p *Provisioner) Teardown(m *model.KubeResource) error {
	if p.TeardownFn == nil {
		return nil
	}
	return p.TeardownFn(m)
}
