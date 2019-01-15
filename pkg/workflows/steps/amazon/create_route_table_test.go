package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
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
	testCases := []struct {
		description         string
		existingRouteTable  string
		getSvcError         error
		createOut           *ec2.CreateRouteTableOutput
		createRouteTableErr error
		tagErr              error
		createRouteErr      error
		errMsg              string
	}{
		{
			description:        "existing route table",
			existingRouteTable: "1234",
		},
		{
			description: "get service error",
			getSvcError: errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description:         "create route table error",
			createRouteTableErr: errors.New("message2"),
			errMsg:              "message2",
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
			errMsg:         "message4",
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

func TestInitCreateRouteTable(t *testing.T) {
	InitCreateRouteTable(GetEC2)

	s := steps.GetStep(StepCreateRouteTable)

	if s == nil {
		t.Errorf("Step %s not found", StepCreateRouteTable)
	}
}

func TestNewCreateRouteTableStep(t *testing.T) {
	step := NewCreateRouteTableStep(GetEC2)

	if step == nil {
		t.Error("Step must not be nil")
	}

	if step.getService == nil {
		t.Errorf("%s getService must not be nil", StepCreateRouteTable)
	}

	if api, err := step.getService(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}

func TestNewCreateRouteTableStepError(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewCreateRouteTableStep(fn)

	if step == nil {
		t.Error("Step must not be nil")
	}

	if step.getService == nil {
		t.Errorf("%s getService must not be nil", StepCreateRouteTable)
	}

	if api, err := step.getService(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}

func TestCreateRouteTableStep_Name(t *testing.T) {
	step := NewCreateRouteTableStep(GetEC2)

	if step.Name() != StepCreateRouteTable {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateRouteTable, step.Name())
	}
}

func TestCreateRouteTableStep_Description(t *testing.T) {
	step := NewCreateRouteTableStep(GetEC2)

	if step.Description() != "Create route table" {
		t.Errorf("Wrong description expected "+
			"Create route table actual %s", step.Description())
	}
}

func TestCreateRouteTableStep_Depends(t *testing.T) {
	step := NewCreateRouteTableStep(GetEC2)

	if deps := step.Depends(); len(deps) != 1 || deps[0] != StepCreateInternetGateway {
		t.Errorf("Wrong dependencies %v expected %v",
			deps, []string{StepCreateInternetGateway})
	}
}

func TestCreateRouteTableStep_Rollback(t *testing.T) {
	step := NewCreateRouteTableStep(GetEC2)

	if err := step.Rollback(context.Background(),
		&bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v while Rollback", err)
	}
}
