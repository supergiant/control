package kube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type kubeServiceMock struct {
	mock.Mock
}

type accServiceMock struct {
	mock.Mock
}

type mockNodeProvisioner struct {
	mock.Mock
}

const (
	serviceCreate            = "Create"
	serviceGet               = "Get"
	serviceListAll           = "ListAll"
	serviceDelete            = "Delete"
	serviceListKubeResources = "ListKubeResources"
	serviceGetKubeResources  = "GetKubeResources"
)

func (m *mockNodeProvisioner) ProvisionNodes(ctx context.Context, nodeProfile []profile.NodeProfile, kube *model.Kube, config *steps.Config) ([]string, error) {
	args := m.Called(ctx, nodeProfile, kube, config)
	val, ok := args.Get(0).([]string)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *kubeServiceMock) Create(ctx context.Context, k *model.Kube) error {
	args := m.Called(ctx, k)
	return args.Error(0)
}
func (m *kubeServiceMock) Get(ctx context.Context, name string) (*model.Kube, error) {
	args := m.Called(ctx, name)
	val, ok := args.Get(0).(*model.Kube)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) ListAll(ctx context.Context) ([]model.Kube, error) {
	args := m.Called(ctx)
	val, ok := args.Get(0).([]model.Kube)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}
func (m *kubeServiceMock) ListKubeResources(ctx context.Context, kname string) ([]byte, error) {
	args := m.Called(ctx, kname)
	val, ok := args.Get(0).([]byte)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) GetKubeResources(ctx context.Context, kname, resource, ns, name string) ([]byte, error) {
	args := m.Called(ctx, kname, resource, ns, name)
	val, ok := args.Get(0).([]byte)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) GetCerts(ctx context.Context, kname, cname string) (*Bundle, error) {
	return nil, nil
}

func (a *accServiceMock) Get(ctx context.Context, name string) (*model.CloudAccount, error) {
	args := a.Called(ctx, name)

	val, ok := args.Get(0).(*model.CloudAccount)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestHandler_createKube(t *testing.T) {
	tcs := []struct {
		rawKube []byte

		serviceError error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			rawKube:         []byte(`{"name":"invalid_json"",,}`),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.InvalidJSON,
		},
		{ // TC#2
			rawKube:         []byte(`{"name":""}`),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{ // TC#3
			rawKube:         []byte(`{"name":"fail_to_put"}`),
			serviceError:    errors.New("error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#4
			rawKube:        []byte(`{"name":"success"}`),
			expectedStatus: http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		req, err := http.NewRequest(http.MethodPost, "/kubes", bytes.NewReader(tc.rawKube))
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceCreate, mock.Anything, mock.Anything).Return(tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}
	}
}

func TestHandler_getKube(t *testing.T) {
	tcs := []struct {
		kubeName string

		serviceKube  *model.Kube
		serviceError error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			kubeName:       "",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeName:        "service_error",
			serviceError:    errors.New("get error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			kubeName:        "not_found",
			serviceError:    sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#4
			kubeName: "success",
			serviceKube: &model.Kube{
				Name: "success",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, "/kubes/"+tc.kubeName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceGet, mock.Anything, tc.kubeName).Return(tc.serviceKube, tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}

		if tc.serviceKube != nil {
			k := new(model.Kube)
			err = json.NewDecoder(rr.Body).Decode(k)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, k, tc.serviceKube, "TC#%d", i+1)
		}
	}
}

func TestHandler_listKubes(t *testing.T) {
	tcs := []struct {
		serviceKubes []model.Kube
		serviceError error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			serviceError:    errors.New("error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#2
			expectedStatus: http.StatusOK,
			serviceKubes: []model.Kube{
				{
					Name: "success",
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, "/kubes", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceListAll, mock.Anything).Return(tc.serviceKubes, tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}

		if tc.serviceKubes != nil {
			kubes := new([]model.Kube)
			err = json.NewDecoder(rr.Body).Decode(kubes)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.serviceKubes, *kubes, "TC#%d", i+1)
		}
	}
}

func TestHandler_deleteKube(t *testing.T) {
	tcs := []struct {
		kubeName string

		serviceError error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			kubeName:       "",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeName:        "service_error",
			serviceError:    errors.New("get error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			kubeName:        "not_found",
			serviceError:    sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#3
			kubeName:       "delete",
			expectedStatus: http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodDelete, "/kubes/"+tc.kubeName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceDelete, mock.Anything, tc.kubeName).Return(tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}
	}
}

func TestHandler_listResources(t *testing.T) {
	tcs := []struct {
		kubeName string

		serviceResources []byte
		serviceError     error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			kubeName:       "",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeName:        "service_error",
			serviceError:    errors.New("get error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			kubeName:        "not_found",
			serviceError:    sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#4
			kubeName:       "list_resources",
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kubes/%s/resources", tc.kubeName), nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceListKubeResources, mock.Anything, tc.kubeName).Return(tc.serviceResources, tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}
	}
}

func TestHandler_getResources(t *testing.T) {
	tcs := []struct {
		kubeName     string
		resourceName string

		serviceResources []byte
		serviceError     error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			kubeName:       "",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeName:        "service_error",
			resourceName:    "service_error",
			serviceError:    errors.New("get error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			kubeName:        "not_found",
			resourceName:    "not_found",
			serviceError:    sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#4
			kubeName:       "list_resources",
			resourceName:   "list_resources",
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kubes/%s/resources/%s", tc.kubeName, tc.resourceName), nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceGetKubeResources, mock.Anything, tc.kubeName, mock.Anything, mock.Anything, mock.Anything).
			Return(tc.serviceResources, tc.serviceError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC#%d", i+1)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC#%d", i+1)
		}
	}
}

func TestAddNodeToKube(t *testing.T) {
	testCases := []struct {
		testName       string
		kubeName       string
		kube           *model.Kube
		kubeServiceErr error

		accountName string
		account     *model.CloudAccount
		accountErr  error

		provisionErr error

		expectedCode int
	}{
		{
			"kube not found",
			"test",
			nil,
			sgerrors.ErrNotFound,
			"",
			nil,
			nil,
			nil,
			http.StatusNotFound,
		},
		{
			"account not found",
			"test",
			&model.Kube{
				AccountName: "test",
			},
			nil,
			"test",
			nil,
			sgerrors.ErrNotFound,
			nil,
			http.StatusNotFound,
		},
		{
			"provision not found",
			"test",
			&model.Kube{
				AccountName: "test",
			},
			nil,
			"test",
			&model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			nil,
			sgerrors.ErrNotFound,
			http.StatusNotFound,
		},
		{
			"provision success",
			"test",
			&model.Kube{
				AccountName: "test",
				Masters:     []*node.Node{{}},
			},
			nil,
			"test",
			&model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			nil,
			nil,
			http.StatusAccepted,
		},
	}

	nodeProfile := []profile.NodeProfile{
		{
			"size":  "s-2vcpu-4gb",
			"image": "ubuntu-18-04-x64",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.testName)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, testCase.kubeName).
			Return(testCase.kube, testCase.kubeServiceErr)

		accService := new(accServiceMock)
		accService.On("Get", mock.Anything, testCase.accountName).
			Return(testCase.account, testCase.accountErr)

		mockProvisioner := new(mockNodeProvisioner)
		mockProvisioner.On("ProvisionNodes",
			mock.Anything, nodeProfile, testCase.kube, mock.Anything).
			Return(mock.Anything, testCase.provisionErr)

		h := NewHandler(svc, accService, mockProvisioner, nil)

		data, _ := json.Marshal(nodeProfile)
		b := bytes.NewBuffer(data)

		req, _ := http.NewRequest(http.MethodPost,
			fmt.Sprintf("/kubes/%s/nodes", testCase.kubeName),
			b)
		rec := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/kubes/{kname}/nodes", h.addNode)
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong error code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}
