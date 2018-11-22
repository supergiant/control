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

type fakeEC2VPC struct {
	ec2iface.EC2API
	createVPCOutput   *ec2.CreateVpcOutput
	describeVPCOutput *ec2.DescribeVpcsOutput
	err               error
}

func (f *fakeEC2VPC) CreateVpcWithContext(aws.Context, *ec2.CreateVpcInput, ...request.Option) (*ec2.CreateVpcOutput, error) {
	return f.createVPCOutput, f.err
}

func (f *fakeEC2VPC) DescribeVpcsWithContext(aws.Context, *ec2.DescribeVpcsInput, ...request.Option) (*ec2.DescribeVpcsOutput, error) {
	return f.describeVPCOutput, f.err
}

func TestCreateVPCStep_Run(t *testing.T) {
	tt := []struct {
		awsFN  GetEC2Fn
		err    error
		awsCfg steps.AWSConfig
	}{
		{
			//happy path
			func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2VPC{
					createVPCOutput: &ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							VpcId: aws.String("ID"),
						},
					},
				}, nil
			},
			nil,
			steps.AWSConfig{},
		},
		{
			//error, vpc id was provided but isn't available in the AWS
			func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2VPC{
					err: errors.New(""),
				}, nil
			},
			ErrReadVPC,
			steps.AWSConfig{
				VPCID: "1",
			},
		},
		{
			func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2VPC{
					describeVPCOutput: &ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							{
								VpcId: aws.String("1"),
							},
						},
					},
				}, nil
			},
			nil,
			steps.AWSConfig{
				VPCID: "1",
			},
		},
		{
			func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2VPC{
					describeVPCOutput: &ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							{
								VpcId:     aws.String("default"),
								IsDefault: aws.Bool(true),
							},
						},
					},
				}, nil
			},
			nil,
			steps.AWSConfig{
				VPCID: "default",
			},
		},
	}

	for i, tc := range tt {
		cfg := steps.NewConfig("TEST", "", "TEST", profile.Profile{
			Region:   "us-east-1",
			Provider: clouds.AWS,
		})
		cfg.AWSConfig = tc.awsCfg

		step := NewCreateVPCStep(tc.awsFN)
		err := step.Run(context.Background(), os.Stdout, cfg)
		if tc.err == nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.True(t, tc.err == errors.Cause(err), "TC%d, %v", i, err)
		}
	}

}
