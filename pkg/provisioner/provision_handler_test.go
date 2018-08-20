package provisioner

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows"
)

type mockProvisioner struct {
	provision func(context.Context, *profile.KubeProfile) ([]*workflows.Task, error)
}

func (m *mockProvisioner) Provision(ctx context.Context, kubeProfile *profile.KubeProfile, credentials model.Credentials) ([]*workflows.Task, error) {
	return m.provision(ctx, kubeProfile)
}

type mockKubeProfileGetter struct {
	get func(context.Context, string) (*profile.KubeProfile, error)
}

func (m *mockKubeProfileGetter) Get(ctx context.Context, id string) (*profile.KubeProfile, error) {
	return m.get(ctx, id)
}

type mockAccountGetter struct {
	get func(context.Context, string) (*model.CloudAccount, error)
}

func (m *mockAccountGetter) Get(ctx context.Context, id string) (*model.CloudAccount, error) {
	return m.get(ctx, id)
}

func TestProvisionHandler(t *testing.T) {
	p := &ProvisionRequest{
		"1234",
		"abcd",
	}

	validBody, _ := json.Marshal(p)

	testCases := []struct {
		expectedCode int

		body       []byte
		getProfile func(context.Context, string) (*profile.KubeProfile, error)
		getAccount func(context.Context, string) (*model.CloudAccount, error)
		provision  func(context.Context, *profile.KubeProfile) ([]*workflows.Task, error)
	}{
		{
			body:         []byte(`{`),
			expectedCode: http.StatusBadRequest,
		},
		{
			body:         validBody,
			expectedCode: http.StatusNotFound,
			getAccount: func(context.Context, string) (*model.CloudAccount, error){
				return nil, nil
			},
			getProfile: func(context.Context, string) (*profile.KubeProfile, error) {
				return nil, sgerrors.ErrNotFound
			},
		},
		{
			body:         validBody,
			expectedCode: http.StatusInternalServerError,
			getAccount: func(context.Context, string) (*model.CloudAccount, error){
				return &model.CloudAccount{}, nil
			},
			getProfile: func(context.Context, string) (*profile.KubeProfile, error) {
				return &profile.KubeProfile{}, nil
			},
			provision: func(context.Context, *profile.KubeProfile) ([]*workflows.Task, error) {
				return nil, sgerrors.ErrInvalidCredentials
			},
		},
		{
			body:         validBody,
			expectedCode: http.StatusAccepted,
			getAccount: func(context.Context, string) (*model.CloudAccount, error){
				return &model.CloudAccount{}, nil
			},
			getProfile: func(context.Context, string) (*profile.KubeProfile, error) {
				return &profile.KubeProfile{}, nil
			},
			provision: func(context.Context, *profile.KubeProfile) ([]*workflows.Task, error) {
				return []*workflows.Task{
					{
						ID: "master-task-id-1",
					},
					{
						ID: "node-task-id-2",
					},
				}, nil
			},
		},
	}

	provisioner := &mockProvisioner{}
	profileGetter := &mockKubeProfileGetter{}
	accGetter := &mockAccountGetter{}

	for _, testCase := range testCases {
		provisioner.provision = testCase.provision
		profileGetter.get = testCase.getProfile
		accGetter.get = testCase.getAccount

		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(testCase.body))
		rec := httptest.NewRecorder()

		handler := ProvisionHandler{
			provisioner:   provisioner,
			profileGetter: profileGetter,
			accountGetter: accGetter,
		}

		handler.Provision(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong status code expected %d actual %d", testCase.expectedCode, rec.Code)
		}

		if testCase.expectedCode == http.StatusAccepted {
			resp := ProvisionResponse{}

			err := json.NewDecoder(rec.Body).Decode(&resp)

			if err != nil {
				t.Errorf("Unepxpected error while decoding response %v", err)
			}
		}
	}
}
