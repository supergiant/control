package amazon

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockELBService struct {
	mock.Mock
}

func (m *mockELBService) CreateLoadBalancerWithContext(ctx aws.Context,
	input *elb.CreateLoadBalancerInput, opts ...request.Option) (*elb.CreateLoadBalancerOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*elb.CreateLoadBalancerOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockELBService) DeleteLoadBalancerWithContext(ctx aws.Context, input *elb.DeleteLoadBalancerInput, opts ...request.Option) (*elb.DeleteLoadBalancerOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*elb.DeleteLoadBalancerOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockELBService) RegisterInstancesWithLoadBalancerWithContext(ctx aws.Context, input *elb.RegisterInstancesWithLoadBalancerInput, opts ...request.Option) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*elb.RegisterInstancesWithLoadBalancerOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestInitCreateLoadBalancer(t *testing.T) {
	InitCreateLoadBalancer(GetELB)

	s := steps.GetStep(StepCreateLoadBalancer)

	if s == nil {
		t.Errorf("Step %s not found", StepCreateLoadBalancer)
	}
}

func TestNewCreateLoadBalancerStep(t *testing.T) {
	step := NewCreateLoadBalancerStep(GetELB)

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

func TestNewCreateLoadBalancerStepError(t *testing.T) {
	fn := func(steps.AWSConfig) (*elb.ELB, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewCreateLoadBalancerStep(fn)

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

func TestCreateLoadBalancerStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		getSvcErr error

		createExternalLB    *elb.CreateLoadBalancerOutput
		createExternalLBErr error

		createInternalLB    *elb.CreateLoadBalancerOutput
		createInternalLBErr error

		errMsg string
	}{
		{
			description: "Error getting ELB svc",
			getSvcErr:   errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description:         "error creating external LB",
			createExternalLBErr: errors.New("error2"),
			errMsg:              "error2",
		},
		{
			description: "error creating internal LB",
			createExternalLB: &elb.CreateLoadBalancerOutput{
				DNSName: aws.String("hello.world"),
			},
			createInternalLBErr: errors.New("error3"),
			errMsg:              "error3",
		},
		{
			description: "success",
			createExternalLB: &elb.CreateLoadBalancerOutput{
				DNSName: aws.String("hello.world"),
			},
			createInternalLB: &elb.CreateLoadBalancerOutput{
				DNSName: aws.String("internal.dns.name"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(mockELBService)

		svc.On("CreateLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.createExternalLB, testCase.createExternalLBErr).Once()

		svc.On("CreateLoadBalancerWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.createInternalLB, testCase.createInternalLBErr).Once()

		step := &CreateLoadBalancerStep{
			getLoadBalancerService: func(cfg steps.AWSConfig) (LoadBalancerCreater, error) {
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

		if testCase.errMsg == "" &&
			config.ExternalDNSName == "" {
			t.Errorf("Wrong ExternalDNSName must not be empty")
		}

		if testCase.errMsg == "" &&
			config.InternalDNSName == "" {
			t.Errorf("Wrong InternalDNSName must not be empty")
		}
	}
}

func TestCreateLoadBalancerStep_Rollback(t *testing.T) {
	step := &CreateLoadBalancerStep{}

	if err := step.Rollback(context.Background(), nil, nil); err != nil {
		t.Errorf("Unexpected error %v while rolling back", err)
	}
}

func TestCreateLoadBalancerStep_Depends(t *testing.T) {
	step := &CreateLoadBalancerStep{}

	if deps := step.Depends(); len(deps) != 2 {
		t.Errorf("Dependencies must be %v actual %v",
			[]string{StepCreateSubnets, StepCreateSecurityGroups}, deps)
	}
}

func TestCreateLoadBalancerStep_Name(t *testing.T) {
	step := &CreateLoadBalancerStep{}

	if step.Name() != StepCreateLoadBalancer {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateLoadBalancer, step.Name())
	}
}

func TestCreateLoadBalancerStep_Description(t *testing.T) {
	step := &CreateLoadBalancerStep{}

	if step.Description() != "Create ELB external and internal load balancers for master nodes" {
		t.Errorf("Wrong step description expected Create ELB external and internal load balancers for master nodes %s",
			step.Description())
	}
}
