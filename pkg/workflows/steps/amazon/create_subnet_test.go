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

	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/account"
)

type fakeEC2Subnet struct {
	ec2iface.EC2API
	output *ec2.CreateSubnetOutput
	err    error
}

func (f *fakeEC2Subnet) CreateSubnetWithContext(aws.Context, *ec2.CreateSubnetInput, ...request.Option) (*ec2.CreateSubnetOutput, error) {
	return f.output, f.err
}

type mockZoneGetter struct{
	zones []string
	err error
}

func (m *mockZoneGetter) GetZones(context.Context, steps.Config) ([]string, error) {
	return m.zones, m.err
}

func TestCreateSubnetStep_Run(t *testing.T) {
	tt := []struct {
		fn  GetEC2Fn
		err error
		cfg steps.AWSConfig
	}{
		{
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2Subnet{
					output: &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:            aws.String("1"),
							AvailabilityZone: aws.String("my-az"),
							SubnetId:         aws.String("mysubnetid"),
						},
					},
				}, nil
			},
			err: nil,
			cfg: steps.AWSConfig{
				VPCCIDR: "10.0.0.0/16",
			},
		},
		{
			cfg: steps.AWSConfig{
				VPCCIDR: "10.0.0.0/16",
			},
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &fakeEC2Subnet{
					output: nil,
					err:    errors.New("fail!"),
				}, nil
			},
			err: ErrCreateSubnet,
		},
	}

	for i, tc := range tt {
		cfg := steps.NewConfig("", "", "", profile.Profile{})
		cfg.AWSConfig = tc.cfg

		step := &CreateSubnetsStep{
			GetEC2: tc.fn,
			accSvc: nil,
			zoneGetterFactory: func(ctx context.Context, accSvc *account.Service,
				cfg *steps.Config) (account.ZonesGetter, error){
				return &mockZoneGetter{
					zones: []string{"eu-west-1a", "eu-west-1b"},
					err: nil,
				}, nil
			},
		}
		err := step.Run(context.Background(), os.Stdout, cfg)
		if tc.err == nil {
			require.NoError(t, err, "TC%d, %v", i, err)
		} else {
			require.True(t, tc.err == errors.Cause(err), "TC%d, %v", i, err)
		}
	}
}
