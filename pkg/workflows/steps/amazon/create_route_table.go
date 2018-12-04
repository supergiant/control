package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateRouteTable = "create_route_table"

type CreateRouteTableStep struct {
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitCreateRouteTable(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateRouteTable, NewCreateRouteTableStep(ec2fn))
}

func NewCreateRouteTableStep(ec2fn GetEC2Fn) *CreateRouteTableStep {
	return &CreateRouteTableStep{
		GetEC2: ec2fn,
	}
}

func (s *CreateRouteTableStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateRouteTable)
	ec2Client, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	//  route table already exists
	if cfg.AWSConfig.RouteTableID != "" {
		return nil
	}

	createResp, err := ec2Client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(cfg.AWSConfig.VPCID),
	})

	if err != nil {
		logrus.Errorf("Error creating route table %v", err)
		return err
	}

	cfg.AWSConfig.RouteTableID = *createResp.RouteTable.RouteTableId

	// Tag route table
	ec2Tags := []*ec2.Tag{
		{
			Key:   aws.String("KubernetesCluster"),
			Value: aws.String(cfg.ClusterName),
		},
		{
			Key:   aws.String(clouds.ClusterIDTag),
			Value: aws.String(cfg.ClusterID),
		},
		{
			Key: aws.String("Name"),
			Value: aws.String(fmt.Sprintf("route-table-%s",
				cfg.ClusterID)),
		},
	}

	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(cfg.AWSConfig.RouteTableID)},
		Tags:      ec2Tags,
	}
	_, err = ec2Client.CreateTags(input)

	if err != nil {
		logrus.Errorf("Error tagging route table %s %v",
			cfg.AWSConfig.RouteTableID, err)
		return err
	}

	// Create route for external connectivity
	_, err = ec2Client.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(cfg.AWSConfig.RouteTableID),
		GatewayId:            aws.String(cfg.AWSConfig.InternetGatewayID),
	})

	if err != nil {
		logrus.Debugf("Error creating rule for internet gateway %v", err)
		return err
	}

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
	return []string{StepCreateInternetGateway}
}
