package digitalocean

import (
	"net/http"
	"testing"

	"bytes"
	"context"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"time"
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
		clusterName   string
		dropletErrors []error
		responses     []*godo.Response
	}{
		{
			description:   "empty tag",
			clusterName:   "",
			dropletErrors: []error{errors.New(""), errors.New(""), errors.New("")},
			responses:     []*godo.Response{nil, nil, nil},
		},
		{
			description:   "retry exceeded",
			clusterName:   "fail",
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
			clusterName:   "success",
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

			step := NewDeleteClusterStep(time.Microsecond * 1)
			step.getDeleteService = func(string) DeleteService {
				return svc
			}
			err := step.Run(context.Background(), &bytes.Buffer{}, &steps.Config{
				ClusterName: testCase.clusterName,
			})

			if err != testCase.dropletErrors[i] {
				t.Errorf("Wrong error expected %v actual %v", testCase.dropletErrors[i], err)
			}
		}
	}
}
