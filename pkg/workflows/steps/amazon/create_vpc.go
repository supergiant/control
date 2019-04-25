package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateVPC = "aws_create_vpc"

//CreateVPCStep represents creation of an virtual private cloud in AWS
type CreateVPCStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateVPCStep(fn GetEC2Fn) *CreateVPCStep {
	return &CreateVPCStep{
		GetEC2: fn,
	}
}

func InitCreateVPC(fn GetEC2Fn) {
	steps.RegisterStep(StepCreateVPC, NewCreateVPCStep(fn))
}

func (c *CreateVPCStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := c.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	//A user doesn't specified that she wants to use preexisting VPC
	//creating a new one for a cluster
	if cfg.AWSConfig.VPCID == "" {
		log.Infof("[%s] - no VPC id specified, creating now...", c.Name())

		input := &ec2.CreateVpcInput{
			CidrBlock: &cfg.AWSConfig.VPCCIDR,
		}
		out, err := EC2.CreateVpcWithContext(ctx, input)
		if err != nil {
			return errors.Wrap(ErrCreateVPC, err.Error())
		}
		cfg.AWSConfig.VPCID = *out.Vpc.VpcId

		vpcattr := &ec2.ModifyVpcAttributeInput{
			EnableDnsHostnames: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
			VpcId: &cfg.AWSConfig.VPCID,
		}
		_, err = EC2.ModifyVpcAttributeWithContext(ctx, vpcattr)
		if err != nil {
			return errors.Wrap(ErrCreateVPC, err.Error())
		}

		desc := &ec2.DescribeVpcsInput{
			VpcIds: []*string{aws.String(cfg.AWSConfig.VPCID)},
		}
		err = EC2.WaitUntilVpcExistsWithContext(ctx, desc)

		if err != nil {
			logrus.Debugf("error waiting for vpc %s %s",
				cfg.AWSConfig.VPCID, err.Error())
			return errors.Wrapf(err, "create vpc error wait")
		}
		log.Infof("[%s] - created a VPC with ID %s and CIDR %s",
			c.Name(), cfg.AWSConfig.VPCID, cfg.AWSConfig.VPCCIDR)
	} else {
		out, err := EC2.DescribeVpcsWithContext(ctx, &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("isDefault"),
					Values: aws.StringSlice([]string{
						"true",
					}),
				},
			},
		})
		if err != nil {
			log.Errorf("[%s] - failed to read VPC data", c.Name())
			return errors.Wrap(ErrReadVPC, err.Error())
		}

		var defaultVPCID string
		var defaultVPCCIDR string
		for _, vpc := range out.Vpcs {
			if *vpc.IsDefault {
				defaultVPCID = *vpc.VpcId
				defaultVPCCIDR = *vpc.CidrBlock
				break
			}
		}

		cfg.AWSConfig.VPCID = defaultVPCID
		cfg.AWSConfig.VPCCIDR = defaultVPCCIDR
	}

	return nil
}

func (*CreateVPCStep) Name() string {
	return StepCreateVPC
}

func (*CreateVPCStep) Description() string {
	return "create a vpc in aws or reuse existing one"
}

func (*CreateVPCStep) Depends() []string {
	return nil
}

func (*CreateVPCStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
