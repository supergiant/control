package digitalocean

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestDeleteMachineRun(t *testing.T) {
	testCases := []struct {
		description   string
		machineName   string
		dropletErrors []error
		responses     []*godo.Response
	}{
		{
			description:   "empty tag error",
			dropletErrors: []error{errors.New(""), errors.New(""), errors.New("")},
			responses:     []*godo.Response{nil, nil, nil},
		},
		{
			description:   "",
			machineName:   "test",
			dropletErrors: []error{errors.New(""), errors.New(""), errors.New("")},
			responses: []*godo.Response{
				{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
					},
				},
			},
		},
		{
			description:   "",
			machineName:   "test",
			dropletErrors: []error{errors.New(""), errors.New(""), nil},
			responses: []*godo.Response{
				{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusNoContent,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		for i := 0; i < 3; i++ {
			svc := new(mockDeleteService)
			svc.On("DeleteByTag", mock.Anything, testCase.machineName).
				Return(testCase.responses[i], testCase.dropletErrors[i])

			step := NewDeleteMachineStep(time.Microsecond * 1)
			step.getDeleteService = func(string) DeleteService {
				return svc
			}
			err := step.Run(context.Background(), &bytes.Buffer{}, &steps.Config{
				Node: node.Node{
					Name: testCase.machineName,
				},
			})

			if err != testCase.dropletErrors[i] {
				t.Errorf("Wrong error expected %v actual %v", testCase.dropletErrors[i], err)
			}
		}
	}
}

func TestStepDeleteMachineName(t *testing.T) {
	s := DeleteMachineStep{}

	if s.Name() != DeleteMachineStepName {
		t.Errorf("Unexpected step name expected %s actual %s", DeleteMachineStepName, s.Name())
	}
}

func TestDeleteMachineDepends(t *testing.T) {
	s := DeleteMachineStep{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}

func TestStep_Rollback(t *testing.T) {
	s := DeleteMachineStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
