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

type mockDeleteSubnetService struct {
	mock.Mock
}

func (m *mockDeleteSubnetService) DeleteSubnet(
	input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DeleteSubnetOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteSubnets_Run(t *testing.T) {
	testCases := []struct {
		description string
		subnets     map[string]string

		existingID string

		getSvcErr error
		deleteErr error

		errMsg string
	}{
		{
			description: "skip empty",
		},
		{
			description: "get svc error",
			subnets: map[string]string{
				"az1": "subnet1",
			},
			existingID: "1234",
			getSvcErr:  errors.New("message1"),
			errMsg:     "message1",
		},
		{
			description: "delete error",
			subnets: map[string]string{
				"az1": "subnet1",
			},
			existingID: "1234",
			deleteErr:  errors.New("message2"),
		},
		{
			description: "success",
			subnets: map[string]string{
				"az1": "subnet1",
			},
			existingID: "1234",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockDeleteSubnetService{}
		svc.On("DeleteSubnet", mock.Anything).
			Return(mock.Anything, testCase.deleteErr)

		step := &DeleteSubnets{
			getSvc: func(config steps.AWSConfig) (deleteSubnetesSvc, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				Subnets: testCase.subnets,
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s does not have expected %s",
				err.Error(), testCase.errMsg)
		}
	}
}

func TestInitDeleteSubnets(t *testing.T) {
	InitDeleteSubnets(GetEC2)

	s := steps.GetStep(DeleteSubnetsStepName)

	if s == nil {
		t.Errorf("step must not be nil")
	}
}

func TestNewDeleteSubnets(t *testing.T) {
	s := NewDeleteSubnets(GetEC2)

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

func TestNewDeleteSubnetsErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDeleteSubnets(fn)

	if s == nil {
		t.Errorf("Step must not be nil")
	}

	if s.getSvc == nil {
		t.Errorf("getSvc must not be nil")
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Unexpected values %v %v", api, err)
	}
}

func TestDeleteSubnets_Depends(t *testing.T) {
	s := &DeleteSubnets{}

	if deps := s.Depends(); deps == nil ||
		len(deps) != 1 || deps[0] != DeleteSecurityGroupsStepName {
		t.Errorf("Wrong dependencies expected %v actual %v",
			[]string{DeleteSecurityGroupsStepName}, deps)
	}
}

func TestDeleteSubnets_Description(t *testing.T) {
	s := &DeleteSubnets{}

	if desc := s.Description(); desc != "Deletes security groups" {
		t.Errorf("Wrong description expected "+
			"Deletes security groups actual %s", desc)
	}
}

func TestDeleteSubnets_Name(t *testing.T) {
	s := &DeleteSubnets{}

	if name := s.Name(); name != DeleteSubnetsStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteSubnetsStepName, name)
	}
}

func TestDeleteSubnets_Rollback(t *testing.T) {
	s := &DeleteSubnets{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("unexpected error when rollbacl %v", err)
	}
}
