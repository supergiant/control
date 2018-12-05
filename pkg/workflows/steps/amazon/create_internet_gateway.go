package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateInternetGateway = "create_internet_gateway"

type CreateInternetGatewayStep struct {
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitCreateInternetGateway(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateInternetGateway, NewCreateInternetGatewayStep(ec2fn))
}

func NewCreateInternetGatewayStep(ec2fn GetEC2Fn) *CreateInternetGatewayStep {
	return &CreateInternetGatewayStep{
		GetEC2: ec2fn,
	}
}

func (s *CreateInternetGatewayStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateInternetGateway)
	ec2Client, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	// Internet gateway already exists
	if cfg.AWSConfig.InternetGatewayID != "" {
		return nil
	} else {
		// Use default gateway for VPC
		output, err := ec2Client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{aws.String(cfg.AWSConfig.VPCID)},
				},
			},
		})

		if err != nil {
			logrus.Errorf("error getting internet gateway for vpc %s %v",
				cfg.AWSConfig.VPCID, err)
			return err
		}

		if len(output.InternetGateways) == 0 {
			return errors.Wrapf(sgerrors.ErrNotFound,
				"not found gateways for vpc id %s",
				cfg.AWSConfig.VPCID)
		}

		cfg.AWSConfig.InternetGatewayID = *output.InternetGateways[0].InternetGatewayId
	}

	return nil
}

func (s *CreateInternetGatewayStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*CreateInternetGatewayStep) Name() string {
	return StepNameCreateEC2Instance
}

func (*CreateInternetGatewayStep) Description() string {
	return "Create internet gateway"
}

func (*CreateInternetGatewayStep) Depends() []string {
	return nil
}
