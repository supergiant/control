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

//CreateVPCStep represents creation of an virtual private cloud in AWS
type CreateVPCStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateVPCStep(fn GetEC2Fn) *CreateVPCStep {
	return &CreateVPCStep{
		GetEC2: fn,
	}
}

func (c *CreateVPCStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(err, "aws: authorization")
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
			return errors.Wrap(err, "aws: create vpc")
		}
		cfg.AWSConfig.VPCID = *out.Vpc.VpcId

		log.Infof("[%s] - created a VPC with ID %s and CIDR %s",
			c.Name(),
			cfg.AWSConfig.VPCID,
			cfg.AWSConfig.VPCCIDR)
	} else {
		if cfg.AWSConfig.VPCID != "default" {
			//if a user specified that there is a vpc already exists it should be verified
			out, err := EC2.DescribeVpcsWithContext(ctx, &ec2.DescribeVpcsInput{
				VpcIds: aws.StringSlice([]string{cfg.AWSConfig.VPCID}),
			})
			if err != nil {
				log.Errorf("[%s] - failed to read VPC data", c.Name())
				return errors.Wrap(err, "aws: read vpc")
			}
			if len(out.Vpcs) == 0 {
				return errors.Wrap(err, "aws: read vpc")
			}
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
				return errors.Wrap(err, "aws: read vpc")
			}
			if len(out.Vpcs) == 0 {
				return errors.Wrap(err, "aws: read vpc")
			}

			var defaultVPCID string
			for _, vpc := range out.Vpcs {
				if *vpc.IsDefault {
					defaultVPCID = *vpc.VpcId
					break
				}
			}

			if defaultVPCID == "" {
				return errors.Wrap(err, "aws: read vpc")
			}

			cfg.AWSConfig.VPCID = defaultVPCID
		}
	}

	return nil
}

func (*CreateVPCStep) Name() string {
	return "aws_create_vpc"
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
