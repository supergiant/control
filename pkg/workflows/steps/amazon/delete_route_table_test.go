package amazon

import (
	"github.com/stretchr/testify/mock"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"context"
	"bytes"
	"strings"
	"time"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockDeleteRouteTableService struct{
	mock.Mock
}

func (m *mockDeleteRouteTableService) DeleteRouteTable(input *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DeleteRouteTableOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDeleteRouteTable_Run(t *testing.T) {
	testCases := []struct{
		description string
		existingID string
		getSvcErr error
		deleteErr error
		errMsg string
	}{
		{
			description: "skip step",
		},
		{
			description:"get service error",
			existingID: "1234",
			getSvcErr: errors.New("message1"),
			errMsg: "message1",
		},
		{
			description: "delete error",
			existingID: "1234",
			deleteErr: errors.New("message2"),
		},
		{
			description: "success",
			existingID: "1234",
		},
	}

	deleteRouteTimeout = time.Nanosecond
	deleteRouteAttemptCount = 1

	for _, testCase := range testCases {
		svc := &mockDeleteRouteTableService{}
		svc.On("DeleteRouteTable", mock.Anything).
			Return(mock.Anything, testCase.deleteErr)

		step := &DeleteRouteTable{
			getSvc: func(config steps.AWSConfig) (deleteRouteTableSvc, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				RouteTableID: testCase.existingID,
			},
		}
		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Message %s not foun in error %s",
				testCase.errMsg, err.Error())
		}
	}
}

func TestInitDeleteRouteTable(t *testing.T) {
	InitDeleteRouteTable(GetEC2)

	s := steps.GetStep(DeleteRouteTableStepName)

	if s ==  nil {
		t.Errorf("Step %s must  not be nil", DeleteRouteTableStepName)
	}
}

func TestNewDeleteRouteTableStep(t *testing.T) {
	step := NewDeleteRouteTableStep(GetEC2)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Errorf("Step get SVC func must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Wrong values %v %v", err, api)
	}
}

func TestNewDeleteRouteTableStepError(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	step := NewDeleteRouteTableStep(fn)

	if step == nil {
		t.Errorf("Step must not be nil")
	}

	if step.getSvc == nil {
		t.Errorf("Step get SVC func must not be nil")
	}

	if api, err := step.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Wrong values %v %v", err, api)
	}
}

func TestDeleteRouteTable_Depends(t *testing.T) {
	step := &DeleteRouteTable{}

	if deps := step.Depends(); deps == nil || len(deps) != 1 || deps[0] != DeleteSubnetsStepName {
		t.Errorf("Wrong deps %v expected %v", deps, []string{DeleteSubnetsStepName})
	}
}

func TestDeleteRouteTable_Name(t *testing.T) {
	step := &DeleteRouteTable{}

	if name := step.Name(); name != DeleteRouteTableStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteRouteTableStepName, step.Name())
	}
}

func TestDeleteRouteTable_Rollback(t *testing.T) {
	step := &DeleteRouteTable{}

	if err := step.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v while rolling back", err)
	}
}

func TestDeleteRouteTable_Description(t *testing.T) {
	step := &DeleteRouteTable{}

	if desc := step.Description(); desc != "Delete route table from vpc" {
		t.Errorf("Wrong description expected Delete route " +
			"table from vpc actual %s", desc)
	}
}
