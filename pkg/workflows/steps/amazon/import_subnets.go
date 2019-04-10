package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"io"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	ImportSubnetsStepName = "import_subnets_aws"
)

type SubnetDescriber interface {
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
}

type ImportSubnetsStep struct {
	getSvc func(config steps.AWSConfig) (SubnetDescriber, error)
}

func NewImportSubnetsStep(fn GetEC2Fn) *ImportSubnetsStep {
	return &ImportSubnetsStep{
		getSvc: func(config steps.AWSConfig) (describer SubnetDescriber, e error) {
			EC2, err := fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func InitImportSubnetDescriber(fn GetEC2Fn) {
	steps.RegisterStep(ImportSubnetsStepName, NewImportSubnetsStep(fn))
}

func (s ImportSubnetsStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	ec2Svc, err := s.getSvc(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	req := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(cfg.AWSConfig.VPCID)},
			},
		},
	}
	output, err := ec2Svc.DescribeSubnets(req)

	if err != nil {
		return errors.Wrapf(err, "")
	}

	for _, subnet := range output.Subnets {
		cfg.AWSConfig.Subnets[*subnet.AvailabilityZone] = *subnet.SubnetId
	}

	return nil
}

func (s ImportSubnetsStep) Name() string {
	return ImportClusterMachinesStepName
}

func (s ImportSubnetsStep) Description() string {
	return ImportClusterMachinesStepName
}

func (s ImportSubnetsStep) Depends() []string {
	return nil
}

func (s ImportSubnetsStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
