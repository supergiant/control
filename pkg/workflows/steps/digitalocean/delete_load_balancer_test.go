package digitalocean

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestNewDeleteLoadBalancerStep(t *testing.T) {
	step := NewCreateLoadBalancerStep()

	if step.getServices == nil {
		t.Errorf("get services must not be nil")
	}

	if client := step.getServices("access token"); client == nil {
		t.Errorf("Client must not be nil")
	}
}

func TestDeleteLoadBalancerStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		externalResp *godo.Response
		deleteExErr  error

		internalResp *godo.Response
		deleteInErr  error

		errMsg string
	}{
		{
			description: "Error deleting external LB",

			deleteExErr: errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description: "Error deleting internal LB",

			externalResp: &godo.Response{},
			deleteInErr:  errors.New("error2"),

			errMsg: "error2",
		},
		{
			description: "success",

			internalResp: &godo.Response{},
			externalResp: &godo.Response{},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &MockLBService{}

		svc.On("Delete", mock.Anything, mock.Anything).
			Return(testCase.externalResp, testCase.deleteExErr).Once()

		svc.On("Delete", mock.Anything, mock.Anything).
			Return(testCase.internalResp, testCase.deleteInErr).Once()

		step := &DeleteLoadBalancerStep{
			getServices: func(accessToken string) LoadBalancerService {
				return svc
			},
		}

		config := &steps.Config{}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if testCase.errMsg != "" && err == nil {
			t.Errorf("Error not must be nil")
		}

		if testCase.errMsg == "" && err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s must contain %s", err.Error(), testCase.errMsg)
		}
	}
}

func TestDeleteLoadBalancerStep_Name(t *testing.T) {
	step := NewDeleteLoadBalancerStep()

	if step.Name() != DeleteLoadBalancerStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteLoadBalancerStepName, step.Name())
	}
}

func TestDeleteLoadBalancerStep_Depends(t *testing.T) {
	step := NewDeleteLoadBalancerStep()

	if step.Depends() != nil {
		t.Error("Delete load balancer step depends must be nil")
	}
}

func TestDeleteLoadBalancerStep_Description(t *testing.T) {
	step := NewDeleteLoadBalancerStep()

	if step.Description() != "Delete external and internal load balancers in Digital Ocean" {
		t.Errorf("Wrong step description expected Delete external and internal load balancers in Digital Ocean %s", step.Description())
	}
}

func TestDeleteLoadBalancerStep_Rollback(t *testing.T) {
	s := DeleteLoadBalancerStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
