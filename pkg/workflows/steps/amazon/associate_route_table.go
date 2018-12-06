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

// InitAssociateRouteTable adds the step to the registry
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

	if cfg.AWSConfig.RouteTableAssociationIDs == nil {
		cfg.AWSConfig.RouteTableAssociationIDs = make(map[string]string)
	}

	for az, subnet := range cfg.AWSConfig.Subnets {
		logrus.Debugf("Associate route table %s with subnet %s",
			cfg.AWSConfig.RouteTableID, cfg.AWSConfig.Subnets[cfg.AWSConfig.AvailabilityZone])

		// Associate route table with subnet
		associtationResponse, err := ec2Client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: aws.String(cfg.AWSConfig.RouteTableID),
			SubnetId:     aws.String(subnet),
		})

		// Skip it since by default route table is associated with default subnet
		if err != nil {
			logrus.Debugf("error associating route table %s with subnet %s in az %s %v",
				cfg.AWSConfig.RouteTableID,
				subnet,
				az,
				err)
			return nil
		}

		// Save this for later
		cfg.AWSConfig.RouteTableAssociationIDs[az] = *associtationResponse.AssociationId
	}

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
	return []string{StepCreateSubnets}
}
