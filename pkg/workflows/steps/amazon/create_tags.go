package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepCreateTags = "create_tags"

type TagService interface {
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
}

type CreateTagsStep struct {
	getService func(config steps.AWSConfig) (TagService, error)
	GetEC2     GetEC2Fn
}

// InitCreateRouteTable adds the step to the registry
func InitCreateTagsStep(ec2fn GetEC2Fn) {
	steps.RegisterStep(StepCreateTags, NewCreateTagsStep(ec2fn))
}

func NewCreateTagsStep(ec2fn GetEC2Fn) *CreateTagsStep {
	return &CreateTagsStep{
		getService: func(cfg steps.AWSConfig) (TagService, error) {
			ec2Client, err := ec2fn(cfg)

			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepCreateTags, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return ec2Client, nil
		},
	}
}

func (s *CreateTagsStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	logrus.Debugf(StepCreateTags)

	svc, err := s.getService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("error getting service on step %s %v",
			StepCreateTags, err)
		return errors.Wrapf(err, "error getting service on step %s",
			StepCreateTags)
	}

	resourceIds := make([]*string, 0)
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.RouteTableID))
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.VPCID))
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.RouteTableID))
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.InternetGatewayID))
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.NodesSecurityGroupID))
	resourceIds = append(resourceIds, aws.String(cfg.AWSConfig.MastersSecurityGroupID))

	for _, subnetId := range cfg.AWSConfig.Subnets {
		resourceIds = append(resourceIds, aws.String(subnetId))
	}

	input := &ec2.CreateTagsInput{
		Resources: resourceIds,
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("KubernetesCluster"),
				Value: aws.String(cfg.ClusterName),
			},
			{
				Key:   aws.String(clouds.TagClusterID),
				Value: aws.String(cfg.ClusterID),
			},
		},
	}

	_, err = svc.CreateTags(input)

	if err != nil {
		logrus.Debugf("Error creating tags for aws entities %v", err)
		return err
	}

	return nil
}

func (s *CreateTagsStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*CreateTagsStep) Name() string {
	return StepCreateTags
}

func (*CreateTagsStep) Description() string {
	return "Create tags for all infra objects"
}

func (*CreateTagsStep) Depends() []string {
	return nil
}
