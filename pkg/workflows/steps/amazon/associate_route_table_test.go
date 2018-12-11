package amazon


import (
	"testing"
	"context"
	"bytes"


	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

type mockAssociater struct{
	association *ec2.AssociateRouteTableOutput
	err error
}

func (m *mockAssociater) AssociateRouteTable(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	return m.association, m.err
}

func TestNewAssociateRouteTableStep(t *testing.T) {
	step := NewAssociateRouteTableStep(GetEC2)

	if step.getRouteTableSvc == nil {
		t.Error("getRouteTableSvc must not be nil")
	}
}

func TestAssociateRouteTableStep_Depends(t *testing.T) {
	step := &AssociateRouteTableStep{}

	if deps := step.Depends(); len(deps) != 1 || deps[0] != StepCreateSubnets {
		t.Errorf("Wron dependencies expected %v actual %v",
			[]string{StepCreateSubnets}, deps)
	}
}
 func TestAssociateRouteTableStep_Name(t *testing.T) {
	 step := &AssociateRouteTableStep{}

	 if step.Name() != StepAssociateRouteTable {
	 	t.Errorf("Wrong name expected %s actual %s",
	 		StepAssociateRouteTable, step.Name())
	 }
}

func TestAssociateRouteTableStep_Rollback(t *testing.T) {
	step := &AssociateRouteTableStep{}

	if err := step.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error %v when roll back", err)
	}
}

func TestInitAssociateRouteTable(t *testing.T) {
	InitAssociateRouteTable(nil)

	s := steps.GetStep(StepAssociateRouteTable)

	if s == nil {
		t.Errorf("Step %s not found", StepAssociateRouteTable)
	}
}

func TestAssociateRouteTableStep_Description(t *testing.T) {
	step := &AssociateRouteTableStep{}

	if step.Description() != "Associate route table with all subnets in VPC" {
		t.Errorf("Wrong description expected Associate route table " +
			"with all subnets in VPC actual %s", step.Description())
	}
}


func TestAssociateRouteTableStep_Run(t *testing.T) {
	testCases := []struct{
		description string
		getSvcErr error
		associationID string
		associateErr error
		errMsg string
	}{
		{
			description: "get svc error",
			getSvcErr: errors.New("message1"),
			errMsg: "message1",
		},
		{
			description: "create association error",
			associateErr: errors.New("message2"),
			errMsg: "message2",
		},
		{
			description: "success",
			associationID: "1234",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &mockAssociater{
			association: &ec2.AssociateRouteTableOutput{
				AssociationId: aws.String(testCase.associationID),
			},
			err:testCase.associateErr,
		}

		step := &AssociateRouteTableStep{
			getRouteTableSvc: func(cfg steps.AWSConfig) (Associater, error) {
				return svc, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			AWSConfig: steps.AWSConfig{
				Subnets: map[string]string{
					"az1": "subnet1",
					"az2": "subnet2",
					"az3": "subnet3",
				},
			},
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if testCase.errMsg == "" && err != nil {
			t.Errorf("Unexpected error %v", err)
			continue
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s does not contain %s",
				err.Error(), testCase.errMsg)
			continue
		}

		if testCase.errMsg == "" &&
			len(config.AWSConfig.RouteTableAssociationIDs) != len(config.AWSConfig.Subnets) {
			t.Errorf("Route table must be size %d of one actual %d",
				len(config.AWSConfig.Subnets), len(config.AWSConfig.RouteTableAssociationIDs))
		}
	}
}