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

const StepAssociateRouteTable = "associate_route_table"

type AssociateRouteTableStep struct {
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitAssociateRouteTable(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepAssociateRouteTable, NewAssociateRouteTableStep(ec2fn))
}

func NewAssociateRouteTableStep(ec2fn GetEC2Fn) *AssociateRouteTableStep {
	return &AssociateRouteTableStep{
		GetEC2: ec2fn,
	}
}

func (s *AssociateRouteTableStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepAssociateRouteTable)
	ec2Client, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	// Associate route table with subnet
	associtationResponse, err := ec2Client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(cfg.AWSConfig.RouteTableID),
		SubnetId:     aws.String(cfg.AWSConfig.SubnetID),
	})

	// Skip it since by default route table is associated with default subnet
	if err != nil {
		logrus.Debugf("error associating route table %s with subnet %s %v",
			cfg.AWSConfig.RouteTableID,
			cfg.AWSConfig.SubnetID,
			err)
		return nil
	}

	cfg.AWSConfig.RouteTableSubnetAssociationID = *associtationResponse.AssociationId

	return nil
}

func (s *AssociateRouteTableStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*AssociateRouteTableStep) Name() string {
	return StepAssociateRouteTable
}

func (*AssociateRouteTableStep) Description() string {
	return "Create route table"
}

func (*AssociateRouteTableStep) Depends() []string {
	return []string{StepCreateSubnet}
}
