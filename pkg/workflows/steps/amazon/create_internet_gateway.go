package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateInternetGateway = "create_internet_gateway"

type CreateInternetGatewayStep struct {
	getIGWService func(cfg steps.AWSConfig) (InternetGatewayCreater, error)
}

type InternetGatewayCreater interface {
	CreateInternetGateway(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error)
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
	AttachInternetGateway(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error)
}

//InitCreateMachine adds the step to the registry
func InitCreateInternetGateway(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateInternetGateway, NewCreateInternetGatewayStep(ec2fn))
}

func NewCreateInternetGatewayStep(ec2fn GetEC2Fn) *CreateInternetGatewayStep {
	return &CreateInternetGatewayStep{
		getIGWService: func(cfg steps.AWSConfig) (InternetGatewayCreater, error) {
			ec2Client, err := ec2fn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepCreateInternetGateway, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return ec2Client, nil
		},
	}
}

func (s *CreateInternetGatewayStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateInternetGateway)

	// Internet gateway already exists
	if cfg.AWSConfig.InternetGatewayID != "" {
		logrus.Debugf("use internet gateway %s",
			cfg.AWSConfig.InternetGatewayID)
		return nil
	} else {
		svc, err := s.getIGWService(cfg.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "error getting IGW service %s",
				StepCreateInternetGateway)
		}

		// Use default gateway for VPC
		resp, err := svc.CreateInternetGateway(new(ec2.CreateInternetGatewayInput))
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
		_, err = svc.CreateTags(tagInput)

		if err != nil {
			logrus.Errorf("Error tagging route table %s %v",
				cfg.AWSConfig.RouteTableID, err)
			return err
		}

		logrus.Debugf("Attach Internet GW %s to VPC %s",
			cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)
		// Attach GW to VPC
		attachGw := &ec2.AttachInternetGatewayInput{
			VpcId:             aws.String(cfg.AWSConfig.VPCID),
			InternetGatewayId: aws.String(cfg.AWSConfig.InternetGatewayID),
		}
		if _, err := svc.AttachInternetGateway(attachGw);
		err != nil && !strings.Contains(err.Error(), "already has an internet gateway attached") {
			logrus.Errorf("Error attaching GW %s to VPC %s", cfg.AWSConfig.InternetGatewayID, cfg.AWSConfig.VPCID)
			return err
		}
	}

	return nil
}

func (s *CreateInternetGatewayStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*CreateInternetGatewayStep) Name() string {
	return StepCreateInternetGateway
}

func (*CreateInternetGatewayStep) Description() string {
	return "Create internet gateway"
}

func (*CreateInternetGatewayStep) Depends() []string {
	return nil
}
