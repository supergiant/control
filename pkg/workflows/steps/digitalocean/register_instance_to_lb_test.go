package digitalocean

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
	"strings"
	"github.com/supergiant/control/pkg/model"
)

func TestNewRegisterInstanceToLBStep(t *testing.T) {
	step := NewRegisterInstanceToLBStep()

	if step.getServices == nil {
		t.Errorf("get services must not be nil")
	}

	if client := step.getServices("access token"); client == nil {
		t.Errorf("Client must not be nil")
	}
}

func TestRegisterInstanceToLBStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		InstanceID string
		externalResp *godo.Response
		deleteExErr  error

		internalResp *godo.Response
		deleteInErr  error

		errMsg string
	}{
		{
			description: "Error converting instance id from string to id",

			InstanceID: "not a number",
			errMsg:      "converting",
		},
		{
			description: "Error registering instance to external LB",

			InstanceID: "1",

			deleteExErr: errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description: "Error registering instance to internal LB",

			InstanceID: "1",

			externalResp: &godo.Response{},
			deleteInErr:  errors.New("error2"),

			errMsg: "error2",
		},
		{
			description: "success",

			InstanceID: "1",

			internalResp: &godo.Response{},
			externalResp: &godo.Response{},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &MockLBService{}

		svc.On("AddDroplets", mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.externalResp, testCase.deleteExErr).Once()

		svc.On("AddDroplets", mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.internalResp, testCase.deleteInErr).Once()

		step := &RegisterInstanceToLBStep{
			getServices: func(accessToken string) LoadBalancerService {
				return svc
			},
		}

		config := &steps.Config{
			Node: model.Machine{
				ID: testCase.InstanceID,
			},
		}

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

func TestRegisterInstanceToLBStep_Name(t *testing.T) {
	step := NewRegisterInstanceToLBStep()

	if step.Name() != RegisterInstanceToLB {
		t.Errorf("Wrong step name expected %s actual %s",
			RegisterInstanceToLB, step.Name())
	}
}

func TestRegisterInstanceToLBStep_Depends(t *testing.T) {
	step := NewRegisterInstanceToLBStep()

	if step.Depends() != nil {
		t.Error("Register instance to load balancer step depends must be nil")
	}
}

func TestRegisterInstanceToLBStep_Description(t *testing.T) {
	step := NewRegisterInstanceToLBStep()

	if step.Description() != "Register instance to load balancers in Digital Ocean" {
		t.Errorf("Wrong step description expected Register instance to load balancers in Digital Ocean %s", step.Description())
	}
}

func TestRegisterInstanceToLBStep_Rollback(t *testing.T) {
	s := RegisterInstanceToLBStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
