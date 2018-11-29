package amazon

import (
	"context"
	"io"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/sirupsen/logrus"
)

const DeleteSecurityGroupsStepName = "aws_delete_security_groups"

type DeleteSecurityGroup struct {
	GetEC2 GetEC2Fn
}

func InitDeleteSecurityGroup(fn GetEC2Fn) {
	steps.RegisterStep(DeleteSecurityGroupsStepName, &DeleteSecurityGroup{
		GetEC2: fn,
	})
}

func (s *DeleteSecurityGroup) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Delete security group %s", cfg.AWSConfig.MastersSecurityGroupID)
	reqMaster := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.MastersSecurityGroupID),
	}

	_, err = EC2.DeleteSecurityGroup(reqMaster)

	if err != nil {
		return errors.Wrapf(err, "%s: master security groups", DeleteSecurityGroupsStepName)
	}

	reqNode := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(cfg.AWSConfig.NodesSecurityGroupID),
	}

	logrus.Debugf("Delete security group %s", cfg.AWSConfig.NodesSecurityGroupID)
	_, err = EC2.DeleteSecurityGroup(reqNode)

	if err != nil {
		return errors.Wrapf(err, "%s: nodes security groups", DeleteSecurityGroupsStepName)
	}

	logrus.Debugf("Deleting security group finished")
	return nil
}

func (*DeleteSecurityGroup) Name() string {
	return DeleteSecurityGroupsStepName
}

func (*DeleteSecurityGroup) Depends() []string {
	return []string{DeleteClusterMachinesStepName}
}

func (*DeleteSecurityGroup) Description() string {
	return "Deletes security groups"
}

func (*DeleteSecurityGroup) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
