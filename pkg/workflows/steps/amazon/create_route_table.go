package amazon

import (
	"github.com/supergiant/control/pkg/workflows/steps"
	"io"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const CreateRouteTableStepName = "create_route_table"

type CreateRouteTableStep struct{
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitCreateRouteTable(ec2fn GetEC2Fn) {
	steps.RegisterStep(CreateRouteTableStepName, NewCreateRouteTableStep(ec2fn))
}


func NewCreateRouteTableStep(ec2fn GetEC2Fn) *StepCreateInstance {
	return &StepCreateInstance{
		GetEC2: ec2fn,
	}
}

func (s *CreateRouteTableStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	ec2Client, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	//  route table already exists
	if cfg.AWSConfig.RouteTableID != ""{
		return nil
	}

	resp, err := ec2Client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(cfg.AWSConfig.VPCID),
	})

	if err != nil {
		logrus.Errorf("Error creating route table %v", err)
		return err
	}

	cfg.AWSConfig.RouteTableID = *resp.RouteTable.RouteTableId

	// TODO(stgleb): tag route table


	associtationResponse, err := ec2Client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(cfg.AWSConfig.RouteTableID),
		SubnetId:     aws.String(cfg.AWSConfig.SubnetID),
	})
	if err != nil {
		return err
	}

	cfg.AWSConfig.RouteTableSubnetAssociationID = *associtationResponse.AssociationId


	return nil
}

func (s *CreateRouteTableStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}


func (*CreateRouteTableStep) Name() string {
	return StepNameCreateEC2Instance
}

func (*CreateRouteTableStep) Description() string {
	return "Create route table"
}

func (*CreateRouteTableStep) Depends() []string {
	return []string{StepCreateSubnet}
}

