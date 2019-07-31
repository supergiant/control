package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockEC2 struct {
	mock.Mock
}

func (m *mockEC2) RunInstancesWithContext(ctx aws.Context,
	req *ec2.RunInstancesInput, opts ...request.Option) (*ec2.Reservation, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.Reservation)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockEC2) DescribeInstancesWithContext(ctx aws.Context,
	req *ec2.DescribeInstancesInput, opts ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.DescribeInstancesOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockEC2) WaitUntilInstanceRunningWithContext(ctx aws.Context,
	req *ec2.DescribeInstancesInput, opts ...request.WaiterOption) error {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func TestStepCreateInstance_Run(t *testing.T) {
	testCases := []struct {
		description       string
		isMaster          bool
		getSvcErr         error
		runInstanceErr    error
		runInstanceResp   *ec2.Reservation
		waitErr           error
		describeErr       error
		describeInstances *ec2.DescribeInstancesOutput
		errMsg            string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description:    "run instance error",
			runInstanceErr: errors.New("message2"),
			errMsg:         "message2",
		},
		{
			description:    "run master instance error",
			isMaster:       true,
			runInstanceErr: errors.New("message2"),
			errMsg:         "message2",
		},
		{
			description: "error not instance created",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{},
			},
			errMsg: "no instances created",
		},
		{
			description: "error wait",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("1234"),
					},
				},
			},
			waitErr: errors.New("message3"),
			errMsg:  "message3",
		},
		{
			description: "describe instances error",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("1234"),
					},
				},
			},
			describeErr: errors.New("message4"),
			errMsg:      "message4",
		},
		{
			description: "no public ip found",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("1234"),
					},
				},
			},
			describeInstances: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId: aws.String("1234"),
							},
						},
					},
				},
			},
			errMsg: "no public IP",
		},
		{
			description: "success",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("1234"),
						LaunchTime: &time.Time{},
					},
				},
			},
			describeInstances: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId:       aws.String("1234"),
								PublicIpAddress:  aws.String("10.20.30.40"),
								PrivateIpAddress: aws.String("172.16.0.1"),
								LaunchTime:       &time.Time{},
							},
						},
					},
				},
			},
		},
		{
			description: "success master",
			runInstanceResp: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("1234"),
						LaunchTime: &time.Time{},
					},
				},
			},
			isMaster: true,
			describeInstances: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId:       aws.String("1234"),
								PublicIpAddress:  aws.String("10.20.30.40"),
								PrivateIpAddress: aws.String("172.16.0.1"),
								LaunchTime:       &time.Time{},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		config, err := steps.NewConfig("test", "", profile.Profile{})

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		config.TaskID = uuid.New()
		config.Kube.ID = uuid.New()
		config.IsMaster = testCase.isMaster

		ec2Svc := &mockEC2{}
		ec2Svc.On("RunInstancesWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.runInstanceResp, testCase.runInstanceErr)
		ec2Svc.On("DescribeInstancesWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.describeInstances, testCase.describeErr)
		ec2Svc.On("WaitUntilInstanceRunningWithContext",
			mock.Anything, mock.Anything, mock.Anything).Return(testCase.waitErr)

		step := &StepCreateInstance{
			getSvc: func(steps.AWSConfig) (instanceService, error) {
				return ec2Svc, testCase.getSvcErr
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for {
				select {
				case <-config.NodeChan():
				case <-ctx.Done():
				}
			}
		}()

		err = step.Run(ctx, &bytes.Buffer{}, config)
		cancel()

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message '%s' does not contain '%s'",
				err.Error(), testCase.errMsg)
		}
	}
}

func TestCreateInstanceStepName(t *testing.T) {
	s := StepCreateInstance{}

	if s.Name() != StepNameCreateEC2Instance {
		t.Errorf("Unexpected step name expected %s actual %s", StepNameCreateEC2Instance, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := StepCreateInstance{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}

func TestNewCreateInstance(t *testing.T) {
	step := NewCreateInstance(GetEC2)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Errorf("Wrong GetEC2 function must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateInstanceError(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewCreateInstance(fn)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Errorf("Wrong GetEC2 function must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestInitCreateMachine(t *testing.T) {
	InitCreateMachine(GetEC2)

	s := steps.GetStep(StepNameCreateEC2Instance)

	if s == nil {
		t.Errorf("Step value must not be nil")
	}
}

func TestStepCreateInstance_Depends(t *testing.T) {
	s := &StepCreateInstance{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Unexpected deps value must not be nil")
	}
}

func TestStepCreateInstance_Description(t *testing.T) {
	s := &StepCreateInstance{}

	if desc := s.Description(); desc != "Create EC2 Instance" {
		t.Errorf("Wrong description value expected "+
			"Create EC2 Instance actual %s", s.Description())
	}
}

func TestStepCreateInstance_Name(t *testing.T) {
	s := &StepCreateInstance{}

	if name := s.Name(); name != StepNameCreateEC2Instance {
		t.Errorf("Wrong name value expected %s actual %s",
			StepNameCreateEC2Instance, s.Name())
	}
}

func TestStepCreateInstance_Rollback(t *testing.T) {
	step := &StepCreateInstance{}

	if err := step.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
