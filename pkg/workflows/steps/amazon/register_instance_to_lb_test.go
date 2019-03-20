package amazon

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestInitRegisterInstance(t *testing.T) {
	InitRegisterInstance(GetELB)

	s := steps.GetStep(RegisterInstanceStepName)

	if s == nil {
		t.Errorf("Step %s not found", RegisterInstanceStepName)
	}
}

func TestNewRegisterInstanceStep(t *testing.T) {
	step := NewRegisterInstanceStep(GetELB)

	if step == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if step.getLoadBalancerService == nil {
		t.Errorf("getLoadBalancerService must not be nil")
	}

	if api, err := step.getLoadBalancerService(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewRegisterInstanceStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (*elb.ELB, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewRegisterInstanceStep(fn)

	if step == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if step.getLoadBalancerService == nil {
		t.Errorf("getLoadBalancerService must not be nil")
	}

	if api, err := step.getLoadBalancerService(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestRegisterInstanceStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		getSvcErr error

		registerExternalLB    *elb.RegisterInstancesWithLoadBalancerOutput
		registerExternalLBErr error

		deleteInternalLB    *elb.RegisterInstancesWithLoadBalancerOutput
		deleteInternalLBErr error

		errMsg string
	}{
		{
			description: "Error getting ELB svc",
			getSvcErr:   errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description:           "error deleting external LB",
			registerExternalLBErr: errors.New("error2"),
			errMsg:                "error2",
		},
		{
			description:         "error deleting internal LB",
			registerExternalLB:  &elb.RegisterInstancesWithLoadBalancerOutput{},
			deleteInternalLBErr: errors.New("error3"),
			errMsg:              "error3",
		},
		{
			description:        "success",
			registerExternalLB: &elb.RegisterInstancesWithLoadBalancerOutput{},
			deleteInternalLB:   &elb.RegisterInstancesWithLoadBalancerOutput{},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(mockELBService)

		svc.On("RegisterInstancesWithLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.registerExternalLB, testCase.registerExternalLBErr).Once()

		svc.On("RegisterInstancesWithLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.deleteInternalLB, testCase.deleteInternalLBErr).Once()

		step := &RegisterInstanceStep{
			getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerRegister, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			ClusterID: "1234",
			AWSConfig: steps.AWSConfig{},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err != nil && testCase.errMsg == "" {
			t.Errorf("Unexpected error %v", err)
			continue
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Wrong error must contain %s actual %s",
				testCase.errMsg, err.Error())
			continue
		}
	}
}

func TestRegisterInstanceStep_Rollback(t *testing.T) {
	step := &RegisterInstanceStep{}

	if err := step.Rollback(context.Background(), nil, nil); err != nil {
		t.Errorf("Unexpected error %v while rolling back", err)
	}
}

func TestRegisterInstanceStep_Depends(t *testing.T) {
	step := &RegisterInstanceStep{}

	if deps := step.Depends(); deps != nil {
		t.Error("Dependencies must ben nil")
	}
}

func TestRegisterInstanceStep_Name(t *testing.T) {
	step := &RegisterInstanceStep{}

	if step.Name() != RegisterInstanceStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			RegisterInstanceStepName, step.Name())
	}
}

func TestRegisterInstanceStep_Description(t *testing.T) {
	step := &RegisterInstanceStep{}

	if step.Description() != "Register node to external and internal Load balancers" {
		t.Errorf("Wrong step description expected Register node to external and internal Load balancers %s",
			step.Description())
	}
}
