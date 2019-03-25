package amazon

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestInitDeleteLoadBalancer(t *testing.T) {
	InitDeleteLoadBalancer(GetELB)

	s := steps.GetStep(DeleteLoadBalancerStepName)

	if s == nil {
		t.Errorf("Step %s not found", DeleteLoadBalancerStepName)
	}
}

func TestNewDeleteLoadBalancerStep(t *testing.T) {
	step := NewDeleteLoadBalancerStep(GetELB)

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

func TestNewDeleteLoadBalancerStepError(t *testing.T) {
	fn := func(steps.AWSConfig) (*elb.ELB, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewDeleteLoadBalancerStep(fn)

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

func TestDeleteLoadBalancerStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		getSvcErr error

		deleteExternalLB    *elb.DeleteLoadBalancerOutput
		deleteExternalLBErr error

		deleteInternalLB    *elb.DeleteLoadBalancerOutput
		deleteInternalLBErr error

		errMsg string
	}{
		{
			description: "Error getting ELB svc",
			getSvcErr:   errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description:         "error deleting external LB",
			deleteExternalLBErr: errors.New("error2"),
			errMsg:              "error2",
		},
		{
			description:         "error deleting internal LB",
			deleteExternalLB:    &elb.DeleteLoadBalancerOutput{},
			deleteInternalLBErr: errors.New("error3"),
			errMsg:              "error3",
		},
		{
			description:      "success",
			deleteExternalLB: &elb.DeleteLoadBalancerOutput{},
			deleteInternalLB: &elb.DeleteLoadBalancerOutput{},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(mockELBService)

		svc.On("DeleteLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.deleteExternalLB, testCase.deleteExternalLBErr).Once()

		svc.On("DeleteLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.deleteInternalLB, testCase.deleteInternalLBErr).Once()

		step := &DeleteLoadBalancerStep{
			getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerDeleter, error) {
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

func TestDeleteLoadBalancerStep_Rollback(t *testing.T) {
	step := &DeleteLoadBalancerStep{}

	if err := step.Rollback(context.Background(), nil, nil); err != nil {
		t.Errorf("Unexpected error %v while rolling back", err)
	}
}

func TestDeleteLoadBalancerStep_Depends(t *testing.T) {
	step := &DeleteLoadBalancerStep{}

	if deps := step.Depends(); deps != nil {
		t.Error("Dependencies must ben nil")
	}
}

func TestDeleteLoadBalancerStep_Name(t *testing.T) {
	step := &DeleteLoadBalancerStep{}

	if step.Name() != DeleteLoadBalancerStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteLoadBalancerStepName, step.Name())
	}
}

func TestDeleteLoadBalancerStep_Description(t *testing.T) {
	step := &DeleteLoadBalancerStep{}

	if step.Description() != "Delete external and internal ELB load balancers for master nodes" {
		t.Errorf("Wrong step description expected Delete external and internal ELB load balancers for master nodes actual %s",
			step.Description())
	}
}
