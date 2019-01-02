package amazon

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
	"strings"
	"time"
)

type mockSecurityGroupSvc struct {
	mock.Mock
}

func (m *mockSecurityGroupSvc) CreateSecurityGroupWithContext(ctx aws.Context,
	req *ec2.CreateSecurityGroupInput, opts ...request.Option) (*ec2.CreateSecurityGroupOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.CreateSecurityGroupOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockSecurityGroupSvc) AuthorizeSecurityGroupIngressWithContext(ctx aws.Context,
	req *ec2.AuthorizeSecurityGroupIngressInput, opts ...request.Option) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.AuthorizeSecurityGroupIngressOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestCreateSecurityGroupsStep_Run(t *testing.T) {
	testCases := []struct {
		description string
		getSvcErr   error

		createMasterGroupErr    error
		createMasterGroupOutput *ec2.CreateSecurityGroupOutput

		createNodeGroupErr    error
		createNodeGroupOutput *ec2.CreateSecurityGroupOutput

		authorizeMasterSshErr error
		authorizeNodeSshErr   error

		authorizeAllErr1 error
		authorizeAllErr2 error

		findOutboundIP func() (string, error)

		whiteListErr1  error
		whiteListErr2  error

		errMsg string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description:          "create master error",
			createMasterGroupErr: errors.New("message2"),
			errMsg:               "message2",
		},
		{
			description: "create node error",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupErr: errors.New("message3"),
			errMsg:             "message3",
		},
		{
			description: "authorize ssh error #1",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			authorizeMasterSshErr: errors.New("message4"),
			errMsg:                "message4",
		},
		{
			description: "authorize ssh error #2",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			authorizeNodeSshErr: errors.New("message5"),
			errMsg:              "message5",
		},
		{
			description: "allow all traffic #1",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			authorizeAllErr1: errors.New("message5"),
			errMsg:            "message5",
		},
		{
			description: "allow all traffic #2",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			authorizeAllErr2: errors.New("message6"),
			errMsg:            "message6",
		},
		{
			description: "find outbound error",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			findOutboundIP: func() (string, error) {
				return "", errors.New("message9")
			},
			errMsg: "message9",
		},
		{
			description: "white list SG IP #1",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			findOutboundIP: func() (string, error) {
				return "10.20.30.40", nil
			},
			whiteListErr1: errors.New("message7"),
			errMsg:        "message7",
		},
		{
			description: "white list SG IP #2",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			findOutboundIP: func() (string, error) {
				return "10.20.30.40", nil
			},
			whiteListErr2: errors.New("message8"),
			errMsg:        "message8",
		},
		{
			description: "success",
			createMasterGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("masterID"),
			},
			createNodeGroupOutput: &ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("nodeID"),
			},
			findOutboundIP: func() (string, error) {
				return "10.20.30.40", nil
			},
		},
	}

	// Set timeout and attempts count for finding outbound IP
	attempts = 1
	timeout = time.Nanosecond

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockSecurityGroupSvc{}
		svc.On("CreateSecurityGroupWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.createMasterGroupOutput,
				testCase.createMasterGroupErr).Once()

		svc.On("CreateSecurityGroupWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.createNodeGroupOutput,
				testCase.createNodeGroupErr).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.authorizeMasterSshErr).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.authorizeNodeSshErr).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.authorizeAllErr1).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.authorizeAllErr2).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.whiteListErr1).Once()

		svc.On("AuthorizeSecurityGroupIngressWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything,
				testCase.whiteListErr2).Once()

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				VPCID: "1234",
			},
		}

		step := &CreateSecurityGroupsStep{
			getSvc: func(config steps.AWSConfig) (secGroupService, error) {
				return svc, testCase.getSvcErr
			},
			findOutboundIP: testCase.findOutboundIP,
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Wrong error message %v must contain %s",
				err, testCase.errMsg)
		}
	}
}

func TestInitCreateSecurityGroups(t *testing.T) {
	InitCreateSecurityGroups(GetEC2)

	s := steps.GetStep(StepCreateSecurityGroups)

	if s == nil {
		t.Errorf("Step must not be nil")
	}
}

func TestNewCreateSecurityGroupsStep(t *testing.T) {
	s := NewCreateSecurityGroupsStep(GetEC2)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewCreateSecurityGroupsStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewCreateSecurityGroupsStep(fn)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("get svc must not be nil")
	}

	if s.findOutboundIP == nil {
		t.Errorf("FindOutboundIP must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestCreateSecurityGroupsStep_Depends(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Deps must be nil")
	}
}

func TestCreateSecurityGroupsStep_Name(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if name := s.Name(); name != StepCreateSecurityGroups {
		t.Errorf("Wrong step name expected %s actual %s",
			StepCreateSecurityGroups, s.Name())
	}
}

func TestCreateSecurityGroupsStep_Rollback(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected value of err %v", err)
	}
}

func TestCreateSecurityGroupsStep_Description(t *testing.T) {
	s := &CreateSecurityGroupsStep{}

	if desc := s.Description(); desc != "Create security groups" {
		t.Errorf("Wrong step desc expected Create security groups"+
			" actual %s", s.Description())
	}
}
