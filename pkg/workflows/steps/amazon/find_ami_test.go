package amazon

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"go.uber.org/zap/buffer"
	"strings"
	"testing"
)

type mockImageService struct {
	output *ec2.DescribeImagesOutput
	err    error
}

func (m *mockImageService) DescribeImagesWithContext(ctx aws.Context, input *ec2.DescribeImagesInput,
	opts ...request.Option) (*ec2.DescribeImagesOutput, error) {
	return m.output, m.err
}

func TestFindAMIStep_Run(t *testing.T) {
	imageID := "1234"

	testCases := []struct {
		description  string
		getFinderErr error
		output       *ec2.DescribeImagesOutput
		err          error
		errMsg       string
	}{
		{
			description:  "error getting finder",
			getFinderErr: errors.New("error obtaining image finder"),
			errMsg:       "error obtaining image finder",
		},
		{
			description: "error while getting image",
			err:         errors.New("something went wrong"),
			errMsg:      "something went wrong",
		},
		{
			description: "image not found",
			output: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{
					{
						Description: aws.String("UNSUPPORTED"),
					},
				},
			},
			err:    nil,
			errMsg: "can't find supported image",
		},
		{
			description: "success",
			output: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{
					{
						ImageId:     aws.String(imageID),
						Description: aws.String("Ubuntu 16.04"),
					},
				},
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockImageService{
			output: testCase.output,
			err:    testCase.err,
		}

		step := &FindAMIStep{
			getImageService: func(config steps.AWSConfig) (ImageFinder, error) {
				return svc, testCase.getFinderErr
			},
		}

		config := &steps.Config{}
		err := step.Run(context.Background(), &buffer.Buffer{}, config)

		if testCase.errMsg != "" && err == nil {
			t.Error("Error must not be nil")
			continue
		}

		if testCase.errMsg != "" && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Not found expected message %s in err %s",
				testCase.errMsg, err.Error())
			continue
		}

		if testCase.errMsg == "" && err != nil {
			t.Errorf("Unexpected error %v", err)
			continue
		}

		if err == nil && config.AWSConfig.ImageID != imageID {
			t.Errorf("Wrong image id expected %s actual %s",
				imageID, config.AWSConfig.ImageID)
		}
	}
}

func TestNewFindAMIStep(t *testing.T) {
	step := NewFindAMIStep(GetEC2)

	if step.getImageService == nil {
		t.Error("getImageService must not be nil")
	}
}

func TestFindAMIStep_Name(t *testing.T) {
	step := &FindAMIStep{}

	if step.Name() != StepFindAMI {
		t.Errorf("Wrong step name expected %s actual %s",
			StepFindAMI, step.Name())
	}
}

func TestFindAMIStep_Depends(t *testing.T) {
	step := &FindAMIStep{}

	if step.Depends() != nil {
		t.Errorf("Depens list must be nil")
	}
}

func TestFindAMIStep_Rollback(t *testing.T) {
	step := &FindAMIStep{}

	err := step.Rollback(context.Background(), &buffer.Buffer{}, nil)

	if err != nil {
		t.Errorf("Unexpected error while Rollback %v", err)
	}
}

func TestFindAMIStep_Description(t *testing.T) {
	step := &FindAMIStep{}

	if !strings.Contains(step.Description(), "Amazon Machine Image") {
		t.Errorf("Step description %s doesn't contain Amazon "+
			"Machine Image", step.Description())
	}
}
