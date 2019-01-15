package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepFindAMI = "find_amazon_machine_image"

type ImageFinder interface {
	DescribeImagesWithContext(aws.Context, *ec2.DescribeImagesInput,
		...request.Option) (*ec2.DescribeImagesOutput, error)
}

type FindAMIStep struct {
	getImageService func(config steps.AWSConfig) (ImageFinder, error)
}

func NewFindAMIStep(fn GetEC2Fn) *FindAMIStep {
	return &FindAMIStep{
		getImageService: func(config steps.AWSConfig) (ImageFinder, error) {
			EC2, err := fn(config)
			if err != nil {
				logrus.Errorf("[%s] - failed to authorize in AWS: %v",
					StepFindAMI, err)
				return nil, errors.Wrap(ErrAuthorization, err.Error())
			}

			return EC2, nil
		},
	}
}

func InitFindAMI(fn GetEC2Fn) {
	steps.RegisterStep(StepFindAMI, NewFindAMIStep(fn))
}

func (s *FindAMIStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	finder, err := s.getImageService(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("[%s] - failed to authorize in AWS: %v",
			s.Name(), err)
		return errors.Wrap(err, StepFindAMI)
	}

	imageID, err := s.FindAMI(ctx, w, finder)
	logrus.Debugf("Found image id %s", cfg.AWSConfig.ImageID)

	if err != nil {
		logrus.Errorf("[%s] - failed to find AMI for Ubuntu: %v",
			s.Name(), err)
		return errors.Wrap(err, "failed to find AMI")
	}

	if err == nil && imageID == "" {
		logrus.Debugf("[%s] - can't find supported image", s.Name())
		return errors.New(fmt.Sprintf("[%s] - can't find "+
			"supported image", s.Name()))
	}

	logrus.Debugf("Use image id %s", imageID)
	cfg.AWSConfig.ImageID = imageID

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

func (s *FindAMIStep) FindAMI(ctx context.Context, w io.Writer, finder ImageFinder) (string, error) {
	// TODO: should it be configurable?
	out, err := finder.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
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
