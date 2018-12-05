package amazon

import (
	"context"
	"io"
	"math/rand"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/account"
)

const StepCreateSubnets = "create_subnet_steps"

type CreateSubnetsStep struct {
	GetEC2 GetEC2Fn
	zoneGetter account.ZonesGetter
}

func NewCreateSubnetStep(fn GetEC2Fn, zoneGetter account.ZonesGetter) *CreateSubnetsStep {
	return &CreateSubnetsStep{
		GetEC2: fn,
		zoneGetter: zoneGetter,
	}
}

func InitCreateSubnet(fn GetEC2Fn, zoneGetter account.ZonesGetter) {
	steps.RegisterStep(StepCreateSubnets, NewCreateSubnetStep(fn, zoneGetter))
}

func (s *CreateSubnetsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	if len(cfg.AWSConfig.Subnets) == 0 {
		logrus.Debugf(cfg.AWSConfig.VPCCIDR)
		logrus.Debugf("Create subnet in VPC %s", cfg.AWSConfig.VPCID)

		logrus.Debugf("get zones for region %s", cfg.AWSConfig.Region)
		zones, err := s.zoneGetter.GetZones(ctx, *cfg)

		if err != nil {
			logrus.Errorf("Error getting zones for region %s",
				cfg.AWSConfig.Region)
		}

		// Make sure we create a subnet map
		if cfg.AWSConfig.Subnets == nil {
			cfg.AWSConfig.Subnets = make(map[string]string)
		}

		// Create subnet for each availability zone
		for _, zone := range zones {
			_, cidrIP, _ := net.ParseCIDR(cfg.AWSConfig.VPCCIDR)

			subnetCidr, err := cidr.Subnet(cidrIP, 8, rand.Int()%256)
			logrus.Debugf("Subnet cidr %s", subnetCidr)

			if err != nil {
				logrus.Debugf("Calculating subnet cidr caused %s", err.Error())
			}

			input := &ec2.CreateSubnetInput{
				VpcId:            aws.String(cfg.AWSConfig.VPCID),
				AvailabilityZone: aws.String(zone),
				CidrBlock:        aws.String(subnetCidr.String()),
			}
			out, err := EC2.CreateSubnetWithContext(ctx, input)
			if err != nil {
				if err, ok := err.(awserr.Error); ok {
					logrus.Debugf("Create subnet cause error %s", err.Message())
				}
				return errors.Wrap(ErrCreateSubnet, err.Error())
			}

			// Store subnet in subnets map
			cfg.AWSConfig.Subnets[zone] = *out.Subnet.SubnetId
		}

		return nil
	} else {
		log.Infof("[%s] - using subnets %s", s.Name(), cfg.AWSConfig.Subnets)
	}

	return nil
}

func (*CreateSubnetsStep) Name() string {
	return StepCreateSubnets
}

func (*CreateSubnetsStep) Description() string {
	return "Step create subnets in all availability zones for Region"
}

func (*CreateSubnetsStep) Depends() []string {
	return nil
}

func (*CreateSubnetsStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
