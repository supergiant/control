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

const CreateSubnetStepName = "create_subnet_step"

type CreateSubnetStep struct {
	GetEC2 GetEC2Fn
}

func NewCreateSubnetStep(fn GetEC2Fn) *CreateSubnetStep {
	return &CreateSubnetStep{
		GetEC2: fn,
	}
}

func InitCreateSubnet(fn GetEC2Fn) {
	steps.RegisterStep(CreateSubnetStepName, NewCreateSubnetStep(fn))
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
	} else {
		log.Infof("[%s] - using subnet %s", s.Name(), cfg.AWSConfig.SubnetID)
	}
	return nil
}

func (*CreateSubnetStep) Name() string {
	return CreateSubnetStepName
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
