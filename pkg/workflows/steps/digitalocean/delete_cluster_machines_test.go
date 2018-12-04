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

	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockDeleteService struct {
	mock.Mock
}

func (m *mockDeleteService) DeleteByTag(ctx context.Context, tag string) (*godo.Response, error) {
	args := m.Called(ctx, tag)
	val, ok := args.Get(0).(*godo.Response)

	if !ok {
		return nil, args.Error(1)
	}

	return val, args.Error(1)
}

func TestDeleteClusterRun(t *testing.T) {
	testCases := []struct {
		description   string
		clusterID     string
		dropletErrors []error
		responses     []*godo.Response
	}{
		{
			description:   "empty tag",
			clusterID:     "",
			dropletErrors: []error{errors.New(""), errors.New(""), errors.New("")},
			responses:     []*godo.Response{nil, nil, nil},
		},
		{
			description:   "retry exceeded",
			clusterID:     "fail",
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
			description:   "success",
			clusterID:     "success",
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
			svc.On("DeleteByTag", mock.Anything, mock.Anything).
				Return(testCase.responses[i], testCase.dropletErrors[i])

			step := NewDeletemachinesStep(time.Microsecond * 1)
			step.getDeleteService = func(string) DeleteService {
				return svc
			}
			err := step.Run(context.Background(), &bytes.Buffer{}, &steps.Config{
				ClusterID: testCase.clusterID,
			})

			if err != testCase.dropletErrors[i] {
				t.Errorf("Wrong error expected %v actual %v", testCase.dropletErrors[i], err)
			}
		}
	}
}

func TestDeleteClusterStepName(t *testing.T) {
	s := DeleteMachinesStep{}

	if s.Name() != DeleteClusterMachines {
		t.Errorf("Unexpected step name expected %s actual %s", DeleteClusterMachines, s.Name())
	}
}

func TestDeleteClusterDepends(t *testing.T) {
	s := DeleteMachinesStep{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}

func TestStepDeleteCluster_Rollback(t *testing.T) {
	s := DeleteMachinesStep{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}
