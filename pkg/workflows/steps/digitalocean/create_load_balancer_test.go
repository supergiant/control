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
	"time"
)

type MockLBService struct {
	mock.Mock
}

func (m *MockLBService) Create(ctx context.Context, req *godo.LoadBalancerRequest) (*godo.LoadBalancer, *godo.Response, error) {
	args := m.Called(ctx, req)
	val, ok := args.Get(0).(*godo.LoadBalancer)
	if !ok {
		return nil, nil, args.Error(1)
	}
	return val, nil, args.Error(1)
}

func (m *MockLBService) Delete(ctx context.Context, lbID string) (*godo.Response, error) {
	args := m.Called(ctx, lbID)
	val, ok := args.Get(0).(*godo.Response)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *MockLBService) Get(ctx context.Context, lbID string) (*godo.LoadBalancer, *godo.Response, error) {
	args := m.Called(ctx, lbID)
	val, ok := args.Get(0).(*godo.LoadBalancer)
	if !ok {
		return nil, nil, args.Error(1)
	}
	return val, nil, args.Error(1)
}

func TestNewCreateLoadBalancerStep(t *testing.T) {
	step := NewCreateLoadBalancerStep()

	if step.getServices == nil {
		t.Errorf("get services must not be nil")
	}

	if client := step.getServices("access token"); client == nil {
		t.Errorf("Client must not be nil")
	}
}

func TestCreateLoadBalancerStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		createExternalLB    *godo.LoadBalancer
		createExternalLBErr error

		getExternalLB *godo.LoadBalancer
		getExternalLBErr error

		createInternalLB    *godo.LoadBalancer
		createInternalLBErr error

		getInternalLB *godo.LoadBalancer
		getInternalLBErr error


		errMsg string
	}{
		{
			description: "Error creating external LB",

			createExternalLBErr: errors.New("error1"),
			errMsg:              "error1",
		},
		{
			description: "Error getting external LB",

			createExternalLB:    &godo.LoadBalancer{},
			getExternalLBErr: errors.New("error2"),

			errMsg: "error2",
		},
		{
			description: "Error creating internal LB",

			createExternalLB:    &godo.LoadBalancer{},
			getExternalLB: &godo.LoadBalancer{
				Status: "active",
				IP: "10.20.30.40",
			},

			createInternalLBErr: errors.New("error3"),

			errMsg: "error3",
		},
		{
			description: "Error getting internal LB",

			createExternalLB: &godo.LoadBalancer{},
			getExternalLB: &godo.LoadBalancer{
				Status: "active",
				IP: "10.20.30.40",
			},
			createInternalLB: &godo.LoadBalancer{},
			getInternalLBErr: errors.New("error4"),
			errMsg: "error4",
		},
		{
			description: "inactive LB",

			createExternalLB: &godo.LoadBalancer{},
			getExternalLB: &godo.LoadBalancer{
				Status: "",
				IP: "",
			},
			errMsg: "IP must not be empty",
		},
		{
			description: "success",

			createExternalLB: &godo.LoadBalancer{},
			getExternalLB: &godo.LoadBalancer{
				Status: "active",
				IP: "10.20.30.40",
			},
			createInternalLB: &godo.LoadBalancer{},
			getInternalLB: &godo.LoadBalancer{
				Status: "active",
				IP: "11.22.33.44",
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &MockLBService{}

		svc.On("Create", mock.Anything, mock.Anything).
			Return(testCase.createExternalLB, testCase.createExternalLBErr).Once()

		svc.On("Get", mock.Anything, mock.Anything).
			Return(testCase.getExternalLB, testCase.getExternalLBErr).Once()

		svc.On("Create", mock.Anything, mock.Anything).
			Return(testCase.createInternalLB, testCase.createInternalLBErr).Once()

		svc.On("Get", mock.Anything, mock.Anything).
			Return(testCase.getInternalLB, testCase.getInternalLBErr).Once()

		step := &CreateLoadBalancerStep{
			getServices: func(accessToken string) LoadBalancerService {
				return svc
			},
			Attempts: 1,
			Timeout: time.Nanosecond * 1,
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

		if err == nil && (strings.Compare(config.ExternalDNSName, testCase.getExternalLB.IP) != 0 ||
			strings.Compare(config.InternalDNSName, testCase.getInternalLB.IP) != 0) {
			t.Log(testCase.getExternalLB.ID, testCase.getInternalLB.ID)
			t.Errorf("External or Internal DNS names do not correspond actual output")
		}
	}
}

func TestCreateLoadBalancerStep_Name(t *testing.T) {
	step := NewCreateLoadBalancerStep()

	if step.Name() != CreateLoadBalancerStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			CreateLoadBalancerStepName, step.Name())
	}
}

func TestCreateLoadBalancerStep_Depends(t *testing.T) {
	step := NewCreateLoadBalancerStep()

	if step.Depends() != nil {
		t.Error("Create load balancer step depends must be nil")
	}
}

func TestCreateLoadBalancerStep_Description(t *testing.T) {
	step := NewCreateLoadBalancerStep()

	if step.Description() != "Create load balancer in Digital Ocean" {
		t.Errorf("Wrong step description expected Create load balancer in Digital Ocean %s", step.Description())
	}
}

func TestCreateLoadBalancerStep_Rollback(t *testing.T) {
	s := CreateLoadBalancerStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
