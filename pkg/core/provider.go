package core

import (
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
)

type Provider interface {
	ValidateAccount(*model.CloudAccount) error

	CreateKube(*model.Kube, *Action) error
	DeleteKube(*model.Kube) error

	CreateNode(*model.Node, *Action) error
	DeleteNode(*model.Node) error

	CreateVolume(*model.Volume, *Action) error
	KubernetesVolumeDefinition(*model.Volume) *kubernetes.Volume
	WaitForVolumeAvailable(*model.Volume, *Action) error
	ResizeVolume(*model.Volume, *Action) error
	DeleteVolume(*model.Volume) error

	CreateEntrypoint(*model.Entrypoint, *Action) error
	DeleteEntrypoint(*model.Entrypoint) error

	CreateEntrypointListener(*model.EntrypointListener) error
	DeleteEntrypointListener(*model.EntrypointListener) error
}
