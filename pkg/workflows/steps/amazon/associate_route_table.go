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

type Associater interface {
	AssociateRouteTable(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error)
}

type AssociateRouteTableStep struct {
	getRouteTableSvc func(config steps.AWSConfig) (Associater, error)
}

// InitAssociateRouteTable adds the step to the registry
func InitAssociateRouteTable(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepAssociateRouteTable, NewAssociateRouteTableStep(ec2fn))
}

func NewAssociateRouteTableStep(ec2fn GetEC2Fn) *AssociateRouteTableStep {
	return &AssociateRouteTableStep{
		getRouteTableSvc: func(config steps.AWSConfig) (Associater, error) {
			ec2Client, err := ec2fn(config)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepAssociateRouteTable, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return ec2Client, nil
		},
	}
}

func (s *AssociateRouteTableStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepAssociateRouteTable)

	if cfg.AWSConfig.RouteTableAssociationIDs == nil {
		cfg.AWSConfig.RouteTableAssociationIDs = make(map[string]string)
	}

	associater, err := s.getRouteTableSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Debugf("error getting associater service %v", err)
		return errors.Wrapf(err, "error getting associater service %s",
			StepAssociateRouteTable)
	}

	for az, subnet := range cfg.AWSConfig.Subnets {
		logrus.Debugf("Associate route table %s with subnet %s",
			cfg.AWSConfig.RouteTableID, subnet)

		// Associate route table with subnet
		associtationResponse, err := associater.AssociateRouteTable(
			&ec2.AssociateRouteTableInput{
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
	return "Associate route table with all subnets in VPC"
}

func (*AssociateRouteTableStep) Depends() []string {
	return []string{StepCreateSubnets}
}
