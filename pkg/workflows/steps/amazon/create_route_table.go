package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateRouteTable = "create_route_table"

type Service interface {
	CreateRouteTable(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error)
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
	CreateRoute(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error)
}

type CreateRouteTableStep struct {
	getService func(config steps.AWSConfig) (Service, error)
	GetEC2     GetEC2Fn
}

// InitCreateRouteTable adds the step to the registry
func InitCreateRouteTable(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateRouteTable, NewCreateRouteTableStep(ec2fn))
}

func NewCreateRouteTableStep(ec2fn GetEC2Fn) *CreateRouteTableStep {
	return &CreateRouteTableStep{
		getService: func(cfg steps.AWSConfig) (Service, error) {
			ec2Client, err := ec2fn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepCreateRouteTable, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return ec2Client, nil
		},
	}
}

func (s *CreateRouteTableStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateRouteTable)

	//  route table already exists
	if cfg.AWSConfig.RouteTableID != "" {
		return nil
	}

	svc, err := s.getService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("error getting service on step %s %v",
			StepCreateRouteTable, err)
		return errors.Wrapf(err, "error getting service on step %s",
			StepCreateRouteTable)
	}

	createResp, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(cfg.AWSConfig.VPCID),
	})

	if err != nil {
		logrus.Errorf("Error creating route table %v", err)
		return err
	}

	cfg.AWSConfig.RouteTableID = *createResp.RouteTable.RouteTableId
	logrus.Infof("Create route table %s", cfg.AWSConfig.RouteTableID)

	// Tag route table
	ec2Tags := []*ec2.Tag{
		{
			Key:   aws.String("KubernetesCluster"),
			Value: aws.String(cfg.Kube.Name),
		},
		{
			Key:   aws.String(clouds.TagClusterID),
			Value: aws.String(cfg.Kube.ID),
		},
		{
			Key: aws.String("Name"),
			Value: aws.String(fmt.Sprintf("route-table-%s",
				cfg.Kube.ID)),
		},
	}

	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(cfg.AWSConfig.RouteTableID)},
		Tags:      ec2Tags,
	}
	_, err = svc.CreateTags(input)

	if err != nil {
		logrus.Errorf("Error tagging route table %s %v",
			cfg.AWSConfig.RouteTableID, err)
		return err
	}

	// Create route for external connectivity
	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
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
	return StepCreateRouteTable
}

func (*CreateRouteTableStep) Description() string {
	return "Create route table"
}

func (*CreateRouteTableStep) Depends() []string {
	return []string{StepCreateInternetGateway}
}
