package amazon

import (
	"context"
	"bytes"
	"strings"
	"time"
	"testing"

	"github.com/pkg/errors"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockDeleteVpcSvc struct{
	mock.Mock
}

func (m *mockDeleteVpcSvc) DeleteVpcWithContext(ctx aws.Context, input *ec2.DeleteVpcInput, opts ...request.Option) (*ec2.DeleteVpcOutput, error) {
	args := m.Called(ctx, input, opts)
	val, ok := args.Get(0).(*ec2.DeleteVpcOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}


func TestDeleteVPC_Run(t *testing.T) {
	testCases := []struct{
		description string
		existingID string

		getSvcErr error
		deleteErr error

		errMsg string
	}{
		{
			description: "skip delete",
		},
		{
			description: "get svc error",
			existingID: "1234",
			getSvcErr: errors.New("message1"),
			errMsg: "message1",
		},
		{
			description: "delete error",
			existingID: "1234",
			deleteErr: errors.New("message2"),
			errMsg: "message2",
		},
		{
			description: "success",
			existingID: "1234",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockDeleteVpcSvc{}
		svc.On("DeleteVpcWithContext", mock.Anything,
			mock.Anything, mock.Anything).Return(mock.Anything,
				testCase.deleteErr)
		step := &DeleteVPC{
			getSvc: func(config steps.AWSConfig) (vpcSvc, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				VPCID: testCase.existingID,
			},
		}

		deleteVPCAttemptCount  =1
		deleteVPCTimeout = time.Nanosecond

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %v must contain %s",
				err, testCase.errMsg)
		}
	}
}

func TestInitDeleteVPC(t *testing.T) {
	InitDeleteVPC(GetEC2)

	s := steps.GetStep(DeleteVPCStepName)

	if s == nil {
		t.Error("step must not be nil")
	}
}

func TestNewDeleteVPC(t *testing.T) {
	s := NewDeleteVPC(GetEC2)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", err, api)
	}
}


func TestNewDeleteVPCErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDeleteVPC(fn)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", err, api)
	}
}

func TestDeleteVPC_Name(t *testing.T) {
	s := DeleteVPC{}

	if name := s.Name(); name != DeleteVPCStepName {
		t.Errorf("Wrong name expected %s actual %s",
			DeleteVPCStepName, name)
	}
}

func TestDeleteVPC_Depends(t *testing.T) {
	s := DeleteVPC{}

	if deps := s.Depends(); deps == nil ||
		len(deps) != 1 || deps[0] != DeleteSecurityGroupsStepName {
		t.Errorf("Wrong dependencies expected %v actual %v",
			deps, []string{DeleteSecurityGroupsStepName})
	}
}

func TestDeleteVPC_Rollback(t *testing.T) {
	s := DeleteVPC{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestDeleteVPC_Description(t *testing.T) {
	s := DeleteVPC{}

	if desc := s.Description(); desc != "Delete vpc" {
		t.Errorf("Wrong description %s expected Delete vpc", desc)
	}
}