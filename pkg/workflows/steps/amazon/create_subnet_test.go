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

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"bytes"
)

type fakeEC2Subnet struct {
	ec2iface.EC2API
	output *ec2.CreateSubnetOutput
	err    error
}

func (f *fakeEC2Subnet) CreateSubnetWithContext(aws.Context, *ec2.CreateSubnetInput, ...request.Option) (*ec2.CreateSubnetOutput, error) {
	return f.output, f.err
}

type mockZoneGetter struct {
	zones []string
	err   error
}

func (m *mockZoneGetter) GetZones(context.Context, steps.Config) ([]string, error) {
	return m.zones, m.err
}

func TestCreateSubnetStep_Run(t *testing.T) {
	tt := []struct {
		description   string
		fn            GetEC2Fn
		err           error
		createZoneErr error
		getZoneErr    error
		cfg           steps.AWSConfig
	}{
		{
			description: "success",
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
			description: "err creating subnet",
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
		{
			description: "err create zone getter",
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
			cfg: steps.AWSConfig{
				VPCCIDR: "10.0.0.0/16",
			},
			createZoneErr: sgerrors.ErrInvalidCredentials,
			err:           sgerrors.ErrInvalidCredentials,
		},
		{
			description: "err getting zone",
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
			cfg: steps.AWSConfig{
				VPCCIDR: "10.0.0.0/16",
			},
			getZoneErr: sgerrors.ErrNotFound,
			err:        sgerrors.ErrNotFound,
		},
	}

	for i, tc := range tt {
		t.Log(tc.description)
		cfg := steps.NewConfig("", "", "", profile.Profile{})
		cfg.AWSConfig = tc.cfg

		step := &CreateSubnetsStep{
			GetEC2: tc.fn,
			accSvc: nil,
			zoneGetterFactory: func(ctx context.Context, accSvc *account.Service,
				cfg *steps.Config) (account.ZonesGetter, error) {
				return &mockZoneGetter{
					zones: []string{"eu-west-1a", "eu-west-1b"},
					err:   tc.getZoneErr,
				}, tc.createZoneErr
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

func TestInitCreateSubnet(t *testing.T) {
	InitCreateSubnet(GetEC2, nil)

	s := steps.GetStep(StepCreateSubnets)

	if s == nil {
		t.Errorf("Step %s not found", StepCreateSubnets)
	}
}

func TestNewCreateSubnetStep(t *testing.T) {
	accSvc := &account.Service{}

	step := NewCreateSubnetStep(GetEC2, accSvc)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.accSvc != accSvc {
		t.Errorf("account service value is wrong exepected %v actual %v",
			step.accSvc, accSvc)
	}

	if step.GetEC2 == nil {
		t.Errorf("Wrong get EC2 function must not be nil")
	}

	if api, err := step.GetEC2(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}


func TestNewCreateSubnetStepErr(t *testing.T) {
	fn := func(steps.AWSConfig)(ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	accSvc := &account.Service{}

	step := NewCreateSubnetStep(fn, accSvc)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.accSvc != accSvc {
		t.Errorf("account service value is wrong exepected %v actual %v",
			step.accSvc, accSvc)
	}

	if step.GetEC2 == nil {
		t.Errorf("Wrong get EC2 function must not be nil")
	}

	if api, err := step.GetEC2(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestCreateSubnetsStep_Name(t *testing.T) {
	step := &CreateSubnetsStep{}

	if step.Name() != StepCreateSubnets {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateSubnets, step.Name())
	}
}


func TestCreateSubnetsStep_Description(t *testing.T) {
	step := &CreateSubnetsStep{}

	if step.Description() != "Step create subnets in all availability zones for Region" {
		t.Errorf("Wrong step description expected Step create subnets in " +
			"all availability zones for Region actual %s", step.Description())
	}
}

func TestCreateSubnetsStep_Depends(t *testing.T) {
	step := &CreateSubnetsStep{}

	if deps := step.Depends(); deps != nil {
		t.Errorf("%s dependencies must be nil", StepCreateSubnets)
	}
}

func TestCreateSubnetsStep_Rollback(t *testing.T) {
	step := &CreateSubnetsStep{}

	if err := step.Rollback(context.Background(),
		&bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v when rollback", err)
	}
}