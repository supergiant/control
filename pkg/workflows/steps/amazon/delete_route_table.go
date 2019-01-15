package amazon

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const DeleteRouteTableStepName = "aws_delete_route_table"

var (
	deleteRouteAttemptCount = 5
	deleteRouteTimeout      = time.Minute
)

type deleteRouteTableSvc interface{
	DeleteRouteTable(*ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error)
}

type DeleteRouteTable struct {
	getSvc func(steps.AWSConfig) (deleteRouteTableSvc, error)
}

func InitDeleteRouteTable(fn GetEC2Fn) {
	steps.RegisterStep(DeleteRouteTableStepName, NewDeleteRouteTableStep(fn))
}

func NewDeleteRouteTableStep(fn GetEC2Fn) *DeleteRouteTable {
	return &DeleteRouteTable{
		getSvc: func(config steps.AWSConfig) (deleteRouteTableSvc, error) {
			EC2, err := fn(config)
			if err != nil {
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}
func (s *DeleteRouteTable) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.AWSConfig.RouteTableID == "" {
		logrus.Debug("Skip deleting empty route table")
		return nil
	}

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Error getting delete service %v", err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	var (
		deleteErr error
		timeout   = deleteRouteTimeout
	)

	// Disassociating of route table and subnets can take a while, we need to be patient
	for i := 0; i < deleteRouteAttemptCount; i++ {
		logrus.Debugf("Delete route table %s from VPC %s",
			cfg.AWSConfig.RouteTableID, cfg.AWSConfig.VPCID)
		_, deleteErr = svc.DeleteRouteTable(&ec2.DeleteRouteTableInput{
			RouteTableId: aws.String(cfg.AWSConfig.RouteTableID),
		})

		if deleteErr != nil {
			logrus.Debugf("Delete route table %s caused %s sleep for %v",
				cfg.AWSConfig.RouteTableID, deleteErr.Error(), timeout)
			time.Sleep(timeout)
			timeout = timeout * 2
		} else {
			break
		}
	}

	return nil
}

func (*DeleteRouteTable) Name() string {
	return DeleteRouteTableStepName
}

func (*DeleteRouteTable) Depends() []string {
	return []string{DeleteSubnetsStepName}
}

func (*DeleteRouteTable) Description() string {
	return "Delete route table from vpc"
}

func (*DeleteRouteTable) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
