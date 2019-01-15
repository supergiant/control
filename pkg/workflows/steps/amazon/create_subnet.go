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
	"github.com/supergiant/control/pkg/model"
	"github.com/aws/aws-sdk-go/aws/request"
)

const StepCreateSubnets = "create_subnet_steps"

type accountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type subnetSvc interface {
	CreateSubnetWithContext(aws.Context, *ec2.CreateSubnetInput,
		...request.Option) (*ec2.CreateSubnetOutput, error)
}

type CreateSubnetsStep struct {
	accountGetter accountGetter
	getSvc func(steps.AWSConfig) (subnetSvc, error)
	zoneGetterFactory func(context.Context, accountGetter, *steps.Config) (account.ZonesGetter, error)
}

func NewCreateSubnetStep(fn GetEC2Fn, getter accountGetter) *CreateSubnetsStep {
	return &CreateSubnetsStep{
		accountGetter: getter,
		getSvc: func(config steps.AWSConfig) (subnetSvc, error) {
			client, err := fn(config)

			if err != nil {
				return nil, ErrAuthorization
			}

			return client, nil
		},
		zoneGetterFactory: func(ctx context.Context, accountGetter accountGetter,
			cfg *steps.Config) (account.ZonesGetter, error) {
			acc, err := accountGetter.Get(ctx, cfg.CloudAccountName)

			if err != nil {
				logrus.Errorf("Get cloud account %s caused error %v",
					cfg.CloudAccountName, err)
				return nil, errors.Wrapf(err, "Get cloud account")
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
	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("%s error getting service %v",
			StepCreateSubnets, err)
		return errors.Wrapf(err, "%s error getting service",
			StepCreateSubnets)
	}

	zoneGetter, err := s.zoneGetterFactory(ctx, s.accountGetter, cfg)

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
		_, cidrIP, err := net.ParseCIDR(cfg.AWSConfig.VPCCIDR)

		if err != nil {
			logrus.Errorf("Error parsing VPC cidr %s",
				cfg.AWSConfig.VPCCIDR)
			return errors.Wrapf(err, "Error parsing VPC cidr %s",
				cfg.AWSConfig.VPCCIDR)
		}

		logrus.Info(cidrIP)
		subnetCidr, err := cidr.Subnet(cidrIP, 8, rand.Int()%256)
		logrus.Debugf("Subnet cidr %s", subnetCidr)

		if err != nil {
			logrus.Debugf("Calculating subnet cidr caused %s", err.Error())
			return errors.Wrapf(err, "%s Calculating subnet" +
				" cidr caused error", StepCreateSubnets)
		}

		input := &ec2.CreateSubnetInput{
			VpcId:            aws.String(cfg.AWSConfig.VPCID),
			AvailabilityZone: aws.String(zone),
			CidrBlock:        aws.String(subnetCidr.String()),
		}
		out, err := svc.CreateSubnetWithContext(ctx, input)
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
