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
	"bytes"
)

type FakeEC2SecurityGroups struct {
	ec2iface.EC2API
	createOutput   *ec2.CreateSecurityGroupOutput
	describeOutput *ec2.DescribeSecurityGroupsOutput
	err            error
}

func (f *FakeEC2SecurityGroups) CreateSecurityGroupWithContext(aws.Context, *ec2.CreateSecurityGroupInput, ...request.Option) (*ec2.CreateSecurityGroupOutput, error) {
	return f.createOutput, f.err
}

func (f *FakeEC2SecurityGroups) AuthorizeSecurityGroupIngressWithContext(aws.Context, *ec2.AuthorizeSecurityGroupIngressInput, ...request.Option) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return nil, nil
}

func (f *FakeEC2SecurityGroups) DescribeSecurityGroupsWithContext(aws.Context, *ec2.DescribeSecurityGroupsInput, ...request.Option) (*ec2.DescribeSecurityGroupsOutput, error) {
	return f.describeOutput, f.err
}

func TestCreateSecurityGroupsStep_Run(t *testing.T) {
	tt := []struct {
		description string
		fn          GetEC2Fn
		err         error
		cfg         steps.AWSConfig
	}{
		{
			description: "error authorization",
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return nil, nil
			},

			err: ErrAuthorization,
			cfg: steps.AWSConfig{
				VPCID: "",
			},
		},
		{
			description: "success",
			fn: func(config steps.AWSConfig) (ec2iface.EC2API, error) {
				return &FakeEC2SecurityGroups{
					createOutput: &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("MYID"),
					},
					describeOutput: &ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{
								GroupId:   aws.String("MYID"),
								GroupName: aws.String("GROUPNAME"),
							},
						},
					},
				}, nil
			},
			cfg: steps.AWSConfig{
				VPCID:                  "ID",
				NodesSecurityGroupID:   "",
				MastersSecurityGroupID: "",
			},
		},
	}

	for i, tc := range tt {
		t.Log(tc.description)
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

func TestInitCreateSecurityGroups(t *testing.T) {
	InitCreateSecurityGroups(GetEC2)

	s := steps.GetStep(StepCreateSecurityGroups)

	if s == nil {
		t.Errorf("Step must not be nil")
	}
}

func TestNewCreateSecurityGroupsStep(t *testing.T) {
	s := NewCreateSecurityGroupsStep(GetEC2)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.GetEC2 == nil {
		t.Errorf("GetEC2 must not be nil")
	}

	if api, err := s.GetEC2(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateSecurityGroupsStepErr(t *testing.T) {
	fn := func(steps.AWSConfig)(ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewCreateSecurityGroupsStep(fn)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.GetEC2 == nil {
		t.Errorf("GetEC2 must not be nil")
	}

	if api, err := s.GetEC2(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestCreateSecurityGroupsStep_Depends(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Deps must be nil")
	}
}

func TestCreateSecurityGroupsStep_Name(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if name := s.Name(); name != StepCreateSecurityGroups {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateSecurityGroups, s.Name())
	}
}

func TestCreateSecurityGroupsStep_Rollback(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected value of err %v", err)
	}
}

func TestCreateSecurityGroupsStep_Description(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if desc := s.Description(); desc != "Create security groups" {
		t.Errorf("Wrong step desc expected Create security groups" +
			" actual %s", s.Description())
	}
}