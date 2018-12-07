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

const DisassociateRouteTableStepName = "aws_disassociate_route_table"

type DisassociateRouteTable struct {
	GetEC2 GetEC2Fn
}

func InitDisassociateRouteTable(fn GetEC2Fn) {
	steps.RegisterStep(DisassociateRouteTableStepName, &DisassociateRouteTable{
		GetEC2: fn,
	})
}

func (s *DisassociateRouteTable) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	for _, associationID := range cfg.AWSConfig.RouteTableAssociationIDs {
		if associationID == "" {
			continue
		}

		disReq := &ec2.DisassociateRouteTableInput{
			AssociationId: aws.String(associationID),
		}

		_, err = EC2.DisassociateRouteTable(disReq)

		if err, ok := err.(awserr.Error); ok {
			logrus.Debugf("DisassociateRouteTable caused %s", err.Message())
		}
	}

	return nil
}

func (*DisassociateRouteTable) Name() string {
	return DisassociateRouteTableStepName
}

func (*DisassociateRouteTable) Depends() []string {
	return []string{DeleteSecurityGroupsStepName}
}

func (*DisassociateRouteTable) Description() string {
	return "Disassociate route table with subnets"
}

func (*DisassociateRouteTable) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
