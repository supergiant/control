package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	ImportInternetGatewayStepName = "import_internet_gateway_aws"
)

type GatewayDescriber interface {
	DescribeInternetGateways(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
}

type ImportInternetGatewayStep struct {
	getSvc func(config steps.AWSConfig) (GatewayDescriber, error)
}

func NewImportInternetGatewayStep(fn GetEC2Fn) *ImportInternetGatewayStep {
	return &ImportInternetGatewayStep{
		getSvc: func(config steps.AWSConfig) (describer GatewayDescriber, e error) {
			EC2, err := fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func InitImportInternetGatewayStep(fn GetEC2Fn) {
	steps.RegisterStep(ImportInternetGatewayStepName, NewImportInternetGatewayStep(fn))
}

func (s ImportInternetGatewayStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	ec2Svc, err := s.getSvc(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	req := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []*string{aws.String(cfg.AWSConfig.VPCID)},
			},
		},
	}
	output, err := ec2Svc.DescribeInternetGateways(req)

	if err != nil {
		return errors.Wrapf(err, "")
	}

	cfg.AWSConfig.InternetGatewayID = *output.InternetGateways[0].InternetGatewayId

	return nil
}

func (s ImportInternetGatewayStep) Name() string {
	return ImportInternetGatewayStepName
}

func (s ImportInternetGatewayStep) Description() string {
	return ImportInternetGatewayStepName
}

func (s ImportInternetGatewayStep) Depends() []string {
	return nil
}

func (s ImportInternetGatewayStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
