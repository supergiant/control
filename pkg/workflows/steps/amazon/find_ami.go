package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepFindAMI = "find_amazon_machine_image"

type FindAMIStep struct {
	GetEC2 GetEC2Fn
}

func NewFindAMIStep(fn GetEC2Fn) *FindAMIStep {
	return &FindAMIStep{
		GetEC2: fn,
	}
}

func InitFindAMI(fn GetEC2Fn) {
	steps.RegisterStep(StepFindAMI, NewFindAMIStep(fn))
}

func (s *FindAMIStep) Run(ctx context.Context, w io.Writer,
	cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v",
			s.Name(), err)
		return errors.Wrap(ErrAuthorization, err.Error())
	}

	cfg.AWSConfig.ImageID, err = s.FindAMI(ctx, w, EC2)
	logrus.Debugf("Found image id %s", cfg.AWSConfig.ImageID)
	if err != nil {
		logrus.Errorf("[%s] - failed to find AMI for Ubuntu: %v",
			s.Name(), err)
		return errors.Wrap(err, "failed to find AMI")
	}

	return nil
}

func (*FindAMIStep) Name() string {
	return StepFindAMI
}

func (*FindAMIStep) Description() string {
	return "Step looks for Amazon Machine Image with specified parameters"
}

func (*FindAMIStep) Depends() []string {
	return nil
}

func (*FindAMIStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *FindAMIStep) FindAMI(ctx context.Context, w io.Writer, EC2 ec2iface.EC2API) (string, error) {
	// TODO: should it be configurable?
	out, err := EC2.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("architecture"),
				Values: []*string{
					aws.String("x86_64"),
				},
			},
			{
				Name: aws.String("virtualization-type"),
				Values: []*string{
					aws.String("hvm"),
				},
			},
			{
				Name: aws.String("root-device-type"),
				Values: []*string{
					aws.String("ebs"),
				},
			},
			//Owner should be Canonical
			{
				Name: aws.String("owner-id"),
				Values: []*string{
					aws.String("099720109477"),
				},
			},
			{
				Name: aws.String("description"),
				Values: []*string{
					aws.String("Canonical, Ubuntu, 16.04*"),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	amiID := ""

	log := util.GetLogger(w)

	for _, img := range out.Images {
		if img.Description == nil {
			continue
		}
		if strings.Contains(*img.Description, "UNSUPPORTED") {
			continue
		}
		amiID = *img.ImageId

		logMessage := fmt.Sprintf("[%s] - using AMI (ID: %s) %s", s.Name(), amiID, *img.Description)
		log.Info(logMessage)
		logrus.Info(logMessage)

		break
	}

	return amiID, nil
}
