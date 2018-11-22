package amazon

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type FakeEC2DeleteCluster struct {
	ec2iface.EC2API

	instancesOutput  *ec2.DescribeInstancesOutput
	terminatedOutput *ec2.TerminateInstancesOutput
	err              error
}

func (f *FakeEC2DeleteCluster) DescribeInstancesWithContext(aws.Context, *ec2.DescribeInstancesInput, ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return f.instancesOutput, f.err
}

func (f *FakeEC2DeleteCluster) TerminateInstancesWithContext(aws.Context, *ec2.TerminateInstancesInput, ...request.Option) (*ec2.TerminateInstancesOutput, error) {
	return f.terminatedOutput, f.err
}

func TestDeleteClusterStep_Run(t *testing.T) {
	tt := []struct {
		fn     GetEC2Fn
		err    error
		awsCfg steps.AWSConfig
	}{
		//No instances found
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &FakeEC2DeleteCluster{
					instancesOutput: &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					},
				}, nil
			},
			err:    nil,
			awsCfg: steps.AWSConfig{},
		},
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &FakeEC2DeleteCluster{
					instancesOutput: &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										InstanceId: aws.String("test"),
									},
								},
							},
						},
					},
					terminatedOutput: &ec2.TerminateInstancesOutput{
						TerminatingInstances: []*ec2.InstanceStateChange{
							{
								InstanceId: aws.String("test"),
							},
						},
					},
				}, nil
			},
			err:    nil,
			awsCfg: steps.AWSConfig{},
		},
	}

	for i, tc := range tt {
		cfg := steps.NewConfig("TEST", "", "TEST", profile.Profile{
			Region:   "us-east-1",
			Provider: clouds.AWS,
		})
		cfg.AWSConfig.Region = "us-east-1"
		cfg.AWSConfig = tc.awsCfg

		step := DeleteClusterStep{
			GetEC2: tc.fn,
		}

		err := step.Run(context.Background(), os.Stdout, cfg)
		if tc.err == nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.True(t, tc.err == errors.Cause(err), "TC%d, %v", i, err)
		}

	}
}
