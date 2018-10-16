package amazon

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type fakeEC2CreateMachine struct {
	ec2iface.EC2API

	runInstancesOutput  *ec2.Reservation
	descInstancesOutput *ec2.DescribeInstancesOutput
	err                 error
}

func (f *fakeEC2CreateMachine) RunInstancesWithContext(aws.Context, *ec2.RunInstancesInput, ...request.Option) (*ec2.Reservation, error) {
	return f.runInstancesOutput, f.err
}

func (f *fakeEC2CreateMachine) DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return f.descInstancesOutput, f.err
}

func TestStepCreateInstance_Run(t *testing.T) {
	tt := []struct {
		fn          GetEC2Fn
		err         error
		hasPublicIP bool
	}{
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2CreateMachine{
					runInstancesOutput: &ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								LaunchTime:      aws.Time(time.Now()),
								InstanceId:      aws.String("ID"),
								PublicIpAddress: aws.String("127.0.0.1"),
							},
						},
					},
				}, nil
			},
			err:         nil,
			hasPublicIP: false,
		},
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2CreateMachine{
					runInstancesOutput: &ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								LaunchTime:      aws.Time(time.Now()),
								InstanceId:      aws.String("ID"),
								PublicIpAddress: aws.String("127.0.0.1"),
							},
						},
					},
					descInstancesOutput: &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										LaunchTime:       aws.Time(time.Now()),
										InstanceId:       aws.String("ID"),
										PublicIpAddress:  aws.String("127.0.0.1"),
										PrivateIpAddress: aws.String("127.0.0.1"),
									},
								},
							},
						},
					},
				}, nil
			},
			err:         nil,
			hasPublicIP: true,
		},
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2CreateMachine{
					runInstancesOutput: &ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								LaunchTime: aws.Time(time.Now()),
								InstanceId: aws.String("ID"),
							},
						},
					},
					descInstancesOutput: &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										LaunchTime: aws.Time(time.Now()),
										InstanceId: aws.String("ID"),
									},
								},
							},
						},
					},
				}, nil
			},
			err:         ErrNoPublicIP,
			hasPublicIP: true,
		},
	}

	for i, tc := range tt {
		cfg := steps.NewConfig("test", "", "", profile.Profile{
			NodesProfiles: []profile.NodeProfile{
				map[string]string{"test": "test"},
				map[string]string{"test": "test"},
				map[string]string{"test": "test"},
			},
		})
		cfg.TaskID = "task"
		if tc.hasPublicIP {
			cfg.AWSConfig.EC2Config.HasPublicAddr = true
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		step := NewCreateInstance(tc.fn)
		err := step.Run(ctx, os.Stdout, cfg)
		if tc.err == nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.Truef(t, tc.err == errors.Cause(err), "TC%d, %v", i, err)
		}
	}
}
