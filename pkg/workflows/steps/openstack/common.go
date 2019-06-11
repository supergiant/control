package openstack

import (
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	StatusOnline               = "ONLINE"
	StatusOffline              = "OFFLINE"
	StatusActive               = "ACTIVE"
)

func Init() {
	steps.RegisterStep(FindImageStepName, NewFindImageStep())
	steps.RegisterStep(CreateNetworkStepName, NewCreateNetworkStep())
	steps.RegisterStep(CreateSubnetStepName, NewCreateSubnetStep())
	steps.RegisterStep(CreateRouterStepName, NewCreateRouterStep())
	steps.RegisterStep(CreateMachineStepName, NewCreateMachineStep())
}
