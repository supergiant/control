package amazon

import (
	"testing"

	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/control/pkg/workflows/steps"
	"strings"
)

type FakeEC2KeyPair struct {
	ec2iface.EC2API

	importOutput   *ec2.ImportKeyPairOutput
	describeOutput *ec2.DescribeKeyPairsOutput

	importErr   error
	describeErr error
}

type mockKeyPairSvc struct {
	mock.Mock
}

func (m *mockKeyPairSvc) ImportKeyPairWithContext(ctx aws.Context,
	req *ec2.ImportKeyPairInput, opts ...request.Option) (*ec2.ImportKeyPairOutput, error) {
	args := m.Called(ctx, req, opts)
	val, ok := args.Get(0).(*ec2.ImportKeyPairOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockKeyPairSvc) WaitUntilKeyPairExists(req *ec2.DescribeKeyPairsInput) error {
	args := m.Called(req)
	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func TestImportKeyPair_Run(t *testing.T) {
	testCases := []struct {
		description string
		getSvcErr   error
		clusterId   string

		importOut *ec2.ImportKeyPairOutput
		importErr error
		waitErr   error
		errMsg    string
	}{
		{
			description: "get service error",
			getSvcErr:   errors.New("message1"),
			errMsg:      "message1",
		},
		{
			description: "cluster id is to short",
			clusterId:   "124",
			errMsg:      "too short",
		},
		{
			description: "import error",
			clusterId:   "12345678",
			importErr:   errors.New("message2"),
			errMsg:      "message2",
		},
		{
			description: "wait error",
			clusterId:   "12345678",
			importOut: &ec2.ImportKeyPairOutput{
				KeyFingerprint: aws.String("fingerprint"),
				KeyName:        aws.String("keyName"),
			},
			waitErr: errors.New("message3"),
			errMsg:  "message3",
		},
		{
			description: "success",
			clusterId:   "12345678",
			importOut: &ec2.ImportKeyPairOutput{
				KeyFingerprint: aws.String("fingerprint"),
				KeyName:        aws.String("keyName"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockKeyPairSvc{}
		svc.On("ImportKeyPairWithContext",
			mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.importOut, testCase.importErr)
		svc.On("WaitUntilKeyPairExists",
			mock.Anything).
			Return(testCase.waitErr)

		config := &steps.Config{
			ClusterName: "test",
			ClusterID:   testCase.clusterId,
			AWSConfig:   steps.AWSConfig{},
		}

		step := KeyPairStep{
			getSvc: func(steps.AWSConfig) (keyImporter, error) {
				return svc, testCase.getSvcErr
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s doesnt not contain %s",
				err.Error(), testCase.errMsg)
		}
	}
}

func TestInitImportKeyPair(t *testing.T) {
	InitImportKeyPair(GetEC2)

	s := steps.GetStep(StepImportKeyPair)

	if s == nil {
		t.Error("Step must not be nil")
	}
}

func TestNewImportKeyPairStep(t *testing.T) {
	s := NewImportKeyPairStep(GetEC2)

	if s == nil {
		t.Error("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}


func TestNewImportKeyPairStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewImportKeyPairStep(fn)

	if s == nil {
		t.Error("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}

func TestKeyPairStep_Depends(t *testing.T) {
	s := &KeyPairStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("dependencies must be nil")
	}
}

func TestKeyPairStep_Description(t *testing.T) {
	s := &KeyPairStep{}

	if desc := s.Description(); desc != "If no keypair is present in config, creates a new keypair" {
		t.Errorf("Description is wrong expected " +
			"If no keypair is present in config, creates a new keypair actuak %s",
			desc)
	}
}

func TestKeyPairStep_Rollback(t *testing.T) {
	s := &KeyPairStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error when rollback %v", err)
	}
}

func TestKeyPairStep_Name(t *testing.T) {
	s := &KeyPairStep{}

	if name := s.Name(); name != StepImportKeyPair {
		t.Errorf("Wrong name expected %s actual %s",
			StepImportKeyPair, name)
	}
}
