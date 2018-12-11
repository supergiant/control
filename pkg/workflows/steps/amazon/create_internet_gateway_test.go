package amazon

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"
	"testing"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/supergiant/control/pkg/workflows/steps"
	"context"
	"bytes"
	"strings"
)

type mockIGWService struct {
	mock.Mock
}

func (m *mockIGWService) CreateInternetGateway(
	input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.CreateInternetGatewayOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockIGWService) CreateTags(
	input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.CreateTagsOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockIGWService) AttachInternetGateway(
	input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.AttachInternetGatewayOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestCreateInternetGatewayStep_Run(t *testing.T) {
	testCases := []struct{
		existingGW string
		getSvcErr error

		createIGWOut *ec2.CreateInternetGatewayOutput
		createIGWErr error

		createTagserr error
		attachErr error

		errMsg string
	}{
		{
			existingGW: "igwID",
		},
		{
			getSvcErr: errors.New("message1"),
			errMsg: "message1",
		},
		{
			createIGWErr: errors.New("message2"),
			errMsg: "message2",
		},
		{
			createIGWOut: &ec2.CreateInternetGatewayOutput{
				InternetGateway: &ec2.InternetGateway{
					InternetGatewayId: aws.String("1234"),
				},
			},
			createIGWErr: errors.New("message3"),
			errMsg: "message3",
		},
		{
			createIGWOut: &ec2.CreateInternetGatewayOutput{
				InternetGateway: &ec2.InternetGateway{
					InternetGatewayId: aws.String("1234"),
				},
			},
			attachErr: errors.New("message4"),
			errMsg: "message4",
		},
		{
			createIGWOut: &ec2.CreateInternetGatewayOutput{
				InternetGateway: &ec2.InternetGateway{
					InternetGatewayId: aws.String("1234"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		svc := &mockIGWService{}
		svc.On("CreateInternetGateway", mock.Anything).
			Return(testCase.createIGWOut, testCase.createIGWErr)
		svc.On("CreateTags", mock.Anything).
			Return(mock.Anything, testCase.attachErr)
		svc.On("AttachInternetGateway", mock.Anything).
			Return(mock.Anything, testCase.attachErr)

		step := &CreateInternetGatewayStep{
			getIGWService: func(cfg steps.AWSConfig) (InternetGatewayCreater, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := steps.Config{
			AWSConfig: steps.AWSConfig{
				InternetGatewayID: testCase.existingGW,
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, &config)

		if err != nil && testCase.errMsg == "" {
			t.Errorf("Unexpected error %v", err)
			continue
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Wrong error must contain %s actual %s",
				testCase.errMsg, err.Error())
		}
	}
}

func TestCreateInternetGatewayStep_Name(t *testing.T) {
	step := &CreateInternetGatewayStep{}

	if step.Name() != StepCreateInternetGateway {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateInternetGateway, step.Name())
	}
}

func TestCreateInternetGatewayStep_Description(t *testing.T) {
	step := &CreateInternetGatewayStep{}

	if step.Description() != "Create internet gateway" {
		t.Errorf("Wrong step description expected Create internet gateway actual %s",
			step.Description())
	}
}

func TestDeleteInternetGateway_Rollback(t *testing.T) {
	step := &CreateInternetGatewayStep{}

	if err := step.Rollback(context.Background(), nil, nil); err != nil {
		t.Errorf("Unexpected error %v while rolling back", err)
	}
}

func TestCreateRouteTableStep_Depends(t *testing.T) {
	step := &CreateInternetGatewayStep{}

	if deps := step.Depends(); deps != nil {
		t.Error("Dependencies must be nil")
	}
}

func TestNewCreateInternetGatewayStep(t *testing.T) {
	step := NewCreateInternetGatewayStep(GetEC2)

	if step == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if step.getIGWService == nil {
		t.Errorf("getIGWService must not be nil")
	}
}