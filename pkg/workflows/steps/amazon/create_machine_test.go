package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
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
