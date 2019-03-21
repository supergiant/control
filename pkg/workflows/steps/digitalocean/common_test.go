package digitalocean

import (
	"testing"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestInit(t *testing.T) {
	Init()

	createMachine := steps.GetStep(CreateMachineStepName)

	if createMachine == nil {
		t.Errorf("%s must not be nil", CreateMachineStepName)
	}
	deleteMachine := steps.GetStep(DeleteMachineStepName)

	if deleteMachine == nil {
		t.Errorf("%s must not be nil", DeleteMachineStepName)
	}

	deleteMachines := steps.GetStep(DeleteClusterMachines)

	if deleteMachines == nil {
		t.Errorf("%s must not be nil", DeleteClusterMachines)
	}

	deleteKeys := steps.GetStep(DeleteDeleteKeysStepName)

	if deleteKeys == nil {
		t.Errorf("%s must not be nil", DeleteDeleteKeysStepName)
	}

	createLB := steps.GetStep(CreateLoadBalancerStepName)

	if createLB == nil {
		t.Errorf("%s must not be nil", CreateLoadBalancerStepName)
	}

	deleteLB := steps.GetStep(DeleteLoadBalancerStepName)

	if deleteLB == nil {
		t.Errorf("%s must not be nil", DeleteLoadBalancerStepName)
	}

	registerInstanceToLB := steps.GetStep(RegisterInstanceToLB)

	if registerInstanceToLB == nil {
		t.Errorf("%s must not be nil", RegisterInstanceToLB)
	}
}
