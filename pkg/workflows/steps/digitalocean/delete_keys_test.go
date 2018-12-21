package digitalocean

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"

	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	testKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDB2ckfv5rVySSq7p9ziEt+waU28aFGo9VNGr9gottC7dew2N+ggLj7DzUUAEI2809qPBxNFN9C/rC2aP+brS8jcvInbcMxOHK/QzxzOSDjQQOfq5tQ451HshkCqRFtz5cIRgrn/yLaPZ+4dr+gspsgu8qvTGZIb8zCyjVPZsfhg70Z8Ql+1kn+1KTljOlvQ6jlxZvZX3o68kMb8wRvkFc8ps4xTyeCfHaCqz6OHWnV9DCtvQYmMmADzezJKOvwAeR6Uf1A1Lwe+B8eUvxtfaeYUZ5pWtHFFfOykmd03Xk0pRYAwtSC9ZWeje6WooyTMf56ErpIUK4qgXmJzG2oHHjD`
)

type mockKeyDeleter struct {
	responses []*godo.Response
	errors    []error
	callCount int
}

func (m *mockKeyDeleter) DeleteByFingerprint(ctx context.Context, fg string) (*godo.Response, error) {
	defer func() { m.callCount++ }()

	if len(m.responses) <= m.callCount || len(m.errors) <= m.callCount {
		panic("illegal call")
	}

	return m.responses[m.callCount], m.errors[m.callCount]
}

func TestDeleteKeysStep_Run(t *testing.T) {
	testCases := []struct {
		responses         []*godo.Response
		errors            []error
		expectedCallCount int
	}{
		{
			responses: []*godo.Response{
				{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
			},
			errors:            []error{sgerrors.ErrNotFound, sgerrors.ErrNotFound},
			expectedCallCount: 2,
		},
		{
			responses: []*godo.Response{
				{
					Response: &http.Response{
						StatusCode: http.StatusAccepted,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
			},
			errors:            []error{nil, sgerrors.ErrNotFound},
			expectedCallCount: 2,
		},
		{
			responses: []*godo.Response{
				{
					Response: &http.Response{
						StatusCode: http.StatusAccepted,
					},
				},
				{
					Response: &http.Response{
						StatusCode: http.StatusAccepted,
					},
				},
			},
			errors:            []error{nil, nil},
			expectedCallCount: 2,
		},
	}

	for _, testCase := range testCases {
		keyDeleterMock := &mockKeyDeleter{
			responses: testCase.responses,
			errors:    testCase.errors,
			callCount: 0,
		}

		step := &DeleteKeysStep{
			getKeyService: func(string) keyDeleter {
				return keyDeleterMock
			},
		}

		config := &steps.Config{
			SshConfig: steps.SshConfig{
				PublicKey:          testKey,
				BootstrapPublicKey: testKey,
			},
		}
		step.Run(context.Background(), ioutil.Discard, config)

		if keyDeleterMock.callCount != testCase.expectedCallCount {
			t.Errorf("Expected call count %d actual %d",
				testCase.expectedCallCount, keyDeleterMock.callCount)
		}
	}
}

func TestNewDeleteKeysStep(t *testing.T) {
	step := NewDeleteKeysStep()

	if step == nil {
		t.Errorf("Step value must not be nil")
		return
	}

	if step.getKeyService == nil {
		t.Errorf("Step value must not be nil")
	}

	if svc := step.getKeyService("token"); svc == nil {
		t.Errorf("key service must be nil")
	}
}

func TestDeleteKeysStep_Rollback(t *testing.T) {
	step := &DeleteKeysStep{}

	if err := step.Rollback(context.Background(), ioutil.Discard, nil); err != nil {
		t.Errorf("Unexpected error value %v", err)
	}
}

func TestDeleteKeysStep_Name(t *testing.T) {
	step := NewDeleteKeysStep()

	if name := step.Name(); name != DeleteDeleteKeysStepName {
		t.Errorf("Wrong step name expected %s actual %s",
			DeleteDeleteKeysStepName, name)
	}
}

func TestDeleteKeysStep_Depends(t *testing.T) {
	step := &DeleteKeysStep{}

	if deps := step.Depends(); deps != nil {
		t.Errorf("Unexpected deps value %v", deps)
	}
}

func TestDeleteKeysStep_Description(t *testing.T) {
	step := &DeleteKeysStep{}

	if desc := step.Description(); desc == "" {
		t.Errorf("Description must not be empty")
	}
}
