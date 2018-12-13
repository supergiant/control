package amazon

import (
	"context"
	"io"
	"math/rand"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateSubnets = "create_subnet_steps"

type CreateSubnetsStep struct {
	GetEC2 GetEC2Fn
	accSvc *account.Service

	zoneGetterFactory func(context.Context, *account.Service, *steps.Config) (account.ZonesGetter, error)
}

func NewCreateSubnetStep(fn GetEC2Fn, accSvc *account.Service) *CreateSubnetsStep {
	return &CreateSubnetsStep{
		GetEC2: fn,
		accSvc: accSvc,

		zoneGetterFactory: func(ctx context.Context, accSvc *account.Service, cfg *steps.Config) (account.ZonesGetter, error) {
			acc, err := accSvc.Get(ctx, cfg.CloudAccountName)

			if err != nil {
				logrus.Errorf("Get cloud account %s caused error %v",
					cfg.CloudAccountName, err)
				return nil, err
			}

			zoneGetter, err := account.NewZonesGetter(acc, cfg)

			return zoneGetter, err
		},
	}
}

func InitCreateSubnet(fn GetEC2Fn, accSvc *account.Service) {
	steps.RegisterStep(StepCreateSubnets, NewCreateSubnetStep(fn, accSvc))
}

func (s *CreateSubnetsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	zoneGetter, err := s.zoneGetterFactory(ctx, s.accSvc, cfg)

	if err != nil {
		logrus.Errorf("Create zone getter caused error %v", err)
		return errors.Wrapf(err, "create subnets for vpc %s", cfg.AWSConfig.VPCID)
	}

	// Make sure we create a subnet map
	if cfg.AWSConfig.Subnets == nil {
		cfg.AWSConfig.Subnets = make(map[string]string)
	}

	logrus.Debugf(cfg.AWSConfig.VPCCIDR)
	logrus.Debugf("Create subnet in VPC %s", cfg.AWSConfig.VPCID)

	logrus.Debugf("get zones for region %s", cfg.AWSConfig.Region)
	zones, err := zoneGetter.GetZones(ctx, *cfg)

	if err != nil {
		logrus.Errorf("Error getting zones for region %s",
			cfg.AWSConfig.Region)
		return errors.Wrapf(err, "Error getting zone for region %s",
			cfg.AWSConfig.Region)
	}

	// Create subnet for each availability zone
	for _, zone := range zones {

		// We already have existing subnet for that
		if cfg.AWSConfig.Subnets[zone] != "" {
			continue
		}

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
			logrus.Debugf("Create subnet cause error %s", err.Error())
			return errors.Wrap(ErrCreateSubnet, err.Error())
		}

		// Store subnet in subnets map
		cfg.AWSConfig.Subnets[zone] = *out.Subnet.SubnetId
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
