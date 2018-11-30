package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
)

const DeleteVPCStepName = "aws_delete_vpc"

type DeleteVPC struct {
	GetEC2 GetEC2Fn
}

func InitDeleteVPC(fn GetEC2Fn) {
	steps.RegisterStep(DeleteVPCStepName, &DeleteVPC{
		GetEC2: fn,
	})
}

func (s *DeleteVPC) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	req := &ec2.DeleteVpcInput{
		VpcId: aws.String(cfg.AWSConfig.VPCID),
	}

	logrus.Debugf("Delete VPC ID: %s", cfg.AWSConfig.VPCID)
	_, err = EC2.DeleteVpc(req)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("Delete VPC %s caused %s",
			cfg.AWSConfig.VPCID, err.Message())
	}

	return nil
}

func (*DeleteVPC) Name() string {
	return DeleteVPCStepName
}

func (*DeleteVPC) Depends() []string {
	return []string{DeleteSecurityGroupsStepName}
}

func (*DeleteVPC) Description() string {
	return "Delete vpc"
}

func (*DeleteVPC) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
