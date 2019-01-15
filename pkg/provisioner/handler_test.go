package provisioner

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type mockProvisioner struct {
	provisionCluster func(context.Context, *profile.Profile, *steps.Config) (map[string][]*workflows.Task, error)
	provisionNode    func(context.Context, profile.NodeProfile, *model.Kube, *steps.Config) (*workflows.Task, error)
}

func (m *mockProvisioner) ProvisionCluster(ctx context.Context, kubeProfile *profile.Profile, config *steps.Config) (map[string][]*workflows.Task, error) {
	return m.provisionCluster(ctx, kubeProfile, config)
}

func (m *mockProvisioner) ProvisionNode(ctx context.Context, nodeProfile profile.NodeProfile, kube *model.Kube, config *steps.Config) (*workflows.Task, error) {
	return m.provisionNode(ctx, nodeProfile, kube, config)
}

type mockAccountGetter struct {
	get func(context.Context, string) (*model.CloudAccount, error)
}

func (m *mockAccountGetter) Get(ctx context.Context, id string) (*model.CloudAccount, error) {
	return m.get(ctx, id)
}

type mockKubeGetter struct {
	get func(context.Context, string) (*model.Kube, error)
}

func (m *mockKubeGetter) Get(ctx context.Context, name string) (*model.Kube, error) {
	return m.get(ctx, name)
}

func TestProvisionBadClusterName(t *testing.T) {
	testCases := []string{"non_Valid`", "_@badClusterName"}

	for _, clusterName := range testCases {
		provisionRequest := ProvisionRequest{
			ClusterName: clusterName,
		}

		bodyBytes, _ := json.Marshal(&provisionRequest)
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		handler := Handler{}
		handler.Provision(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Wrong status code expected %d actual %d", http.StatusBadRequest, rec.Code)
			return
		}
	}
}

func TestProvisionHandler(t *testing.T) {
	p := &ProvisionRequest{
		"test",
		profile.Profile{},
		"1234",
	}

	validBody, _ := json.Marshal(p)

	testCases := []struct {
		description string

		expectedCode int

		body       []byte
		getProfile func(context.Context, string) (*profile.Profile, error)
		kubeGetter func(context.Context, string) (*model.Kube, error)
		getAccount func(context.Context, string) (*model.CloudAccount, error)
		provision  func(context.Context, *profile.Profile, *steps.Config) (map[string][]*workflows.Task, error)
	}{
		{
			description:  "malformed request body",
			body:         []byte(`{`),
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "account not found",
			body:         validBody,
			expectedCode: http.StatusNotFound,
			getAccount: func(context.Context, string) (*model.CloudAccount, error) {
				return nil, sgerrors.ErrNotFound
			},
			kubeGetter: func(context.Context, string) (*model.Kube, error) {
				return nil, nil
			},
		},
		{
			description:  "wrong cloud provider name",
			body:         validBody,
			expectedCode: http.StatusNotFound,
			getAccount: func(context.Context, string) (*model.CloudAccount, error) {
				return &model.CloudAccount{}, nil
			},
			kubeGetter: func(context.Context, string) (*model.Kube, error) {
				return nil, nil
			},
		},
		{
			description:  "invalid credentials when provisionCluster",
			body:         validBody,
			expectedCode: http.StatusInternalServerError,
			getAccount: func(context.Context, string) (*model.CloudAccount, error) {
				return &model.CloudAccount{
					Provider: clouds.DigitalOcean,
				}, nil
			},
			kubeGetter: func(context.Context, string) (*model.Kube, error) {
				return nil, nil
			},
			provision: func(context.Context, *profile.Profile, *steps.Config) (map[string][]*workflows.Task, error) {
				return nil, sgerrors.ErrInvalidCredentials
			},
		},
		{
			body:         validBody,
			expectedCode: http.StatusAccepted,
			getAccount: func(context.Context, string) (*model.CloudAccount, error) {
				return &model.CloudAccount{
					Provider: clouds.DigitalOcean,
				}, nil
			},
			kubeGetter: func(context.Context, string) (*model.Kube, error) {
				return nil, nil
			},
			provision: func(ctx context.Context, profile *profile.Profile, config *steps.Config) (map[string][]*workflows.Task, error) {
				config.ClusterID = uuid.New()
				return map[string][]*workflows.Task{
					"master": {
						{
							ID: "master-task-id-1",
						},
					},
					"node": {
						{
							ID: "node-task-id-2",
						},
					},
					"cluster": {
						{},
					},
				}, nil
			},
		},
	}

	provisioner := &mockProvisioner{}
	kubeGetter := &mockKubeGetter{}
	accGetter := &mockAccountGetter{}

	for _, testCase := range testCases {
		provisioner.provisionCluster = testCase.provision
		accGetter.get = testCase.getAccount
		kubeGetter.get = testCase.kubeGetter

		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(testCase.body))
		rec := httptest.NewRecorder()

		handler := Handler{
			kubeGetter:    kubeGetter,
			provisioner:   provisioner,
			accountGetter: accGetter,
		}

		handler.Provision(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong status code expected %d actual %d", testCase.expectedCode, rec.Code)
			return
		}

		if testCase.expectedCode == http.StatusAccepted {
			resp := ProvisionResponse{}

			err := json.NewDecoder(rec.Body).Decode(&resp)

			if err != nil {
				t.Errorf("Unepxpected error while decoding response %v", err)
			}

			if len(resp.ClusterID) == 0 {
				t.Errorf("ClusterID must not be empty")
			}
		}
	}
}

func TestNewHandler(t *testing.T) {
	accSvc := &account.Service{}
	kubeSvc := &mockKubeService{}
	p := &TaskProvisioner{}
	h := NewHandler(kubeSvc, accSvc, p)

	if h.accountGetter == nil {
		t.Errorf("account getter must not be nil")
	}

	if h.provisioner != p {
		t.Errorf("expected provisioner %v actual %v", p, h.provisioner)
	}
}

func TestHandler_Register(t *testing.T) {
	h := Handler{}
	r := mux.NewRouter()
	h.Register(r)

	expectedRouteCount := 1
	actualRouteCount := 0
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if router != r {
			return errors.New("wrong router")
		}
		actualRouteCount += 1
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error from walk router %v", err)
	}

	if expectedRouteCount != actualRouteCount {
		t.Errorf("Wrong route count expected %d actual %d", expectedRouteCount, actualRouteCount)
	}
}
