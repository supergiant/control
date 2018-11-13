package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepCreateSubnet = "create_subnet_step"

type CreateSubnetStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateSubnetStep(fn GetEC2Fn) *CreateSubnetStep {
	return &CreateSubnetStep{
		GetEC2: fn,
	}
}

func InitCreateSubnet(fn GetEC2Fn) {
	steps.RegisterStep(StepCreateSubnet, NewCreateSubnetStep(fn))
}

func (s *CreateSubnetStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}
	if cfg.AWSConfig.SubnetID == "" {
		input := &ec2.CreateSubnetInput{
			VpcId:            aws.String(cfg.AWSConfig.VPCID),
			AvailabilityZone: aws.String(cfg.AWSConfig.AvailabilityZone),
		}
		out, err := EC2.CreateSubnetWithContext(ctx, input)
		if err != nil {
			return errors.Wrap(ErrCreateSubnet, err.Error())
		}
		cfg.AWSConfig.SubnetID = *out.Subnet.SubnetId
	} else if cfg.AWSConfig.SubnetID == "default" {
		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: aws.StringSlice([]string{cfg.AWSConfig.VPCID}),
				},
				{
					Name:   aws.String("availabilityZone"),
					Values: aws.StringSlice([]string{cfg.AWSConfig.AvailabilityZone}),
				},
				{
					Name:   aws.String("default-for-az"),
					Values: aws.StringSlice([]string{"true"}),
				},
			},
		}
		out, err := EC2.DescribeSubnetsWithContext(ctx, input)
		if err != nil {
			return errors.Wrap(ErrCreateSubnet, err.Error())
		}
		if len(out.Subnets) == 0 {
			return errors.Wrap(ErrCreateSubnet, "no default subnet found")
		}
		cfg.AWSConfig.SubnetID = *out.Subnets[0].SubnetId
	} else {
		log.Infof("[%s] - using subnet %s", s.Name(), cfg.AWSConfig.SubnetID)
	}
	return nil
}

func (*CreateSubnetStep) Name() string {
	return StepCreateSubnet
}

func (*CreateSubnetStep) Description() string {
	return ""
}

func (*CreateSubnetStep) Depends() []string {
	return nil
}

func (*CreateSubnetStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
