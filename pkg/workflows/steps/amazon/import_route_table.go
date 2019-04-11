package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"io"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	ImporRouteTablesStepName = "import_route_tables_aws"
)

type RouteTableDescriber interface {
	DescribeRouteTablesWithContext(aws.Context, *ec2.DescribeRouteTablesInput, ...request.Option) (*ec2.DescribeRouteTablesOutput, error)
}

type ImportRouteTablesStep struct {
	getSvc func(config steps.AWSConfig) (RouteTableDescriber, error)
}

func NewImportRouteTablesStep(fn GetEC2Fn) *ImportRouteTablesStep {
	return &ImportRouteTablesStep{
		getSvc: func(config steps.AWSConfig) (describer RouteTableDescriber, e error) {
			EC2, err := fn(config)

			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func InitImportRouteTablesStep(fn GetEC2Fn) {
	steps.RegisterStep(ImporRouteTablesStepName, NewImportRouteTablesStep(fn))
}

func (s ImportRouteTablesStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	logrus.Info(ImporRouteTablesStepName)
	ec2Svc, err := s.getSvc(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	req := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(cfg.AWSConfig.VPCID)},
			},
		},
	}
	output, err := ec2Svc.DescribeRouteTablesWithContext(ctx, req)

	if err != nil {
		return errors.Wrapf(err, "")
	}

	if len(output.RouteTables) == 0 {
		return errors.Wrapf(err, "route tables not found in %s", cfg.AWSConfig.VPCID)
	}

	routeTable := output.RouteTables[0]
	cfg.AWSConfig.RouteTableID = *routeTable.RouteTableId
	cfg.AWSConfig.RouteTableAssociationIDs = make(map[string]string)

	for _, association := range routeTable.Associations {
		 for az, subnetId := range cfg.AWSConfig.Subnets {
		 	if association != nil && association.SubnetId != nil && association.RouteTableAssociationId != nil && subnetId == *association.SubnetId {
				cfg.AWSConfig.RouteTableAssociationIDs[az] = *association.RouteTableAssociationId
			}
		 }
	}

	return nil
}

func (s ImportRouteTablesStep) Name() string {
	return ImportClusterMachinesStepName
}

func (s ImportRouteTablesStep) Description() string {
	return ImportClusterMachinesStepName
}

func (s ImportRouteTablesStep) Depends() []string {
	return nil
}

func (s ImportRouteTablesStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
