package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"testing"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/pkg/errors"
)

type fakeEC2CreateMachine struct {
	ec2iface.EC2API

	runInstancesOutput  *ec2.Reservation
	descInstancesOutput *ec2.DescribeInstancesOutput
	runErr              error
	descErr             error
}

func (f *fakeEC2CreateMachine) RunInstancesWithContext(aws.Context, *ec2.RunInstancesInput, ...request.Option) (*ec2.Reservation, error) {
	return f.runInstancesOutput, f.runErr
}

func (f *fakeEC2CreateMachine) DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return f.descInstancesOutput, f.descErr
}

//TODO: FIX THE TESTS AFTER THE DEMO

//func TestStepCreateInstance_Run(t *testing.T) {
//	tt := []struct {
//		fn          GetEC2Fn
//		importErr         error
//		hasPublicIP bool
//	}{
//		{
//			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
//				return &fakeEC2CreateMachine{
//					runInstancesOutput: &ec2.Reservation{
//						Instances: []*ec2.Instance{
//							{
//								LaunchTime:      aws.Time(time.Now()),
//								InstanceId:      aws.String("ID"),
//								PublicIpAddress: aws.String("127.0.0.1"),
//							},
//						},
//					},
//				}, nil
//			},
//			importErr:         nil,
//			hasPublicIP: false,
//		},
//		{
//			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
//				return &fakeEC2CreateMachine{
//					runInstancesOutput: &ec2.Reservation{
//						Instances: []*ec2.Instance{
//							{
//								LaunchTime:       aws.Time(time.Now()),
//								InstanceId:       aws.String("ID"),
//								PublicIpAddress:  aws.String("127.0.0.1"),
//								PrivateIpAddress: aws.String("127.0.0.1"),
//							},
//						},
//					},
//					descInstancesOutput: &ec2.DescribeInstancesOutput{
//						Reservations: []*ec2.Reservation{
//							{
//								Instances: []*ec2.Instance{
//									{
//										LaunchTime:       aws.Time(time.Now()),
//										InstanceId:       aws.String("ID"),
//										PublicIpAddress:  aws.String("127.0.0.1"),
//										PrivateIpAddress: aws.String("127.0.0.1"),
//									},
//								},
//							},
//						},
//					},
//				}, nil
//			},
//			importErr:         nil,
//			hasPublicIP: true,
//		},
//		{
//			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
//				return &fakeEC2CreateMachine{
//					runInstancesOutput: &ec2.Reservation{
//						Instances: []*ec2.Instance{
//							{
//								LaunchTime: aws.Time(time.Now()),
//								InstanceId: aws.String("ID"),
//							},
//						},
//					},
//					descInstancesOutput: &ec2.DescribeInstancesOutput{
//						Reservations: []*ec2.Reservation{
//							{
//								Instances: []*ec2.Instance{
//									{
//										LaunchTime: aws.Time(time.Now()),
//										InstanceId: aws.String("ID"),
//									},
//								},
//							},
//						},
//					},
//				}, nil
//			},
//			importErr:         ErrNoPublicIP,
//			hasPublicIP: true,
//		},
//	}
//
//	for i, tc := range tt {
//		cfg := steps.NewConfig("test", "", "", profile.Profile{
//			NodesProfiles: []profile.NodeProfile{
//				map[string]string{"test": "test"},
//				map[string]string{"test": "test"},
//				map[string]string{"test": "test"},
//			},
//		})
//		cfg.TaskID = "task"
//		if tc.hasPublicIP {
//			cfg.AWSConfig.HasPublicAddr = "true"
//		}
//
//		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
//		defer cancel()
//
//		step := NewCreateInstance(tc.fn)
//		importErr := step.Run(ctx, os.Stdout, cfg)
//		if tc.importErr == nil {
//			require.NoError(t, importErr, "TC%d, %v", i, importErr)
//		} else {
//			require.Truef(t, tc.importErr == errors.Cause(importErr), "TC%d, %v", i, importErr)
//		}
//	}
//}

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

	if step.GetEC2 == nil {
		t.Errorf("Wrong GetEC2 function must not be nil")
	}

	if api, err := step.GetEC2(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateInstanceError(t *testing.T) {
	fn := func(steps.AWSConfig)(ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewCreateInstance(fn)

	if step.GetEC2 == nil {
		t.Errorf("Wrong GetEC2 function must not be nil")
	}

	if api, err := step.GetEC2(steps.AWSConfig{}); err == nil || api != nil {
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
		t.Errorf("Wrong description value expected " +
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
