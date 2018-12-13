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

const DisassociateRouteTableStepName = "aws_disassociate_route_table"

type DisassociateService interface {
	DisassociateRouteTable(*ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error)
}

type DisassociateRouteTable struct {
	getSvc func(steps.AWSConfig) (DisassociateService, error)
}

func InitDisassociateRouteTable(fn GetEC2Fn) {
	steps.RegisterStep(DisassociateRouteTableStepName,
		NewDisassociateRouteTableStep(fn))
}

func NewDisassociateRouteTableStep(fn GetEC2Fn) *DisassociateRouteTable {
	return &DisassociateRouteTable{
		getSvc: func(cfg steps.AWSConfig) (DisassociateService, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func (s *DisassociateRouteTable) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("%s error getting service %v",
			DisassociateRouteTableStepName, err)
		return errors.Wrapf(err, "Step %s getting service error",
			DisassociateRouteTableStepName)
	}

	for _, associationID := range cfg.AWSConfig.RouteTableAssociationIDs {
		if associationID == "" {
			continue
		}

		disReq := &ec2.DisassociateRouteTableInput{
			AssociationId: aws.String(associationID),
		}

		_, err = svc.DisassociateRouteTable(disReq)

		if err != nil {
			logrus.Debugf("DisassociateRouteTable caused %s", err.Error())
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
