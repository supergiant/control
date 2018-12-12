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

const DeleteSubnetsStepName = "aws_delete_subnets"

type deleteSubnetesSvc interface {
	DeleteSubnet(*ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error)
}

type DeleteSubnets struct {
	getSvc func(steps.AWSConfig) (deleteSubnetesSvc, error)
	GetEC2 GetEC2Fn
}

func InitDeleteSubnets(fn GetEC2Fn) {
	steps.RegisterStep(DeleteSubnetsStepName, NewDeleteSubnets(fn))
}

func NewDeleteSubnets(fn GetEC2Fn) *DeleteSubnets {
	return &DeleteSubnets{
		getSvc: func(cfg steps.AWSConfig) (deleteSubnetesSvc, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteSubnets) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if len(cfg.AWSConfig.Subnets) == 0 {
		logrus.Debug("Skip deleting empty subnets")
		return nil
	}

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Error getting delete subnets service %v", err)
		return errors.Wrapf(ErrAuthorization, "%s %v",
			DeleteSubnetsStepName, err.Error())
	}

	for az, subnet := range cfg.AWSConfig.Subnets {
		logrus.Debugf("Delete subnet %s in az %s", subnet, az)
		descReq := &ec2.DeleteSubnetInput{
			SubnetId: aws.String(subnet),
		}

		_, err = svc.DeleteSubnet(descReq)

		if err != nil {
			logrus.Debugf("DeleteSubnet caused %s", err.Error())
		}
	}

	return nil
}

func (*DeleteSubnets) Name() string {
	return DeleteSubnetsStepName
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
