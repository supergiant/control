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
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockDisassociateService struct{
	mock.Mock
}

func (m *mockDisassociateService) DisassociateRouteTable(
	input *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	args := m.Called(input)
	val, ok := args.Get(0).(*ec2.DisassociateRouteTableOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestDisassociateRouteTable_Run(t *testing.T) {
	testCases := []struct{
		description string
		existingID string
		getSvcErr error
		disassociateErr error
		errMsg string
	}{
		{
			description: "skip disassociate",
		},
		{
			description: "get service error",
			existingID: "1234",
			getSvcErr: errors.New("message1"),
			errMsg: "message1",
		},
		{
			description: "disassociate error",
			existingID: "1234",
			disassociateErr: errors.New("message2"),
		},
		{
			description: "success",
			existingID: "1234",
		},
	}

	for _, testCase := range testCases {
		svc := &mockDisassociateService{}
		svc.On("DisassociateRouteTable",
			mock.Anything).Return(mock.Anything, testCase.disassociateErr)

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				RouteTableAssociationIDs: map[string]string{
					"subnetId": "routeTableId",
				},
			},
		}

		step := DisassociateRouteTable{
			getSvc: func(cfg steps.AWSConfig) (DisassociateService, error) {
				return svc, testCase.getSvcErr
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s must contain %s", err.Error(),
				testCase.errMsg)
		}
	}
}


func TestInitDisassociateRouteTable(t *testing.T) {
	InitDisassociateRouteTable(GetEC2)

	s := steps.GetStep(DisassociateRouteTableStepName)

	if s == nil {
		t.Errorf("Step must not be nil")
	}
}

func TestNewDisassociateRouteTableStep(t *testing.T) {
	s := NewDisassociateRouteTableStep(GetEC2)

	if s == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if s.getSvc == nil {
		t.Errorf("Get service for %s must not be nil",
			DisassociateRouteTableStepName)
		return
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err != nil || api == nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}

func TestNewDisassociateRouteTableStepErr(t *testing.T) {
	fn := func(steps.AWSConfig) (ec2iface.EC2API, error) {
		return nil, errors.New("errorMessage")
	}

	s := NewDisassociateRouteTableStep(fn)

	if s == nil {
		t.Errorf("Step must not be nil")
		return
	}

	if s.getSvc == nil {
		t.Errorf("Get service for %s must not be nil",
			DisassociateRouteTableStepName)
		return
	}

	if api, err := s.getSvc(steps.AWSConfig{}); err == nil || api != nil {
		t.Errorf("Wrong values %v %v", api, err)
	}
}

func TestDisassociateRouteTable_Depends(t *testing.T) {
	s := &DisassociateRouteTable{}

	if deps := s.Depends(); deps == nil || len(deps) != 1 || deps[0] != DeleteSecurityGroupsStepName {
		t.Errorf("Wrong dependencies expected %v actual %v",
			[]string{DeleteSecurityGroupsStepName}, deps)
	}
}

func TestDisassociateRouteTable_Name(t *testing.T) {
	s := &DisassociateRouteTable{}

	if name := s.Name(); name != DisassociateRouteTableStepName {
		t.Errorf("Wrong name expected %s actual %s",
			DisassociateRouteTableStepName, name)
	}
}

func TestDisassociateRouteTable_Rollback(t *testing.T) {
	s := &DisassociateRouteTable{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error while roll back %v", err)
	}
}

func TestDisassociateRouteTable_Description(t *testing.T) {
	s :=  &DisassociateRouteTable{}

	if desc := s.Description(); desc != "Disassociate route table with subnets" {
		t.Errorf("Wrong description expected Disassociate route table " +
			"with subnets actual %s", desc)
	}
}
