package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
	"time"
)

const DeleteVPCStepName = "aws_delete_vpc"

var (
	deleteVPCTimeout      = time.Second * 30
	deleteVPCAttemptCount = 3
)

type vpcSvc interface {
	DeleteVpcWithContext(aws.Context, *ec2.DeleteVpcInput, ...request.Option) (*ec2.DeleteVpcOutput, error)
}

type DeleteVPC struct {
	getSvc func(steps.AWSConfig) (vpcSvc, error)
	GetEC2 GetEC2Fn
}

func InitDeleteVPC(fn GetEC2Fn) {
	steps.RegisterStep(DeleteVPCStepName, NewDeleteVPC(fn))
}

func NewDeleteVPC(fn GetEC2Fn) *DeleteVPC {
	return &DeleteVPC{
		getSvc: func(cfg steps.AWSConfig) (vpcSvc, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DeleteVPC) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.VPCID == "" {
		logrus.Debug("Skip deleting empty VPC")
		return nil
	}

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	var (
		deleteErr error
		timeout   = deleteVPCTimeout
	)

	for i := 0; i < deleteVPCAttemptCount; i++ {
		req := &ec2.DeleteVpcInput{
			VpcId: aws.String(cfg.AWSConfig.VPCID),
		}

		logrus.Debugf("Delete VPC ID: %s", cfg.AWSConfig.VPCID)
		_, deleteErr = svc.DeleteVpcWithContext(ctx, req)

		if deleteErr != nil {
			logrus.Debugf("Delete VPC %s caused %s retry in %v ",
				cfg.AWSConfig.VPCID, deleteErr.Error(), timeout)
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			break
		}
	}

	return deleteErr
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
