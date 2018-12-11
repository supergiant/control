package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockDeleter struct {
	mock.Mock
}

func (m *mockDeleter) DetachInternetGateway(input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DetachInternetGatewayOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockDeleter) DeleteInternetGateway(input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DeleteInternetGatewayOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteInternetGateway_Run(t *testing.T) {
	testCases := []struct {
		description string
		existingID  string
		getSVCErr   error
		detachErr   error
		deleteErr   error
		errMsg      string
	}{
		{
			description: "skip deleting IGW",
		},
		{
			description: "get service error",
			existingID:  "1234",
			getSVCErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description: "detach error",
			existingID:  "1234",
			detachErr:   errors.New("message2"),
		},

		{
			description: "delete error",
			existingID:  "1234",
			deleteErr:   errors.New("message2"),
		},

		{
			description: "success",
			existingID:  "1234",
		},
	}

	for _, testCase := range testCases {
		svc := &mockDeleter{}
		svc.On("DetachInternetGateway",
			mock.Anything).Return(mock.Anything, testCase.detachErr)
		svc.On("DeleteInternetGateway",
			mock.Anything).Return(mock.Anything, testCase.deleteErr)
		step := DeleteInternetGateway{
			getIGWService: func(cfg steps.AWSConfig) (IGWDeleter, error) {
				return svc, testCase.getSVCErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				InternetGatewayID: testCase.existingID,
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Wrong error message expected %s actual %s",
				testCase.errMsg, err.Error())
		}
	}
}

func TestInitDeleteInternetGateWay(t *testing.T) {
	InitDeleteInternetGateWay(GetEC2)

	s := steps.GetStep(DeleteInternetGatewayStepName)

	if s == nil {
		t.Error("Step must not be nil")
	}
}

func TestDeleteInternetGateway_Depends(t *testing.T) {
	s := &DeleteInternetGateway{}

	if deps := s.Depends(); deps == nil ||
		len(deps) != 1 || deps[0] != DeleteSecurityGroupsStepName {
		t.Errorf("Wrong dependencis expected %v actual %v",
			[]string{DeleteSecurityGroupsStepName}, s.Depends())
	}
}

func TestDeleteInternetGateway_Name(t *testing.T) {
	s := &DeleteInternetGateway{}

	if name := s.Name(); name != DeleteInternetGatewayStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteInternetGatewayStepName, s.Name())
	}
}

func TestDeleteInternetGateway_Description(t *testing.T) {
	s := &DeleteInternetGateway{}

	if desc := s.Description(); desc != "Delete internet gateway from VPC" {
		t.Errorf("Wrong step desc expected Delete internet gateway from VPC actual %s",
			s.Description())
	}
}

func TestDeleteInternetGateway_Rollback(t *testing.T) {
	s := &DeleteInternetGateway{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpecter error %v", err)
	}
}

func TestNewDeleteInernetGateway(t *testing.T) {
	step := NewDeleteInernetGateway(GetEC2)

	if step == nil {
		t.Error("step must not be nil")
	}

	if step.getIGWService == nil {
		t.Errorf("getIGWService must not be nil")
	}

	if api, err := step.getIGWService(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewDeleteInernetGatewayErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewDeleteInernetGateway(fn)

	if step == nil {
		t.Error("step must not be nil")
	}

	if step.getIGWService == nil {
		t.Errorf("getIGWService must not be nil")
	}

	if api, err := step.getIGWService(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}
