package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestDeleteClusterMachibesStep_Run(t *testing.T) {
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
		step := DeleteClusterMachines{
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
func TestNewDeleteClusterInstances(t *testing.T) {
	s := NewDeleteClusterInstances(GetEC2)

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

func TestNewDeleteClusterInstancesErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDeleteClusterInstances(fn)

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

func TestInitDeleteClusterMachines(t *testing.T) {
	InitDeleteClusterMachines(GetEC2)

	s := steps.GetStep(DeleteClusterMachinesStepName)

	if s == nil {
		t.Errorf("Step must not be nil")
	}
}

func TestDeleteClusterMachines_Depends(t *testing.T) {
	s := &DeleteClusterMachines{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("depencies must be nil")
	}
}

func TestDeleteClusterMachines_Name(t *testing.T) {
	s := &DeleteClusterMachines{}

	if name := s.Name(); name != DeleteClusterMachinesStepName {
		t.Errorf("Wrong name expected %s actual %s",
			DeleteClusterMachinesStepName, name)
	}
}

func TestDeleteClusterMachines_Rollback(t *testing.T) {
	s := &DeleteClusterMachines{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error while rollback %v", err)
	}
}

func TestDeleteClusterMachines_Description(t *testing.T) {
	s := &DeleteClusterMachines{}

	if desc := s.Description(); desc != "Deletes all nodes in aws cluster" {
		t.Errorf("Wrong description expected Deletes all nodes "+
			"in aws cluster actual %s", desc)
	}
}
