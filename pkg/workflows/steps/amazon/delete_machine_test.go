package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockInstanceDeleter struct {
	mock.Mock
}

func (m *mockInstanceDeleter) DescribeInstancesWithContext(ctx aws.Context,
	req *ec2.DescribeInstancesInput, opts ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.DescribeInstancesOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)

}

func (m *mockInstanceDeleter) TerminateInstancesWithContext(ctx aws.Context,
	req *ec2.TerminateInstancesInput, opts ...request.Option) (*ec2.TerminateInstancesOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.TerminateInstancesOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockInstanceDeleter) CancelSpotInstanceRequestsWithContext(ctx aws.Context, req *ec2.CancelSpotInstanceRequestsInput, opts ...request.Option) (*ec2.CancelSpotInstanceRequestsOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.CancelSpotInstanceRequestsOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteNodeStep_Run(t *testing.T) {
	testCases := []struct {
		description string

		getSvcErr error

		describeErr    error
		describeOutput *ec2.DescribeInstancesOutput

		terminateErr error
		errMsg       string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description: "describe error",
			describeErr: errors.New("message2"),
			errMsg:      "message2",
		},
		{
			description: "reservation empty",
			describeOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{},
			},
			terminateErr: errors.New("message3"),
		},
		{
			description: "instance empty",
			describeOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
			terminateErr: errors.New("message3"),
		},

		{
			description: "terminate error",
			describeOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId: aws.String("instanceID"),
							},
						},
					},
				},
			},
			terminateErr: errors.New("message3"),
			errMsg:       "message3",
		},
		{
			description: "success",
			describeOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId: aws.String("instanceID"),
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockInstanceDeleter{}
		svc.On("DescribeInstancesWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.describeOutput,
			testCase.describeErr)
		svc.On("TerminateInstancesWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything, testCase.terminateErr)

		svc.On("CancelSpotInstanceRequestsWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(nil, nil)
		config := &steps.Config{}
		step := DeleteNodeStep{
			getSvc: func(steps.AWSConfig) (instanceDeleter, error) {
				return svc, testCase.getSvcErr
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s does not contain %s",
				err.Error(), testCase.errMsg)
		}
	}
}

func TestInitDeleteNode(t *testing.T) {
	InitDeleteNode(GetEC2)

	s := steps.GetStep(DeleteNodeStepName)

	if s == nil {
		t.Error("Step must not be nil")
	}
}

func TestNewDeleteNode(t *testing.T) {
	s := NewDeleteNode(GetEC2)

	if s == nil {
		t.Error("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Error("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewDeleteNodeErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDeleteNode(fn)

	if s == nil {
		t.Error("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Error("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestDeleteNodeStep_Depends(t *testing.T) {
	s := &DeleteNodeStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Dependencies must be nil")
	}
}

func TestDeleteNodeStep_Name(t *testing.T) {
	s := &DeleteNodeStep{}

	if name := s.Name(); name != DeleteNodeStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteNodeStepName, name)
	}
}

func TestDeleteNodeStep_Rollback(t *testing.T) {
	s := &DeleteNodeStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error while rollback %v", err)
	}
}

func TestDeleteNodeStep_Description(t *testing.T) {
	s := &DeleteNodeStep{}

	if desc := s.Description(); desc != "Deletes node in aws cluster" {
		t.Errorf("Wrong description expected "+
			"Deletes node in aws cluster actual %s", desc)
	}
}
