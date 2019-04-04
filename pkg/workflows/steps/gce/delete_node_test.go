package gce

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestDeleteNodeStep_Run(t *testing.T) {
	testCases := []struct {
		description string
		getSvcErr   error
		deleteErr   error
		role        model.Role
		errMsg      string
	}{
		{
			description: "delete service",
			getSvcErr:   errors.New("error1"),
			errMsg:      "error1",
		},
		{
			description: "delete master error",
			deleteErr:   errors.New("error2"),
			role:        model.RoleMaster,
			errMsg:      "error2",
		},

		{
			description: "delete node error",
			deleteErr:   errors.New("error3"),
			role:        model.RoleNode,
			errMsg:      "error3",
		},
		{
			description: "success",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		step := &DeleteNodeStep{
			getComputeSvc: func(context.Context, steps.GCEConfig) (*computeService, error) {
				return &computeService{
					deleteInstance: func(string, string, string) (*compute.Operation, error) {
						return nil, testCase.deleteErr
					},
				}, testCase.getSvcErr
			},
		}

		config := &steps.Config{
			Masters: steps.NewMap(map[string]*model.Machine{}),
			Nodes:   steps.NewMap(map[string]*model.Machine{}),
		}

		if testCase.role == model.RoleMaster {
			config.AddMaster(&model.Machine{
				Name: "name",
				Role: testCase.role,
			})
		} else {
			config.AddNode(&model.Machine{
				Name: "name",
				Role: testCase.role,
			})
		}

		err := step.Run(context.Background(), &bytes.Buffer{}, config)

		if err == nil && testCase.errMsg != "" {
			t.Errorf("Error must not be nil")
		}

		if err != nil && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("Error message %s does not contain %s",
				err.Error(), testCase.errMsg)
		}
	}
}

func TestNewDeleteNodeStep(t *testing.T) {
	s := NewDeleteNodeStep()

	if s == nil {
		t.Error("Step must not be nil")
	}

	if s.getComputeSvc == nil {
		t.Errorf("get compute service must not be nil")
	}

	if client, err := s.getComputeSvc(context.Background(), steps.GCEConfig{
		ServiceAccount: steps.ServiceAccount{
			Type: "service_account",
		},
	});

	client == nil || err != nil {
		t.Errorf("Unexpected values %v %v", client, err)
	}
}

func TestDeleteNodeStep_Depends(t *testing.T) {
	s := DeleteNodeStep{}

	if deps := s.Depends(); deps != nil {
		t.Errorf("Dependencies must be nil")
	}
}

func TestDeleteNodeStep_Name(t *testing.T) {
	s := &DeleteNodeStep{}

	if name := s.Name(); name != DeleteNodeStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteNodeStepName, name)
	}
}

func TestDeleteNodeStep_Rollback(t *testing.T) {
	s := &DeleteNodeStep{}

	if err := s.Rollback(context.Background(), &bytes.Buffer{}, &steps.Config{}); err != nil {
		t.Errorf("Unexpected error when rollback %v", err)
	}
}

func TestDeleteNodeStep_Description(t *testing.T) {
	s := &DeleteNodeStep{}

	if desc := s.Description(); desc != "Google compute engine delete instance step" {
		t.Errorf("Wrong description expected "+
			"Google compute engine delete instance step actual %s", desc)
	}
}
