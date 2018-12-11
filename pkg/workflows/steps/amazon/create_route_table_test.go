package amazon

import (
	"context"
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreateRouteTable(input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.CreateRouteTableOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockService) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.CreateTagsOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockService) CreateRoute(input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.CreateRouteOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestCreateRouteTableStep_Run(t *testing.T) {
	testCases := []struct{
		description string
		existingRouteTable string
		getSvcError error
		createOut *ec2.CreateRouteTableOutput
		createRouteTableErr error
		tagErr error
		createRouteErr error
		errMsg string
	}{
		{
			description: "existing route table",
			existingRouteTable: "1234",
		},
		{
			description: "get service error",
			getSvcError: errors.New("message1"),
			errMsg: "message1",
		},
		{
			description: "create route table error",
			createRouteTableErr: errors.New("message2"),
			errMsg: "message2",
		},
		{
			description: "tag error",
			createOut: &ec2.CreateRouteTableOutput{
				RouteTable: &ec2.RouteTable{
					RouteTableId: aws.String("1234"),
				},
			},
			tagErr: errors.New("message3"),
			errMsg: "message3",
		},
		{
			description: "create route error",
			createOut: &ec2.CreateRouteTableOutput{
				RouteTable: &ec2.RouteTable{
					RouteTableId: aws.String("1234"),
				},
			},
			createRouteErr: errors.New("message4"),
			errMsg: "message4",
		},
		{
			description: "success",
			createOut: &ec2.CreateRouteTableOutput{
				RouteTable: &ec2.RouteTable{
					RouteTableId: aws.String("1234"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockService{}
		svc.On("CreateRouteTable", mock.Anything).
			Return(testCase.createOut, testCase.createRouteTableErr)
		svc.On("CreateTags", mock.Anything).
			Return(mock.Anything, testCase.tagErr)
		svc.On("CreateRoute", mock.Anything).
			Return(mock.Anything, testCase.createRouteErr)

		step := &CreateRouteTableStep{
			getService: func(cfg steps.AWSConfig) (Service, error) {
				return svc, testCase.getSvcError
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				RouteTableID: testCase.existingRouteTable,
			},
		}
		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err != nil && testCase.errMsg == "" {
			t.Errorf("Unexpected error %v", err)
			continue
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s does not contain expected text %s",
				err.Error(), testCase.errMsg)
			continue
		}

		if testCase.createOut != nil &&
			config.AWSConfig.RouteTableID != *testCase.createOut.RouteTable.RouteTableId {
			t.Errorf("Wrong Route Table ID expected %s actual %s",
				*testCase.createOut.RouteTable.RouteTableId, config.AWSConfig.RouteTableID)
		}
	}
}