package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateInternetGateway = "create_internet_gateway"

type CreateInternetGatewayStep struct {
	GetEC2 GetEC2Fn
}

//InitCreateMachine adds the step to the registry
func InitCreateInternetGateway(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateInternetGateway, NewCreateInternetGatewayStep(ec2fn))
}

func NewCreateInternetGatewayStep(ec2fn GetEC2Fn) *CreateInternetGatewayStep {
	return &CreateInternetGatewayStep{
		GetEC2: ec2fn,
	}
}

func (s *CreateInternetGatewayStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateInternetGateway)
	ec2Client, err := s.GetEC2(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v", s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	// Internet gateway already exists
	if cfg.AWSConfig.InternetGatewayID != "" {
		logrus.Debugf("use internet gateway %s",
			cfg.AWSConfig.InternetGatewayID)
		return nil
	} else {
		// Use default gateway for VPC
		resp, err := ec2Client.CreateInternetGateway(new(ec2.CreateInternetGatewayInput))
		if err != nil {
			return err
		}
		cfg.AWSConfig.InternetGatewayID = *resp.InternetGateway.InternetGatewayId

		// Tag gateway
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
				Value: aws.String(fmt.Sprintf("inet-gateway-%s",
					cfg.ClusterID)),
			},
		}

		tagInput := &ec2.CreateTagsInput{
			Resources: []*string{aws.String(cfg.AWSConfig.InternetGatewayID)},
			Tags:      ec2Tags,
		}
		_, err = ec2Client.CreateTags(tagInput)

		if err != nil {
			logrus.Errorf("Error tagging route table %s %v",
				cfg.AWSConfig.RouteTableID, err)
			return err
		}

		// Attach GW to VPC
		attachGw := &ec2.AttachInternetGatewayInput{
			VpcId:             aws.String(cfg.AWSConfig.VPCID),
			InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
		}
		if _, err := ec2Client.AttachInternetGateway(attachGw); err != nil && !strings.Contains(err.Error(), "already has an internet gateway attached") {
			return err
		}
	}

	return nil
}

func (s *CreateInternetGatewayStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*CreateInternetGatewayStep) Name() string {
	return StepNameCreateEC2Instance
}

func (*CreateInternetGatewayStep) Description() string {
	return "Create internet gateway"
}

func (*CreateInternetGatewayStep) Depends() []string {
	return nil
}
