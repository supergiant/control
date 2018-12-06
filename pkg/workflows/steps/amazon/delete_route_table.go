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

const DeleteRouteTableStepName = "aws_delete_route_table"

type DeleteRouteTable struct {
	GetEC2 GetEC2Fn
}

func InitDeleteRouteTable(fn GetEC2Fn) {
	steps.RegisterStep(DeleteRouteTableStepName, &DeleteRouteTable{
		GetEC2: fn,
	})
}

func (s *DeleteRouteTable) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.RouteTableID == "" {
		return nil
	}

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	logrus.Debugf("Delete route table %s from VPC %s",
		cfg.AWSConfig.RouteTableID, cfg.AWSConfig.VPCID)
	_, err = EC2.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(cfg.AWSConfig.RouteTableID),
	})

	if err, ok := err.(awserr.Error); ok {
		logrus.Debugf("DisassociateRouteTable caused %s", err.Message())
	}

	return nil
}

func (*DeleteRouteTable) Name() string {
	return DeleteRouteTableStepName
}

func (*DeleteRouteTable) Depends() []string {
	return []string{DisassociateRouteTableStepName}
}

func (*DeleteRouteTable) Description() string {
	return "Delete route table from vpc"
}

func (*DeleteRouteTable) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
