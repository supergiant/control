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
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type FakeEC2SecurityGroups struct {
	ec2iface.EC2API
	output *ec2.CreateSecurityGroupOutput
	err    error
}

func (f *FakeEC2SecurityGroups) CreateSecurityGroupWithContext(aws.Context, *ec2.CreateSecurityGroupInput, ...request.Option) (*ec2.CreateSecurityGroupOutput, error) {
	return f.output, f.err
}

func TestCreateSecurityGroupsStep_Run(t *testing.T) {
	tt := []struct {
		fn  GetEC2Fn
		err error
		cfg steps.AWSConfig
	}{
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return nil, nil
			},
			err: ErrAuthorization,
			cfg: steps.AWSConfig{
				VPCID: "",
			},
		},
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &FakeEC2SecurityGroups{
					output: &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("MYID"),
					},
				}, nil
			},
			cfg: steps.AWSConfig{
				VPCID:                "ID",
				NodesSecurityGroup:   "",
				MastersSecurityGroup: "",
			},
		},
	}

	for i, tc := range tt {
		cfg := steps.NewConfig("", "", "", profile.Profile{})
		cfg.AWSConfig = tc.cfg
		step := NewCreateSecurityGroupsStep(tc.fn)
		err := step.Run(context.Background(), os.Stdout, cfg)
		if tc.err == nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.True(t, tc.err == errors.Cause(err), "TC%d, %v", i, err)
		}
	}
}
