package amazon

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockKeySvc struct {
	mock.Mock
}

func (m *mockKeySvc) DeleteKeyPair(input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DeleteKeyPairOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteKeyPair_Run(t *testing.T) {
	testCases := []struct {
		decription    string
		existingKeyID string
		getSvcErr     error
		deleteErr     error
		errMsg        string
	}{
		{
			decription: "skip deletion",
		},
		{
			decription:    "get service error",
			existingKeyID: "1234",
			getSvcErr:     errors.New("message1"),
			errMsg:        "message1",
		},
		{
			decription:    "delete error",
			existingKeyID: "1234",
			deleteErr:     errors.New("message2"),
		},
		{
			decription:    "success",
			existingKeyID: "1234",
		},
	}

	deleteKeyPairTimeout = time.Nanosecond * 1
	deleteKeyPairAttemptCount = 1

	for _, testCase := range testCases {
		t.Log(testCase.decription)
		svc := &mockKeySvc{}
		svc.On("DeleteKeyPair", mock.Anything).
			Return(mock.Anything, testCase.deleteErr)
		fn := func(config steps.AWSConfig) (KeyService, error) {
			return svc, testCase.getSvcErr
		}

		step := DeleteKeyPair{
			getSvc: fn,
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				KeyID:       testCase.existingKeyID,
				KeyPairName: "test",
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if testCase.errMsg != "" && err == nil {
			t.Errorf("Error must  not be nil")
			continue
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Err message must contain test %s actual %s",
				testCase.errMsg, err.Error())
		}
	}
}

func TestInitDeleteKeyPair(t *testing.T) {
	InitDeleteKeyPair(GetEC2)

	s := steps.GetStep(DeleteKeyPairStepName)

	if s == nil {
		t.Error("step must not be nil")
	}
}

func TestNewDeleteKeyPair(t *testing.T) {
	step := NewDeleteKeyPairStep(GetEC2)

	if step == nil {
		t.Error("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Error("step get service func must  not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestNewDeleteKeyPairError(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewDeleteKeyPairStep(fn)

	if step == nil {
		t.Error("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Error("step get service func must  not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestDeleteKeyPair_Depends(t *testing.T) {
	s := &DeleteKeyPair{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Dependencies must be nil")
	}
}

func TestDeleteKeyPair_Description(t *testing.T) {
	s := &DeleteKeyPair{}

	if desc := s.Description(); desc != "Delete key pair" {
		t.Errorf("Description must be Delete key pair actual %s",
			s.Description())
	}
}

func TestDeleteKeyPair_Name(t *testing.T) {
	s := &DeleteKeyPair{}

	if name := s.Name(); name != DeleteKeyPairStepName {
		t.Errorf("Name must be %s actual %s",
			DeleteKeyPairStepName, s.Name())
	}
}

func TestDeleteKeyPair_Rollback(t *testing.T) {
	s := &DeleteKeyPair{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error while rollback %v", err)
	}
}
