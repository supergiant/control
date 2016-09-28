package core

import "github.com/supergiant/supergiant/pkg/model"

type Provisioner interface {
	Provision(*model.KubeResource) error
	IsRunning(*model.KubeResource) (bool, error)
	Teardown(*model.KubeResource) error
}
