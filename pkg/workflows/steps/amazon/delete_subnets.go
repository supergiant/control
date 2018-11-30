package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteSubnetsStepName = "aws_delete_subnets"

type DeleteSubnets struct {
	GetEC2 GetEC2Fn
}

func InitDeleteSubnets(fn GetEC2Fn) {
	steps.RegisterStep(DeleteSubnetsStepName, &DeleteSubnets{
		GetEC2: fn,
	})
}

func (s *DeleteSubnets) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	descReq := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(cfg.AWSConfig.SubnetID),
	}

	_, err = EC2.DeleteSubnet(descReq)

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("DeleteSubnet caused %s", err.Message())
	}

	return nil
}

func (*DeleteSubnets) Name() string {
	return DeleteSecurityGroupsStepName
}

func (*DeleteSubnets) Depends() []string {
	return []string{DeleteSecurityGroupsStepName}
}

func (*DeleteSubnets) Description() string {
	return "Deletes security groups"
}

func (*DeleteSubnets) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
