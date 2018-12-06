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
	"time"
)

const DeleteVPCStepName = "aws_delete_vpc"

var (
	deleteVPCTimeout      = time.Second * 30
	deleteVPCAttemptCount = 3
)

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

	var (
		deleteErr error
		timeout   = deleteVPCTimeout
	)

	for i := 0; i < deleteVPCAttemptCount; i++ {
		req := &ec2.DeleteVpcInput{
			VpcId: aws.String(cfg.AWSConfig.VPCID),
		}

		logrus.Debugf("Delete VPC ID: %s", cfg.AWSConfig.VPCID)
		_, deleteErr = EC2.DeleteVpcWithContext(ctx, req)

		if err, ok := deleteErr.(awserr.Error); ok {
			logrus.Debugf("Delete VPC %s caused %s retry in %v ",
				cfg.AWSConfig.VPCID, err.Message(), timeout)
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			break
		}
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
