package amazon

import (
	"context"
	"io"
	"math/rand"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
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
		logrus.Debugf(cfg.AWSConfig.VPCCIDR)
		logrus.Debugf("Create subnet in VPC %s", cfg.AWSConfig.VPCID)
		_, cidrIP, _ := net.ParseCIDR(cfg.AWSConfig.VPCCIDR)

		subnetCidr, err := cidr.Subnet(cidrIP, 8, rand.Int()%256)
		logrus.Debugf("Subnet cidr %s", subnetCidr)

		if err != nil {
			logrus.Debugf("Calculating subnet cidr caused %s", err.Error())
		}

		input := &ec2.CreateSubnetInput{
			VpcId:            aws.String(cfg.AWSConfig.VPCID),
			AvailabilityZone: aws.String(cfg.AWSConfig.AvailabilityZone),
			CidrBlock:        aws.String(subnetCidr.String()),
		}
		out, err := EC2.CreateSubnetWithContext(ctx, input)
		if err != nil {
			if err, ok := err.(awserr.Error); ok {
				logrus.Debugf("Create subnet cause error %s", err.Message())
			}
			return errors.Wrap(ErrCreateSubnet, err.Error())
		}

		cfg.AWSConfig.SubnetID = *out.Subnet.SubnetId

		return nil
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
		logrus.Debugf("Take subnet %s", *out.Subnets[0].SubnetId)
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
