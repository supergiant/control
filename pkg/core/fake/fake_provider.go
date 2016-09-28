package fake

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type Provider struct {
	ValidateAccountFn            func(*model.CloudAccount) error
	CreateKubeFn                 func(*model.Kube, *core.Action) error
	DeleteKubeFn                 func(*model.Kube) error
	CreateNodeFn                 func(*model.Node, *core.Action) error
	DeleteNodeFn                 func(*model.Node) error
	CreateVolumeFn               func(*model.Volume, *core.Action) error
	KubernetesVolumeDefinitionFn func(*model.Volume) *kubernetes.Volume
	WaitForVolumeAvailableFn     func(*model.Volume, *core.Action) error
	ResizeVolumeFn               func(*model.Volume, *core.Action) error
	DeleteVolumeFn               func(*model.Volume) error
	CreateEntrypointFn           func(*model.Entrypoint, *core.Action) error
	DeleteEntrypointFn           func(*model.Entrypoint) error
	CreateEntrypointListenerFn   func(*model.EntrypointListener) error
	DeleteEntrypointListenerFn   func(*model.EntrypointListener) error
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

func (p *Provider) DeleteKube(m *model.Kube) error {
	if p.DeleteKubeFn == nil {
		return nil
	}
	return p.DeleteKubeFn(m)
}

func (p *Provider) CreateNode(m *model.Node, a *core.Action) error {
	if p.CreateNodeFn == nil {
		return nil
	}
	return p.CreateNodeFn(m, a)
}

func (p *Provider) DeleteNode(m *model.Node) error {
	if p.DeleteNodeFn == nil {
		return nil
	}
	return p.DeleteNodeFn(m)
}

func (p *Provider) CreateVolume(m *model.Volume, a *core.Action) error {
	if p.CreateVolumeFn == nil {
		return nil
	}
	return p.CreateVolumeFn(m, a)
}

func (p *Provider) KubernetesVolumeDefinition(m *model.Volume) *kubernetes.Volume {
	if p.KubernetesVolumeDefinitionFn == nil {
		return nil
	}
	return p.KubernetesVolumeDefinitionFn(m)
}

func (p *Provider) WaitForVolumeAvailable(m *model.Volume, a *core.Action) error {
	if p.WaitForVolumeAvailableFn == nil {
		return nil
	}
	return p.WaitForVolumeAvailableFn(m, a)
}

func (p *Provider) ResizeVolume(m *model.Volume, a *core.Action) error {
	if p.ResizeVolumeFn == nil {
		return nil
	}
	return p.ResizeVolumeFn(m, a)
}

func (p *Provider) DeleteVolume(m *model.Volume) error {
	if p.DeleteVolumeFn == nil {
		return nil
	}
	return p.DeleteVolumeFn(m)
}

func (p *Provider) CreateEntrypoint(m *model.Entrypoint, a *core.Action) error {
	if p.CreateEntrypointFn == nil {
		return nil
	}
	return p.CreateEntrypointFn(m, a)
}

func (p *Provider) DeleteEntrypoint(m *model.Entrypoint) error {
	if p.DeleteEntrypointFn == nil {
		return nil
	}
	return p.DeleteEntrypointFn(m)
}

func (p *Provider) CreateEntrypointListener(m *model.EntrypointListener) error {
	if p.CreateEntrypointListenerFn == nil {
		return nil
	}
	return p.CreateEntrypointListenerFn(m)
}

func (p *Provider) DeleteEntrypointListener(m *model.EntrypointListener) error {
	if p.DeleteEntrypointListenerFn == nil {
		return nil
	}
	return p.DeleteEntrypointListenerFn(m)
}
