package amazon

import (
	"context"

	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
	"github.com/pkg/errors"
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

	// TODO(stgleb): Filter by VPC-ID here
	EC2.DescribeSubnets(nil)
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
