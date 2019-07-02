package kube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/proxy"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/testutils"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	errFake = errors.New("fake error")

	deployedReleaseInput = `{"chartName":"nginx","namespace":"default","repoName":"fake"}`
	deployedRelease      = &release.Release{
		Name:      "fakeDeployed",
		Namespace: "default",
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name: "nginx",
			},
		},
		Info: &release.Info{
			Status: &release.Status{
				Code: release.Status_DEPLOYED,
			},
		},
	}
	deployedReleaseInfo = &model.ReleaseInfo{
		Name:      "fakeDeleted",
		Namespace: "default",
		Chart:     "nginx",
		Status:    release.Status_Code_name[int32(release.Status_DEPLOYED)],
	}

	deletedReleaseInput = `{"chartName":"esync","namespace":"kube-system","repoName":"fake"}`
	deletedReleaseInfo  = &model.ReleaseInfo{
		Name:      "fakeDeleted",
		Namespace: "kube-system",
		Chart:     "esync",
		Status:    release.Status_Code_name[int32(release.Status_DELETED)],
	}
)

type kubeServiceMock struct {
	mock.Mock
	rls         *release.Release
	rlsInfo     *model.ReleaseInfo
	rlsInfoList []*model.ReleaseInfo
	rlsErr      error
}

type accServiceMock struct {
	mock.Mock
}

type mockNodeProvisioner struct {
	mock.Mock
}

type mockProvisioner struct {
	mock.Mock
}

type mockProfileService struct {
	mock.Mock
}

func (m *mockProfileService) Get(ctx context.Context,
	profileID string) (*profile.Profile, error) {
	args := m.Called(ctx, profileID)

	val, ok := args.Get(0).(*profile.Profile)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockProfileService) Create(ctx context.Context,
	profile *profile.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *mockProvisioner) RestartClusterProvisioning(ctx context.Context,
	clusterProfile *profile.Profile,
	config *steps.Config,
	taskIdMap map[string][]string) error {
	args := m.Called(ctx, clusterProfile, config, taskIdMap)

	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func (m *mockProvisioner) UpgradeCluster(ctx context.Context, nextVersion string, k *model.Kube,
	tasks map[string][]*workflows.Task,  config *steps.Config) {
	m.Called(ctx, nextVersion, tasks, config)
}

type bufferCloser struct {
	bytes.Buffer
	err error
}

func (b *bufferCloser) Close() error {
	return b.err
}

const (
	serviceCreate            = "Create"
	serviceGet               = "Get"
	serviceListAll           = "ListAll"
	serviceDelete            = "Delete"
	serviceListKubeResources = "ListKubeResources"
	serviceListNodes         = "ListNodes"
	serviceKubeConfigFor     = "KubeConfigFor"
	serviceGetKubeResources  = "GetKubeResources"
	serviceGetCerts          = "GetCerts"
)

func (m *mockNodeProvisioner) ProvisionNodes(ctx context.Context, nodeProfile []profile.NodeProfile, kube *model.Kube, config *steps.Config) ([]string, error) {
	args := m.Called(ctx, nodeProfile, kube, config)
	val, ok := args.Get(0).([]string)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func (m *mockNodeProvisioner) Cancel(clusterID string) error {
	args := m.Called(clusterID)
	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func (m *kubeServiceMock) Create(ctx context.Context, k *model.Kube) error {
	args := m.Called(ctx, k)
	val, ok := args.Get(0).(error)
	if !ok {
		return nil
	}
	return val
}
func (m *kubeServiceMock) Get(ctx context.Context, name string) (*model.Kube, error) {
	args := m.Called(ctx, name)
	val, ok := args.Get(0).(*model.Kube)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) KubeConfigFor(ctx context.Context, kname, user string) ([]byte, error) {
	args := m.Called(ctx, kname, user)
	val, ok := args.Get(0).([]byte)
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

func (m *kubeServiceMock) ListNodes(ctx context.Context, k *model.Kube, role string) ([]corev1.Node, error) {
	args := m.Called(ctx, k, role)
	val, ok := args.Get(0).([]corev1.Node)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
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
	args := m.Called(ctx, kname, cname)
	val, ok := args.Get(0).(*Bundle)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}
func (m *kubeServiceMock) InstallRelease(ctx context.Context,
	kname string, rls *ReleaseInput) (*release.Release, error) {
	return m.rls, m.rlsErr
}
func (m *kubeServiceMock) ReleaseDetails(ctx context.Context,
	kname string, rlsName string) (*release.Release, error) {
	return m.rls, m.rlsErr
}
func (m *kubeServiceMock) ListReleases(ctx context.Context,
	kname, ns, offset string, limit int) ([]*model.ReleaseInfo, error) {
	return m.rlsInfoList, m.rlsErr
}
func (m *kubeServiceMock) DeleteRelease(ctx context.Context,
	kname, rlsName string, purge bool) (*model.ReleaseInfo, error) {
	return m.rlsInfo, m.rlsErr
}

type mockContainter struct {
	mock.Mock
}

func (m *mockContainter) RegisterProxies(targets []proxy.Target) error {
	args := m.Called(targets)
	val, ok := args.Get(0).(error)
	if !ok {
		return args.Error(0)
	}
	return val
}

func (m *mockContainter) GetProxies(prefix string) map[string]*proxy.ServiceReverseProxy {
	args := m.Called(prefix)
	val, ok := args.Get(0).(map[string]*proxy.ServiceReverseProxy)
	if !ok {
		return nil
	}
	return val
}

func (m *mockContainter) Shutdown(ctx context.Context) {
	m.Called(ctx)
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

		serviceCreateError error
		serviceGetResp     *model.Kube
		serviceGetError    error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			rawKube:         []byte(`{"name":"invalid_json"",,}`),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.InvalidJSON,
		},
		{
			rawKube: []byte(`{"name":"newKube"}`),
			serviceGetResp: &model.Kube{
				Name: "alreadyExists",
			},
			expectedStatus:  http.StatusConflict,
			expectedErrCode: sgerrors.AlreadyExists,
		},
		{ // TC#2
			rawKube:         []byte(`{"name":""}`),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{ // TC#3
			rawKube:            []byte(`{"name":"fail_to_put"}`),
			serviceCreateError: errors.New("error"),
			expectedStatus:     http.StatusInternalServerError,
			expectedErrCode:    sgerrors.UnknownError,
		},
		{ // TC#4
			rawKube:        []byte(`{"name":"success"}`),
			expectedStatus: http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil,
			nil, nil, nil, nil, nil)

		req, err := http.NewRequest(http.MethodPost, "/kubes",
			bytes.NewReader(tc.rawKube))
		require.Equalf(t, nil, err,
			"TC#%d: create request: %v", i+1, err)

		svc.On(serviceCreate, mock.Anything, mock.Anything).
			Return(tc.serviceCreateError)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(tc.serviceGetResp, tc.serviceGetError)

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
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

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
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

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
		description string
		kubeName    string

		accountName     string
		getAccountError error
		account         *model.CloudAccount

		kube            *model.Kube
		getKubeError    error
		deleteKubeError error

		expectedStatus int
	}{
		{
			description:    "kube not found",
			kubeName:       "test",
			getKubeError:   sgerrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			description: "account not found",
			kubeName:    "service_error",

			accountName:     "test",
			getAccountError: sgerrors.ErrNotFound,
			account:         nil,
			kube: &model.Kube{
				Provider:    clouds.DigitalOcean,
				Name:        "test",
				AccountName: "test",
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			description:     "delete kube err not found",
			kubeName:        "kubeName",
			getAccountError: nil,
			accountName:     "test",
			account: &model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			getKubeError: nil,
			kube: &model.Kube{
				Provider:    clouds.DigitalOcean,
				Name:        "test",
				AccountName: "test",
				Tasks: map[string][]string{},
			},
			deleteKubeError: sgerrors.ErrNotFound,
			expectedStatus:  http.StatusAccepted,
		},
		{
			description:     "success",
			kubeName:        "delete kube error",
			getAccountError: nil,
			accountName:     "test",
			account: &model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			getKubeError: nil,
			kube: &model.Kube{
				Provider:    clouds.DigitalOcean,
				Name:        "test",
				AccountName: "test",
				Tasks: map[string][]string{},
			},
			deleteKubeError: nil,
			expectedStatus:  http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		t.Log(tc.description)
		// setup handler
		svc := new(kubeServiceMock)
		accSvc := new(accServiceMock)

		// prepare
		req, err := http.NewRequest(http.MethodDelete, "/kubes/"+tc.kubeName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceGet, mock.Anything, tc.kubeName).Return(tc.kube, tc.getKubeError)
		svc.On(serviceDelete, mock.Anything, tc.kubeName).Return(tc.deleteKubeError)
		svc.On(serviceCreate, mock.Anything, mock.Anything).Return(nil)

		accSvc.On(serviceGet, mock.Anything, tc.accountName).Return(tc.account, tc.getAccountError)
		mockRepo := new(testutils.MockStorage)
		mockRepo.On("Put", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("GetAll", mock.Anything,
			mock.Anything).Return([][]byte{}, nil)

		workflows.Init()
		workflows.RegisterWorkFlow(workflows.DeleteCluster, []steps.Step{})

		rr := httptest.NewRecorder()

		mockProvisioner := new(mockNodeProvisioner)
		mockProvisioner.On("Cancel", mock.Anything).
			Return(nil)

		h := NewHandler(svc, accSvc, nil,
			mockProvisioner, nil, mockRepo, nil)

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		if tc.expectedStatus != rr.Code {
			t.Errorf("Wrong response code expected %d actual %d",
				tc.expectedStatus, rr.Code)
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
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

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
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

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

func TestHandler_listNodes(t *testing.T) {
	tcs := []struct {
		name            string
		kubeID          string
		svcNodes        []corev1.Node
		svcGetErr       error
		svcListNodesErr error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{
			name:           "invalid kube",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "kube not found",
			kubeID:          "13",
			svcGetErr:       sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{
			name:            "get kube: internal error",
			kubeID:          "13",
			svcGetErr:       sgerrors.ErrNilEntity,
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			name:            "list nodes error",
			kubeID:          "13",
			svcListNodesErr: sgerrors.ErrNilValue,
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			name:           "list nodes",
			kubeID:         "13",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tcs {

		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kubes/%s/nodes", tc.kubeID), nil)
		require.Equalf(t, nil, err, "TC %s: create request: %v", tc.name, err)

		svc.On(serviceGet, mock.Anything, mock.Anything).Return(&model.Kube{}, tc.svcGetErr)
		svc.On(serviceListNodes, mock.Anything, mock.Anything, mock.Anything).Return(tc.svcNodes, tc.svcListNodesErr)
		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC %s: status code", tc.name)

		if tc.expectedErrCode != sgerrors.ErrorCode(0) {
			m := new(message.Message)
			err = json.NewDecoder(rr.Body).Decode(m)
			require.Equalf(t, nil, err, "TC %s: error codemess", tc.name)

			require.Equalf(t, tc.expectedErrCode, m.ErrorCode, "TC %s", tc.name)
		}
	}
}

func TestAddNodeToKube(t *testing.T) {
	testCases := []struct {
		testName       string
		kubeName       string
		kube           *model.Kube
		kubeServiceErr error

		kubeProfile *profile.Profile
		profileErr  error

		accountName string
		account     *model.CloudAccount
		accountErr  error

		provisionErr error

		expectedCode int
	}{
		{
			testName:       "kube not found",
			kubeName:       "test",
			kube:           nil,
			kubeServiceErr: sgerrors.ErrNotFound,
			kubeProfile:    nil,
			profileErr:     nil,
			accountName:    "",
			accountErr:     nil,
			account:        nil,
			provisionErr:   nil,
			expectedCode:   http.StatusNotFound,
		},
		{
			testName: "profile not found",
			kubeName: "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},
			kubeServiceErr: nil,
			kubeProfile:    nil,
			profileErr:     sgerrors.ErrNotFound,
			accountName:    "test",
			account:        nil,
			accountErr:     sgerrors.ErrNotFound,
			provisionErr:   nil,
			expectedCode:   http.StatusNotFound,
		},
		{
			testName: "account not found",
			kubeName: "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},
			kubeServiceErr: nil,
			kubeProfile:    &profile.Profile{},
			profileErr:     nil,
			accountName:    "test",
			account:        nil,
			accountErr:     sgerrors.ErrNotFound,
			provisionErr:   nil,
			expectedCode:   http.StatusNotFound,
		},
		{
			testName: "provision not found",
			kubeName: "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},
			kubeServiceErr: nil,
			kubeProfile:    &profile.Profile{},
			profileErr:     nil,
			accountName:    "test",
			account: &model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			accountErr:   nil,
			provisionErr: sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			testName: "provision success",
			kubeName: "test",
			kube: &model.Kube{
				AccountName: "test",
				Masters: map[string]*model.Machine{
					"": {},
				},
				Tasks: make(map[string][]string),
			},
			kubeServiceErr: nil,
			kubeProfile:    &profile.Profile{},
			profileErr:     nil,
			accountName:    "test",
			account: &model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			accountErr:   nil,
			provisionErr: nil,
			expectedCode: http.StatusAccepted,
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
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kube, testCase.kubeServiceErr)
		svc.On(serviceCreate, mock.Anything, mock.Anything).
			Return(nil)

		profileSvc := new(mockProfileService)
		profileSvc.On("Get", mock.Anything,
			mock.Anything).Return(testCase.kubeProfile,
			testCase.profileErr)

		accService := new(accServiceMock)
		accService.On("Get", mock.Anything, mock.Anything).
			Return(testCase.account, testCase.accountErr)

		mockProvisioner := new(mockNodeProvisioner)
		mockProvisioner.On("ProvisionNodes",
			mock.Anything, nodeProfile, testCase.kube, mock.Anything).
			Return(mock.Anything, testCase.provisionErr)
		mockProvisioner.On("Cancel", mock.Anything).
			Return(nil)
		h := NewHandler(svc, accService, profileSvc,
			mockProvisioner, nil,
			nil, nil)

		data, _ := json.Marshal(nodeProfile)
		b := bytes.NewBuffer(data)

		req, _ := http.NewRequest(http.MethodPost,
			fmt.Sprintf("/kubes/%s/nodes", testCase.kubeName),
			b)
		rec := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/kubes/{kubeID}/nodes", h.addMachine)
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong error code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestDeleteNodeFromKube(t *testing.T) {
	testCases := []struct {
		testName string

		nodeName       string
		kubeName       string
		kube           *model.Kube
		kubeServiceErr error

		accountName string
		account     *model.CloudAccount
		accountErr  error

		getWriter    func(string) (io.WriteCloser, error)
		expectedCode int
	}{
		{
			"kube not found",
			"test",
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
			"get kube unknown error",
			"test",
			"test",
			nil,
			errors.New("unknown"),
			"",
			nil,
			nil,
			nil,
			http.StatusInternalServerError,
		},
		{
			"method not allowed",
			"test",
			"test",
			&model.Kube{
				Masters: map[string]*model.Machine{
					"test": {
						Name: "test",
					},
				},
			},
			nil,
			"",
			nil,
			nil,
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			"node not found",
			"test",
			"test",
			&model.Kube{
				Nodes: map[string]*model.Machine{
					"test2": {
						Name: "test2",
					},
				},
			},
			nil,
			"",
			nil,
			nil,
			nil,
			http.StatusNotFound,
		},
		{
			"account not found",
			"test",
			"test",
			&model.Kube{
				AccountName: "test",
				Nodes: map[string]*model.Machine{
					"test": {
						Name: "test",
					},
				},
			},
			nil,
			"test",
			nil,
			sgerrors.ErrNotFound,
			nil,
			http.StatusNotFound,
		},
		{
			"account unknown error",
			"test",
			"test",
			&model.Kube{
				AccountName: "test",
				Nodes: map[string]*model.Machine{
					"test": {
						Name: "test",
					},
				},
			},
			nil,
			"test",
			nil,
			errors.New("account unknown error"),
			nil,
			http.StatusInternalServerError,
		},
		{
			"success",
			"test",
			"test",
			&model.Kube{
				AccountName: "test",
				Nodes: map[string]*model.Machine{
					"test": {
						Name: "test",
					},
				},
			},
			nil,
			"test",
			&model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"publicKey": "publicKey",
				},
			},
			nil,
			func(string) (io.WriteCloser, error) {
				return &bufferCloser{}, nil
			},
			http.StatusAccepted,
		},
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.DeleteNode, []steps.Step{})

	for _, testCase := range testCases {
		t.Log(testCase.testName)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kube, testCase.kubeServiceErr)
		svc.On(serviceCreate, mock.Anything, testCase.kube).
			Return(mock.Anything)

		accService := new(accServiceMock)
		accService.On("Get", mock.Anything, mock.Anything).
			Return(testCase.account, testCase.accountErr)

		mockRepo := new(testutils.MockStorage)
		mockRepo.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		mockRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		handler := Handler{
			svc:            svc,
			accountService: accService,
			getWriter:      testCase.getWriter,
			repo:           mockRepo,
		}

		router := mux.NewRouter()
		router.HandleFunc("/{kubeID}/nodes/{nodename}", handler.deleteMachine).Methods(http.MethodDelete)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/nodes/%s", testCase.kubeName, testCase.nodeName), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d", testCase.expectedCode, rec.Code)
		}
	}
}

func TestKubeTasks(t *testing.T) {
	testCases := []struct {
		description string
		repoData    []byte
		repoErr     error

		kubeResp *model.Kube
		kubeErr  error

		err string
	}{
		{
			description: "kube not found",
			kubeResp:    nil,
			kubeErr:     sgerrors.ErrNotFound,
		},
		{
			description: "task not found",
			kubeResp: &model.Kube{
				Tasks: map[string][]string{
					workflows.MasterTask: {"taskID"},
				},
			},
			kubeErr: nil,
			repoErr: sgerrors.ErrNotFound,
		},
		{
			description: "marshall error",
			kubeResp: &model.Kube{
				Tasks: map[string][]string{
					workflows.MasterTask: {"taskID"},
				},
			},
			kubeErr:  nil,
			repoData: []byte(`{`),
			repoErr:  nil,
			err:      "unexpected",
		},
		{
			description: "success",
			kubeResp: &model.Kube{
				Tasks: map[string][]string{
					workflows.MasterTask: {"taskID"},
				},
			},
			kubeErr:  nil,
			repoData: []byte(`{"config": {"clusterId":"test"}}`),
			repoErr:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kubeResp, testCase.kubeErr)

		repo := &testutils.MockStorage{}
		repo.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.repoData, testCase.repoErr)
		h := Handler{
			repo: repo,
			svc:  svc,
		}

		_, err := h.getKubeTasks(context.Background(), "test")

		if err != nil && !strings.Contains(err.Error(), testCase.err) {
			t.Errorf("Wrong error error message expected to have %s actual %s",
				testCase.err, err.Error())
		}
	}
}

func TestDeleteKubeTasks(t *testing.T) {
	testCases := []struct {
		description string
		repoData    []byte
		repoErr     error

		kubeResp *model.Kube
		kubeErr  error

		deleteErr error
	}{
		{
			description: "kube not found",
			kubeErr:     sgerrors.ErrNotFound,
			deleteErr:   sgerrors.ErrNotFound,
		},
		{
			description: "repo not found",
			kubeErr:     nil,
			kubeResp: &model.Kube{
				Tasks: map[string][]string{
					workflows.MasterTask: {"not_found_id"},
				},
			},
			repoErr: sgerrors.ErrNotFound,
		},
		{
			description: "success",
			kubeErr:     nil,
			kubeResp: &model.Kube{
				Tasks: map[string][]string{
					workflows.MasterTask: {"1234"},
				},
			},

			repoData: []byte(`{"config": {"clusterId":"test"}}`),
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kubeResp, testCase.kubeErr)

		repo := &testutils.MockStorage{}
		repo.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.repoData, testCase.repoErr)
		repo.On("Delete", mock.Anything,
			mock.Anything, mock.Anything).
			Return(testCase.deleteErr)
		h := Handler{
			repo: repo,
			svc:  svc,
		}

		err := h.deleteClusterTasks(context.Background(), "test")

		if errors.Cause(err) != testCase.deleteErr {
			t.Errorf("Wrong error expected %v actual %v",
				testCase.deleteErr, err)
		}
	}
}

func TestServiceGetCerts(t *testing.T) {
	testCases := []struct {
		kname string
		cname string

		serviceResp  *Bundle
		serviceErr   error
		expectedCode int
	}{
		{
			kname:        "test",
			cname:        "test",
			serviceResp:  nil,
			serviceErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			kname:        "test",
			cname:        "test",
			serviceResp:  nil,
			serviceErr:   errors.New("unknown"),
			expectedCode: http.StatusInternalServerError,
		},

		{
			kname: "test",
			cname: "test",
			serviceResp: &Bundle{
				Cert: []byte(`cert`),
				Key:  []byte(`key`),
			},
			serviceErr:   nil,
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		svc := new(kubeServiceMock)
		svc.On(serviceGetCerts, mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.serviceResp, testCase.serviceErr)

		h := Handler{
			svc: svc,
		}

		req, _ := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/kubes/%s/certs/%s", testCase.kname, testCase.cname),
			nil)
		rec := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/kubes/{kubeID}/certs/{cname}", h.getCerts)
		router.ServeHTTP(rec, req)

		if testCase.expectedCode != rec.Code {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestGetTasks(t *testing.T) {
	testCases := []struct {
		description string
		kubeID      string
		kubeResp    *model.Kube
		kubeErr     error
		repoData    []byte
		repoErr     error

		expectedCode int
	}{
		{
			description:  "kube not found",
			kubeID:       "test",
			kubeErr:      sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "internal error",
			kubeID:      "test",
			kubeResp: &model.Kube{
				ID: "test",
				Tasks: map[string][]string{
					workflows.MasterTask: {"1234"},
				},
			},
			repoData:     []byte(``),
			expectedCode: http.StatusInternalServerError,
		},
		{
			description: "nothing found",
			kubeID:      "test",
			kubeResp: &model.Kube{
				ID: "test",
				Tasks: map[string][]string{
					workflows.MasterTask: {"1234"},
				},
			},
			repoErr:      sgerrors.ErrInvalidJson,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "success",
			kubeID:      "test",
			kubeResp: &model.Kube{
				ID: "test",
				Tasks: map[string][]string{
					workflows.MasterTask: {"1234"},
				},
			},
			repoData:     []byte(`{"config": {"clusterId":"test"}}`),
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kubeResp, testCase.kubeErr)

		repo := &testutils.MockStorage{}
		repo.On("Get", mock.Anything,
			mock.Anything, mock.Anything).
			Return(testCase.repoData, testCase.repoErr)
		h := Handler{
			repo: repo,
			svc:  svc,
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/kubes/%s/tasks", testCase.kubeID),
			nil)

		router := mux.NewRouter()
		router.HandleFunc("/kubes/{kubeID}/tasks", h.getTasks)
		router.ServeHTTP(rec, req)

		if testCase.expectedCode != rec.Code {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestHandler_installRelease(t *testing.T) {
	tcs := []struct {
		rlsInp string

		kubeSvc *kubeServiceMock

		expectedRls     *release.Release
		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{
			rlsInp: "{{}",
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.InvalidJSON,
		},
		{
			rlsInp: "{}",
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{
			rlsInp: deployedReleaseInput,
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			rlsInp: deployedReleaseInput,
			kubeSvc: &kubeServiceMock{
				rls: deployedRelease,
			},
			expectedStatus: http.StatusOK,
			expectedRls:    deployedRelease,
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.kubeSvc}

		router := mux.NewRouter()
		h.Register(router)

		// prepare
		req, err := http.NewRequest(
			http.MethodPost,
			"/kubes/fake/releases",
			strings.NewReader(tc.rlsInp))
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			rlsInfo := &release.Release{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(rlsInfo), "TC#%d: decode chart", i+1)

			require.Equalf(t, tc.expectedRls, rlsInfo, "TC#%d: check release", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_getRelease(t *testing.T) {
	tcs := []struct {
		kubeSvc *kubeServiceMock

		expectedRls     *release.Release
		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			kubeSvc: &kubeServiceMock{
				rls: deployedRelease,
			},
			expectedStatus: http.StatusOK,
			expectedRls:    deployedRelease,
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.kubeSvc}

		router := mux.NewRouter()
		h.Register(router)

		// prepare
		req, err := http.NewRequest(
			http.MethodGet,
			"/kubes/fake/releases/releaseName",
			nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			rlsInfo := &release.Release{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(rlsInfo), "TC#%d: decode chart", i+1)

			require.Equalf(t, tc.expectedRls, rlsInfo, "TC#%d: check release", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_listReleases(t *testing.T) {
	tcs := []struct {
		description string
		kubeSvc *kubeServiceMock

		k *model.Kube
		expectedRlsInfoList []*model.ReleaseInfo
		expectedStatus      int
		expectedErrCode     sgerrors.ErrorCode
	}{
		{
			description: "kube service error",
			k: &model.Kube{
				State: model.StateOperational,
			},
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			description: "status ok",
			k: &model.Kube{
				State: model.StateOperational,
			},
			kubeSvc: &kubeServiceMock{
				rlsInfoList: []*model.ReleaseInfo{deployedReleaseInfo},
			},
			expectedStatus:      http.StatusOK,
			expectedRlsInfoList: []*model.ReleaseInfo{deployedReleaseInfo},
		},
	}

	for i, tc := range tcs {
		tc.kubeSvc.On("Get", mock.Anything, mock.Anything).Return(tc.k, nil)

		t.Log(tc.description)
		// setup handler
		h := &Handler{svc: tc.kubeSvc}

		router := mux.NewRouter()
		h.Register(router)

		// prepare
		req, err := http.NewRequest(
			http.MethodGet,
			"/kubes/fake/releases",
			nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			rlsInfoList := []*model.ReleaseInfo{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(&rlsInfoList), "TC#%d: decode release list", i+1)

			require.Equalf(t, tc.expectedRlsInfoList, rlsInfoList, "TC#%d: check release", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_getKubeconfig(t *testing.T) {
	tcs := []struct {
		kubeID   string
		userName string

		serviceResources []byte
		serviceError     error

		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			kubeID:         "",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeID:         "cluster1",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeID:          "service_error",
			userName:        "uname",
			serviceError:    errors.New("get error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			kubeID:          "not_found",
			userName:        "uname",
			serviceError:    sgerrors.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#4
			kubeID:         "kubeconfig",
			userName:       "uname",
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		svc := new(kubeServiceMock)
		h := NewHandler(svc, nil, nil,
			nil, nil, nil, nil)

		// prepare
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kubes/%s/users/%s/kubeconfig", tc.kubeID, tc.userName), nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		svc.On(serviceKubeConfigFor, mock.Anything, tc.kubeID, tc.userName).Return(tc.serviceResources, tc.serviceError)
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

func TestHandler_deleteRelease(t *testing.T) {
	tcs := []struct {
		kubeSvc *kubeServiceMock

		expectedRlsInfo *model.ReleaseInfo
		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{
			kubeSvc: &kubeServiceMock{
				rlsErr: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			kubeSvc: &kubeServiceMock{
				rlsInfo: deletedReleaseInfo,
			},
			expectedStatus:  http.StatusOK,
			expectedRlsInfo: deletedReleaseInfo,
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.kubeSvc}

		router := mux.NewRouter()
		h.Register(router)

		// prepare
		req, err := http.NewRequest(
			http.MethodDelete,
			"/kubes/fake/releases/releaseName",
			nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			rlsInfoList := &model.ReleaseInfo{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(rlsInfoList), "TC#%d: decode release info", i+1)

			require.Equalf(t, tc.expectedRlsInfo, rlsInfoList, "TC#%d: check release", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestGetClusterMetrics(t *testing.T) {
	testCases := []struct {
		kubeServiceGetResp  *model.Kube
		kubeServiceGetError error
		getMetrics          func(string, *model.Kube) (*MetricResponse, error)
		expectedCode        int
	}{
		{
			kubeServiceGetError: sgerrors.ErrNotFound,
			expectedCode:        http.StatusNotFound,
		},
		{
			kubeServiceGetError: errors.New("unknown error"),
			expectedCode:        http.StatusInternalServerError,
		},
		{
			kubeServiceGetResp: &model.Kube{
				Name: "test",
				Masters: map[string]*model.Machine{
					"master-1": {
						Name:     "master-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			kubeServiceGetError: nil,
			getMetrics: func(string, *model.Kube) (*MetricResponse, error) {
				return nil, sgerrors.ErrInvalidJson
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			kubeServiceGetResp: &model.Kube{
				Name: "test",
				Masters: map[string]*model.Machine{
					"master-1": {
						Name:     "master-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			kubeServiceGetError: nil,
			getMetrics: func(string, *model.Kube) (*MetricResponse, error) {
				return &MetricResponse{
					Data: struct {
						ResultType string `json:"resultType"`
						Result     []struct {
							Metric map[string]string `json:"metric"`
							Value  []interface{}     `json:"value"`
						} `json:"result"`
					}{
						ResultType: "metric",
						Result: []struct {
							Metric map[string]string `json:"metric"`
							Value  []interface{}     `json:"value"`
						}{
							{

								Value: []interface{}{"cpu", 0.42},
							},
							{

								Value: []interface{}{"memory", 0.65},
							},
						},
					},
				}, nil
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		svc := new(kubeServiceMock)
		svc.On("Get", mock.Anything, mock.Anything).
			Return(testCase.kubeServiceGetResp, testCase.kubeServiceGetError)

		handler := Handler{
			svc:        svc,
			getMetrics: testCase.getMetrics,
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/kubes/%s/metrics", "test"), nil)

		router := mux.NewRouter().SkipClean(true)
		handler.Register(router)

		// run
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestGetNodesMetrics(t *testing.T) {
	expectedNodeCount := 3
	testCases := []struct {
		kubeServiceGetResp  *model.Kube
		kubeServiceGetError error
		getMetrics          func(string, *model.Kube) (*MetricResponse, error)
		expectedCode        int
	}{
		{
			kubeServiceGetError: sgerrors.ErrNotFound,
			expectedCode:        http.StatusNotFound,
		},
		{
			kubeServiceGetError: errors.New("unknown error"),
			expectedCode:        http.StatusInternalServerError,
		},
		{
			kubeServiceGetResp: &model.Kube{
				Name: "test",
				Masters: map[string]*model.Machine{
					"master-1": {
						Name:     "master-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			kubeServiceGetError: nil,
			getMetrics: func(string, *model.Kube) (*MetricResponse, error) {
				return nil, sgerrors.ErrInvalidJson
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			kubeServiceGetResp: &model.Kube{
				Name: "test",
				Masters: map[string]*model.Machine{
					"master-1": {
						Name:     "master-1",
						PublicIp: "10.20.30.40",
					},
				},
			},
			kubeServiceGetError: nil,
			getMetrics: func(string, *model.Kube) (*MetricResponse, error) {
				return &MetricResponse{
					Data: struct {
						ResultType string `json:"resultType"`
						Result     []struct {
							Metric map[string]string `json:"metric"`
							Value  []interface{}     `json:"value"`
						} `json:"result"`
					}{
						ResultType: "metric",
						Result: []struct {
							Metric map[string]string `json:"metric"`
							Value  []interface{}     `json:"value"`
						}{
							{
								Metric: map[string]string{
									"node": "node-1",
								},
								Value: []interface{}{"memory", 0.42},
							},
							{

								Metric: map[string]string{
									"node": "node-2",
								},
								Value: []interface{}{"memory", 0.54},
							},
							{

								Metric: map[string]string{
									"node": "master-1",
								},
								Value: []interface{}{"memory", 0.77},
							},
							{
								Metric: map[string]string{
									"node": "node-1",
								},
								Value: []interface{}{"cpu", 0.21},
							},
							{

								Metric: map[string]string{
									"node": "node-2",
								},
								Value: []interface{}{"cpu", 0.35},
							},
							{

								Metric: map[string]string{
									"node": "master-1",
								},
								Value: []interface{}{"cpu", 0.69},
							},
						},
					},
				}, nil
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		svc := new(kubeServiceMock)
		svc.On("Get", mock.Anything, mock.Anything).
			Return(testCase.kubeServiceGetResp, testCase.kubeServiceGetError)

		handler := Handler{
			svc:        svc,
			getMetrics: testCase.getMetrics,
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/kubes/%s/nodes/metrics", "test"), nil)

		router := mux.NewRouter().SkipClean(true)
		handler.Register(router)

		// run
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}

		if testCase.expectedCode == http.StatusOK {
			resp := map[string]interface{}{}

			err := json.NewDecoder(rec.Body).Decode(&resp)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if len(resp) != expectedNodeCount {
				t.Errorf("Unexpected count of nodes expected %d actual %d",
					expectedNodeCount, len(resp))
			}
		}
	}
}

func TestRestarProvisioningKube(t *testing.T) {
	testCases := []struct {
		description string
		kubeName    string

		kube           *model.Kube
		kubeServiceErr error

		kubeProfile *profile.Profile
		profileErr  error

		accountName string
		account     *model.CloudAccount
		accountErr  error

		provisionErr error

		expectedCode int
	}{
		{
			description:    "kube not found",
			kubeName:       "test",
			kubeServiceErr: sgerrors.ErrNotFound,
			expectedCode:   http.StatusNotFound,
		},
		{
			description: "profile not found",
			kubeName:    "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},

			profileErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "account not found",
			kubeName:    "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},

			kubeProfile:  &profile.Profile{},
			accountName:  "not found",
			accountErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "unsupported cloud provider",
			kubeName:    "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},

			kubeProfile: &profile.Profile{},
			accountName: "not found",
			account: &model.CloudAccount{
				Provider: "unsupported",
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			description: "Error while provision",
			kubeName:    "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},

			kubeProfile: &profile.Profile{},
			accountName: "not found",
			account: &model.CloudAccount{
				Provider: "unsupported",
			},
			provisionErr: errors.New("provision error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			description: "Success",
			kubeName:    "test",
			kube: &model.Kube{
				AccountName: "test",
				Tasks:       make(map[string][]string),
			},

			kubeProfile: &profile.Profile{},
			accountName: "not found",
			account: &model.CloudAccount{
				Provider: clouds.AWS,
			},
			expectedCode: http.StatusAccepted,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, mock.Anything).
			Return(testCase.kube, testCase.kubeServiceErr)
		svc.On(serviceCreate, mock.Anything, mock.Anything).
			Return(nil)

		profileSvc := new(mockProfileService)
		profileSvc.On("Get", mock.Anything,
			mock.Anything).Return(testCase.kubeProfile,
			testCase.profileErr)

		accService := new(accServiceMock)
		accService.On("Get", mock.Anything, mock.Anything).
			Return(testCase.account, testCase.accountErr)

		mockProvisioner := new(mockProvisioner)
		mockProvisioner.On("RestartClusterProvisioning",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.provisionErr)

		h := NewHandler(svc, accService, profileSvc,
			nil, mockProvisioner,
			nil, nil)

		req, _ := http.NewRequest(http.MethodPost,
			fmt.Sprintf("/kubes/%s/restart", testCase.kubeName),
			nil)
		rec := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/kubes/{kubeID}/restart", h.restartKubeProvisioning)
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong error code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestGetServices(t *testing.T) {
	testCases := []struct {
		name string

		getKubeErr error
		getKube    *model.Kube

		getServicesErr error
		k8sServices    *corev1.ServiceList

		registerProxiesErr error
		getProxies         map[string]*proxy.ServiceReverseProxy

		expectedCode int
	}{
		{
			name:         "kube not found",
			getKubeErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "kube internal error",
			getKubeErr:   errors.New("unknown"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "get services error",
			getKube: &model.Kube{
				ID: "1234",
				Masters: map[string]*model.Machine{
					"key": {
						ID: "key",
					},
				},
			},
			getServicesErr: errors.New("error"),
			expectedCode:   http.StatusInternalServerError,
		},
		{
			name: "register proxy error",

			getKube: &model.Kube{
				ID: "1234",
				Masters: map[string]*model.Machine{
					"key": {
						ID: "key",
					},
				},
			},
			k8sServices:        &corev1.ServiceList{},
			registerProxiesErr: errors.New("error"),
			expectedCode:       http.StatusInternalServerError,
		},
		{
			name: "success 1",

			getKube: &model.Kube{
				ID: "1234",
				Masters: map[string]*model.Machine{
					"key": {
						ID: "key",
					},
				},
			},
			k8sServices: &corev1.ServiceList{
				Items: []corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								clusterService: "false",
							},
						},
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{
								{
									Name:     "http",
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			getProxies: map[string]*proxy.ServiceReverseProxy{
				"kubeID": {
					ServingBase: "http:/10.20.30.40:9090",
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success 2",

			getKube: &model.Kube{
				ID: "1234",
				Masters: map[string]*model.Machine{
					"key": {
						ID: "key",
					},
				},
			},
			k8sServices: &corev1.ServiceList{
				Items: []corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								clusterService: "true",
							},
						},
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{
								{
									Name:     "http",
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			getProxies: map[string]*proxy.ServiceReverseProxy{
				"kubeID": {
					ServingBase: "http:/10.20.30.40:9090",
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success 3",

			getKube: &model.Kube{
				ID: "1234",
				Masters: map[string]*model.Machine{
					"key": {
						ID: "key",
					},
				},
			},
			k8sServices: &corev1.ServiceList{
				Items: []corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								clusterService: "true",
							},
						},
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{
								{
									Name:     "other",
									Protocol: "unknown",
								},
								{
									Name:     "http",
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			getProxies: map[string]*proxy.ServiceReverseProxy{
				"kubeID": {
					ServingBase: "http:/10.20.30.40:9090",
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		kubeSvc := &kubeServiceMock{}
		kubeSvc.On("Get", mock.Anything, mock.Anything).
			Return(tc.getKube, tc.getKubeErr)
		mockProxies := &mockContainter{}
		mockProxies.On("RegisterProxies",
			mock.Anything).Return(tc.registerProxiesErr)
		mockProxies.On("GetProxies",
			mock.Anything).Return(tc.getProxies)
		getSvc := func(*model.Kube, string) (*corev1.ServiceList, error) {
			return tc.k8sServices, tc.getServicesErr
		}

		handler := &Handler{
			listK8sServices: getSvc,
			svc:             kubeSvc,
			proxies:         mockProxies,
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet,
			"/kubes/kubeID/services", nil)

		router := mux.NewRouter().SkipClean(true)
		handler.Register(router)

		router.ServeHTTP(rec, req)

		if rec.Code != tc.expectedCode {
			t.Errorf("TC: %s: wrong response code expected "+
				"%d actual %d", tc.name, tc.expectedCode, rec.Code)
		}
	}
}

func TestImportKube(t *testing.T) {
	testCases := []struct {
		description string

		req []byte

		accountName string
		account     *model.CloudAccount
		accountErr  error

		profileErr error

		svcNodes  []corev1.Node
		svcGetErr error

		expectedCode int
	}{
		{
			description:  "json error",
			req:          []byte(`{`),
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "bad credentials",
			req:          []byte(`{"kubeconfig":"{}","clusterName":"kubernetes","cloudAccountName":"test"}`),
			accountErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "cloud account not found",
			req:          []byte(`{"kubeconfig":"{\r\n  \"kind\": \"Config\",\r\n  \"apiVersion\": \"v1\",\r\n  \"preferences\": {},\r\n  \"clusters\": [\r\n    {\r\n      \"name\": \"asdfasdf\",\r\n      \"cluster\": {\r\n        \"server\": \"https:\/\/ex-24adfede-130460518.eu-west-2.elb.amazonaws.com\",\r\n        \"certificate-authority-data\": \"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1EUXdPVEUyTlRNeE1sb1hEVEk1TURRd05qRTJOVE14TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSyt5Clc5VTdjZWJQV0FBWXNVR2pZUVM5aC9tblAyWW4zVUtxTStaT3QzQ2Z6MVk5ekhaaTlyK0pObEgrWkwralI2QWYKamYyTzRScEFvSG5uYTUzMGEwM2s3dFp3bTdiNXZCcEZLTmw2aHhoKzU2Y1RzMUxZbVJuTWRERFlRV2JSbXk2bwo3ZFRsaDVBVHY1K21tUlNMMkxja1lraDRqTWhObWFPb1hLVmxzck5SWXZ1NHAvRk5uNHF3OE0xekxXK25uSG5kCmZrSWlJZHRXb2ROMG8yL0Njb2l3QW5uSXpGVmVIYnF5L3ZqTm1aOFc2NU5PbW4yZHk1cnkwd0EvQVFzRDdXUS8Kb3I3c2NVRkEvdnRJNXJ4eWVNM2xhMjFycjVGMnhpbVcrZWNlUVNkY1JvK2RoaXFlYmsrcDJrRnV2SjBBZDlUdgpla0dtS0dhRXRPZnE1R2lkSnBrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFDdmlBNXlBKzhHT1NZWjlyam5tN2h0b3BoQkYKemFGMXNadHBudzZuZzM3cjNSSFJwUVlnZUpMcXZZaGtvQlZ5M0lsc1JTV0c2NGJwYUFubjdEb3JZSjdzYmFmeQpIdFRQYlF4S1Bxa0NyMGlwUkxBZmdtdDlodVNLbTlNQUVwWTlNL1NXdmpvNXVoZUg0RWJFWXViSGJhV0Z4eFpPCktTeXlZeTc5WGpIKy9pQndFemoxcWxYUzVsQ1dIUjN6SUUycnM5cVNKWnA2MW9NWDlmYWFYUElSZHJvOHNpVWIKMG1kOUZFeVcrc05GL05xREtUNTFzbHVYR2lWZ1lrK0diUnJ1L1IrbkpZNnlQUU1uK212UTFWN29Ic2RoUTJ1dApRUllZcytkRCtSMW1tNXdNOEIzL3NPSDRnelpKVmtNOFdteUg5a1RDMGFzbkszNDRGOTdrQktaN2VUST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\"\r\n      }\r\n    }\r\n  ],\r\n  \"users\": [\r\n    {\r\n      \"name\": \"admin@asdfasdf\",\r\n      \"user\": {\r\n        \"client-certificate-data\": \"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJSlFTLzQrYjgyLzR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RBME1Ea3hOalV6TVRKYUZ3MHlNREEwTURneE5qVXpNVEphTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXVDRGlvUHBOcjlnb2kyY3IKUUd6d21tVFU3OUV4WDN0VUZGUUw1clZoQkVTTjdma2k0MWNueCtBRkhCbVRnODRNVStlR0VqditudGYvWEdQSApTR3FiYlZkOUFmM2hMV2dBUnBFdCtVVTZFUUJSTDdtUE9qZFI0WFhRQVk3RlNHam9wUlgvcWdUdFFJZ05MS0tHClVESzhMQVV2bkVoaVQrN0hKUGdlZGVJNG9SeHh4NUpvdXpqUlk0ODkyOGtNTE02Mm1ZMmV1bkFqMi8vWmtna2QKNmhKT3dxN0t5ck9jY0k3NVA3RE0xM3BtUDduZDA2SHp6VFJ5ZGxwbEJQbmErcHAwaDN0Q2xpNG5GZG5yakFwYgpiZWxKYUtDUElseEF1Mk00a1BBWDRZdGUvd0hiRVROQjNHbVhTWjNxQ3hocFhVdnlCYzdBalVWVnZncmJvbU11CkdvR2dBUUlEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFGVEJWM0crNEdQaGtWUW9USU5XZDlTQjE2UTRCZjNGUm5BbApuODRJZzBWNFRkNndCZG1lbVZhYzFjeU94dkFSQmpGMEVoRXFaODFjVDZuK3NoMzZrYmh0Rzl3RDd4WU1lanhRClRBZnZDL01ndFo0YVl5Qnp2Uk5yWmxQYkoyUUlpdXo3RmM0NWFSUnh5LzJEVkVXYTdXaytzbUUrR0dHTnR0OFQKRUQzWjBhSTFWSkxDcDhqR0xVeVg3V3FRNU5YckN0TE95cnd0UHZMNGlLTnNZd2VwYzRYUTBacXBEM0VDMERJdApKZ0Vzb0FybDVYdTVad0oxbWtwS2x4RGhEOVZHTGExRkV1YmtNVTh4cGpQd0JzU2xWb3V3bzVVbVhETEhKUE5kCmlkQjBSdjZvUENzTlZvclpVUDRrR1lMdXA3NTJnc0FRSHQyUm5HcUZWUkJSRTBQWWc3az0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\",\r\n        \"client-key-data\": \"LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdUNEaW9QcE5yOWdvaTJjclFHendtbVRVNzlFeFgzdFVGRlFMNXJWaEJFU043ZmtpCjQxY254K0FGSEJtVGc4NE1VK2VHRWp2K250Zi9YR1BIU0dxYmJWZDlBZjNoTFdnQVJwRXQrVVU2RVFCUkw3bVAKT2pkUjRYWFFBWTdGU0dqb3BSWC9xZ1R0UUlnTkxLS0dVREs4TEFVdm5FaGlUKzdISlBnZWRlSTRvUnh4eDVKbwp1empSWTQ4OTI4a01MTTYybVkyZXVuQWoyLy9aa2drZDZoSk93cTdLeXJPY2NJNzVQN0RNMTNwbVA3bmQwNkh6CnpUUnlkbHBsQlBuYStwcDBoM3RDbGk0bkZkbnJqQXBiYmVsSmFLQ1BJbHhBdTJNNGtQQVg0WXRlL3dIYkVUTkIKM0dtWFNaM3FDeGhwWFV2eUJjN0FqVVZWdmdyYm9tTXVHb0dnQVFJREFRQUJBb0lCQUVCcmVQN2JNK3p1MHBpYgpPdDJxZjY5MDBhOHA0SDFJZDgwNDdvVUVObkk0emZOUmMreFlKTm5nUGNsc0JWbGE4S00yWUlqbXBwbktkbFJMCnNJQmNsQTU0U29zMDJPQjYvNFd3NjRYSHU1NFJIMVowTkhFb3c2UC9CUXhXZUIxeVh0ckxUSXllZHFkYU1rbkkKZnZkMkdMSEtDck5CKyt0OVhmMUlOZGdHa3N1Q3dPaFBuM3Z0cXJ2cktRSzJCQnAwVURJWXdHbVY2dmxCM3JmLwoxZDhmd0VQMEc0YVdzUDR4UnFSemt5bU5YbDNIcmZTUGZQTzZrdjhWMElWWjhFeFFXMjVXcWRsMERVcXhBZjBICmtGSHNtMGFYY1FVQk9ScmhxalByL1g5ZnBxbkNnMlJRSC9SVEV6ZXpoc2NrQWxHUjVFTTcxak42RkxVZFllVEgKRkdjSWdSVUNnWUVBNTk5OW1MVmU2WTd5RHJpL0pnQXRyVWxJWUN3SDNhOENmeVkrOURTWVdSdWFRdktReGVqNQpiSUhBQnBXVGp4YW95MjlYR01meCtxci9wdThzdWlzRjFueWtEUFh5cVRUeHgvMmIwRjg1eXA5M3pRcW9DOWZSCm8yTThpNDkxSWhWODZ6ajk5b0IvWkcvRjZBeGhBMlVOOG9jdnBaTHF0M1M0M2NZeDF2c2F0aE1DZ1lFQXkwbVoKek0wY0J5UTdhZlI3RVpFcUVTNGoxVkxGbTBBVldNcmRJaThuM0JhZEJ3RTFscXJ0NDJndU5VSG9MbUVEZENCcQpnSmYxNXZIYWFBa1Jjem95c1Y5SG5aZVhoQUxNMG11dU01amxpdFNiUEhvczl3WG5RQVJzTFdvNVJlSlpkWVliClhkSWVCLzNNTUtKdUFaejBvcVZaOTZ2SWpwajRQS1pWWThoTFpCc0NnWUVBdEI1UzlUWW14VzFxTU85b1pQK00KZStqS1ZSSy9CWUMyZ3NqVjdHT1MyTjF0Um9ZZzJld3hIUTNwZWZQbFRTaS85RS9JSzVMZU1PZDJjbG1tdC9ORgp0S2piMHNtWE44UE44Wm5hMk5Hd0Zlc3NaOVhZVm1MUEVZbTc5WGw1OXdFVUtiRDY3dXBBaTJlY0o3YStBYUlWClpJbUpCS2lNdGZmd3h5MzNkMVZXR1lzQ2dZQTlQK1hKSVJ1S3cwM3JkTEFIOFBiOXlpc2R3UnlzMURnYVVyVWgKOFpkTzVybFZQUFlLZVdISG5NSWZaY1l4QXlYcFBVTVpqNitWYjlWZ2R5cjh6dElyUXd2dTNaZlhQSWk5OVplOQpFQnBKSkJuSnRQNExSNG9QYmNXeVFVa1VWMGlnOGxFWWlaQm0wLzlMd0FUcEU0Tlo1ZndmZFhDdUZrVGs4VERWCkthb2RkUUtCZ0hkWnc4MU1sV3ZBU3EyUlVRK3BjK0wxL3lvb2F1bm93b0laak5uTVJETFFkYi9nSGhYYmNKUG0KOVRGeXBtYkthMVlGMkFjN2tJbkhxOFNCYUQyanhLaHZLNkVONFNpV2t3MExrb2JiZHR4OXlPVmYxcEc5MVZvUgpVZEhvS20wUnREbmRYMjhaNUxRR1FKL01DZkphMGJHNURKaFZ6SXgyQXBlSTlhWmVNbnN5Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==\"\r\n      }\r\n    }\r\n  ],\r\n  \"contexts\": [\r\n    {\r\n      \"name\": \"admin@asdfasdf\",\r\n      \"context\": {\r\n        \"cluster\": \"asdfasdf\",\r\n        \"user\": \"admin@asdfasdf\"\r\n      }\r\n    }\r\n  ],\r\n  \"current-context\": \"admin@asdfasdf\"\r\n}","clusterName":"kubernetes","cloudAccountName":"test"}`),
			accountErr:   sgerrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "success",
			req:         []byte(`{"kubeconfig":"{\r\n  \"kind\": \"Config\",\r\n  \"apiVersion\": \"v1\",\r\n  \"preferences\": {},\r\n  \"clusters\": [\r\n    {\r\n      \"name\": \"asdfasdf\",\r\n      \"cluster\": {\r\n        \"server\": \"https:\/\/ex-24adfede-130460518.eu-west-2.elb.amazonaws.com\",\r\n        \"certificate-authority-data\": \"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1EUXdPVEUyTlRNeE1sb1hEVEk1TURRd05qRTJOVE14TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSyt5Clc5VTdjZWJQV0FBWXNVR2pZUVM5aC9tblAyWW4zVUtxTStaT3QzQ2Z6MVk5ekhaaTlyK0pObEgrWkwralI2QWYKamYyTzRScEFvSG5uYTUzMGEwM2s3dFp3bTdiNXZCcEZLTmw2aHhoKzU2Y1RzMUxZbVJuTWRERFlRV2JSbXk2bwo3ZFRsaDVBVHY1K21tUlNMMkxja1lraDRqTWhObWFPb1hLVmxzck5SWXZ1NHAvRk5uNHF3OE0xekxXK25uSG5kCmZrSWlJZHRXb2ROMG8yL0Njb2l3QW5uSXpGVmVIYnF5L3ZqTm1aOFc2NU5PbW4yZHk1cnkwd0EvQVFzRDdXUS8Kb3I3c2NVRkEvdnRJNXJ4eWVNM2xhMjFycjVGMnhpbVcrZWNlUVNkY1JvK2RoaXFlYmsrcDJrRnV2SjBBZDlUdgpla0dtS0dhRXRPZnE1R2lkSnBrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFDdmlBNXlBKzhHT1NZWjlyam5tN2h0b3BoQkYKemFGMXNadHBudzZuZzM3cjNSSFJwUVlnZUpMcXZZaGtvQlZ5M0lsc1JTV0c2NGJwYUFubjdEb3JZSjdzYmFmeQpIdFRQYlF4S1Bxa0NyMGlwUkxBZmdtdDlodVNLbTlNQUVwWTlNL1NXdmpvNXVoZUg0RWJFWXViSGJhV0Z4eFpPCktTeXlZeTc5WGpIKy9pQndFemoxcWxYUzVsQ1dIUjN6SUUycnM5cVNKWnA2MW9NWDlmYWFYUElSZHJvOHNpVWIKMG1kOUZFeVcrc05GL05xREtUNTFzbHVYR2lWZ1lrK0diUnJ1L1IrbkpZNnlQUU1uK212UTFWN29Ic2RoUTJ1dApRUllZcytkRCtSMW1tNXdNOEIzL3NPSDRnelpKVmtNOFdteUg5a1RDMGFzbkszNDRGOTdrQktaN2VUST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\"\r\n      }\r\n    }\r\n  ],\r\n  \"users\": [\r\n    {\r\n      \"name\": \"admin@asdfasdf\",\r\n      \"user\": {\r\n        \"client-certificate-data\": \"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJSlFTLzQrYjgyLzR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB4T1RBME1Ea3hOalV6TVRKYUZ3MHlNREEwTURneE5qVXpNVEphTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXVDRGlvUHBOcjlnb2kyY3IKUUd6d21tVFU3OUV4WDN0VUZGUUw1clZoQkVTTjdma2k0MWNueCtBRkhCbVRnODRNVStlR0VqditudGYvWEdQSApTR3FiYlZkOUFmM2hMV2dBUnBFdCtVVTZFUUJSTDdtUE9qZFI0WFhRQVk3RlNHam9wUlgvcWdUdFFJZ05MS0tHClVESzhMQVV2bkVoaVQrN0hKUGdlZGVJNG9SeHh4NUpvdXpqUlk0ODkyOGtNTE02Mm1ZMmV1bkFqMi8vWmtna2QKNmhKT3dxN0t5ck9jY0k3NVA3RE0xM3BtUDduZDA2SHp6VFJ5ZGxwbEJQbmErcHAwaDN0Q2xpNG5GZG5yakFwYgpiZWxKYUtDUElseEF1Mk00a1BBWDRZdGUvd0hiRVROQjNHbVhTWjNxQ3hocFhVdnlCYzdBalVWVnZncmJvbU11CkdvR2dBUUlEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFGVEJWM0crNEdQaGtWUW9USU5XZDlTQjE2UTRCZjNGUm5BbApuODRJZzBWNFRkNndCZG1lbVZhYzFjeU94dkFSQmpGMEVoRXFaODFjVDZuK3NoMzZrYmh0Rzl3RDd4WU1lanhRClRBZnZDL01ndFo0YVl5Qnp2Uk5yWmxQYkoyUUlpdXo3RmM0NWFSUnh5LzJEVkVXYTdXaytzbUUrR0dHTnR0OFQKRUQzWjBhSTFWSkxDcDhqR0xVeVg3V3FRNU5YckN0TE95cnd0UHZMNGlLTnNZd2VwYzRYUTBacXBEM0VDMERJdApKZ0Vzb0FybDVYdTVad0oxbWtwS2x4RGhEOVZHTGExRkV1YmtNVTh4cGpQd0JzU2xWb3V3bzVVbVhETEhKUE5kCmlkQjBSdjZvUENzTlZvclpVUDRrR1lMdXA3NTJnc0FRSHQyUm5HcUZWUkJSRTBQWWc3az0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\",\r\n        \"client-key-data\": \"LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdUNEaW9QcE5yOWdvaTJjclFHendtbVRVNzlFeFgzdFVGRlFMNXJWaEJFU043ZmtpCjQxY254K0FGSEJtVGc4NE1VK2VHRWp2K250Zi9YR1BIU0dxYmJWZDlBZjNoTFdnQVJwRXQrVVU2RVFCUkw3bVAKT2pkUjRYWFFBWTdGU0dqb3BSWC9xZ1R0UUlnTkxLS0dVREs4TEFVdm5FaGlUKzdISlBnZWRlSTRvUnh4eDVKbwp1empSWTQ4OTI4a01MTTYybVkyZXVuQWoyLy9aa2drZDZoSk93cTdLeXJPY2NJNzVQN0RNMTNwbVA3bmQwNkh6CnpUUnlkbHBsQlBuYStwcDBoM3RDbGk0bkZkbnJqQXBiYmVsSmFLQ1BJbHhBdTJNNGtQQVg0WXRlL3dIYkVUTkIKM0dtWFNaM3FDeGhwWFV2eUJjN0FqVVZWdmdyYm9tTXVHb0dnQVFJREFRQUJBb0lCQUVCcmVQN2JNK3p1MHBpYgpPdDJxZjY5MDBhOHA0SDFJZDgwNDdvVUVObkk0emZOUmMreFlKTm5nUGNsc0JWbGE4S00yWUlqbXBwbktkbFJMCnNJQmNsQTU0U29zMDJPQjYvNFd3NjRYSHU1NFJIMVowTkhFb3c2UC9CUXhXZUIxeVh0ckxUSXllZHFkYU1rbkkKZnZkMkdMSEtDck5CKyt0OVhmMUlOZGdHa3N1Q3dPaFBuM3Z0cXJ2cktRSzJCQnAwVURJWXdHbVY2dmxCM3JmLwoxZDhmd0VQMEc0YVdzUDR4UnFSemt5bU5YbDNIcmZTUGZQTzZrdjhWMElWWjhFeFFXMjVXcWRsMERVcXhBZjBICmtGSHNtMGFYY1FVQk9ScmhxalByL1g5ZnBxbkNnMlJRSC9SVEV6ZXpoc2NrQWxHUjVFTTcxak42RkxVZFllVEgKRkdjSWdSVUNnWUVBNTk5OW1MVmU2WTd5RHJpL0pnQXRyVWxJWUN3SDNhOENmeVkrOURTWVdSdWFRdktReGVqNQpiSUhBQnBXVGp4YW95MjlYR01meCtxci9wdThzdWlzRjFueWtEUFh5cVRUeHgvMmIwRjg1eXA5M3pRcW9DOWZSCm8yTThpNDkxSWhWODZ6ajk5b0IvWkcvRjZBeGhBMlVOOG9jdnBaTHF0M1M0M2NZeDF2c2F0aE1DZ1lFQXkwbVoKek0wY0J5UTdhZlI3RVpFcUVTNGoxVkxGbTBBVldNcmRJaThuM0JhZEJ3RTFscXJ0NDJndU5VSG9MbUVEZENCcQpnSmYxNXZIYWFBa1Jjem95c1Y5SG5aZVhoQUxNMG11dU01amxpdFNiUEhvczl3WG5RQVJzTFdvNVJlSlpkWVliClhkSWVCLzNNTUtKdUFaejBvcVZaOTZ2SWpwajRQS1pWWThoTFpCc0NnWUVBdEI1UzlUWW14VzFxTU85b1pQK00KZStqS1ZSSy9CWUMyZ3NqVjdHT1MyTjF0Um9ZZzJld3hIUTNwZWZQbFRTaS85RS9JSzVMZU1PZDJjbG1tdC9ORgp0S2piMHNtWE44UE44Wm5hMk5Hd0Zlc3NaOVhZVm1MUEVZbTc5WGw1OXdFVUtiRDY3dXBBaTJlY0o3YStBYUlWClpJbUpCS2lNdGZmd3h5MzNkMVZXR1lzQ2dZQTlQK1hKSVJ1S3cwM3JkTEFIOFBiOXlpc2R3UnlzMURnYVVyVWgKOFpkTzVybFZQUFlLZVdISG5NSWZaY1l4QXlYcFBVTVpqNitWYjlWZ2R5cjh6dElyUXd2dTNaZlhQSWk5OVplOQpFQnBKSkJuSnRQNExSNG9QYmNXeVFVa1VWMGlnOGxFWWlaQm0wLzlMd0FUcEU0Tlo1ZndmZFhDdUZrVGs4VERWCkthb2RkUUtCZ0hkWnc4MU1sV3ZBU3EyUlVRK3BjK0wxL3lvb2F1bm93b0laak5uTVJETFFkYi9nSGhYYmNKUG0KOVRGeXBtYkthMVlGMkFjN2tJbkhxOFNCYUQyanhLaHZLNkVONFNpV2t3MExrb2JiZHR4OXlPVmYxcEc5MVZvUgpVZEhvS20wUnREbmRYMjhaNUxRR1FKL01DZkphMGJHNURKaFZ6SXgyQXBlSTlhWmVNbnN5Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==\"\r\n      }\r\n    }\r\n  ],\r\n  \"contexts\": [\r\n    {\r\n      \"name\": \"admin@asdfasdf\",\r\n      \"context\": {\r\n        \"cluster\": \"asdfasdf\",\r\n        \"user\": \"admin@asdfasdf\"\r\n      }\r\n    }\r\n  ],\r\n  \"current-context\": \"admin@asdfasdf\"\r\n}","clusterName":"kubernetes","cloudAccountName":"test"}`),
			account: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			expectedCode: http.StatusAccepted,
		},
	}

	workflows.Init()
	workflows.RegisterWorkFlow(workflows.ImportCluster, []steps.Step{})

	for _, testCase := range testCases {
		t.Log(testCase.description)
		svc := &kubeServiceMock{}
		svc.On(serviceListNodes, mock.Anything, mock.Anything, mock.Anything).Return(testCase.svcNodes, testCase.svcGetErr)
		svc.On(serviceCreate, mock.Anything, mock.Anything).Return(nil)
		accSvc := &accServiceMock{}
		accSvc.On("Get", mock.Anything, mock.Anything).
			Return(testCase.account, testCase.accountErr)

		profileSvc := &mockProfileService{}
		profileSvc.On("Create", mock.Anything, mock.Anything).Return(testCase.profileErr)

		mockRepo := new(testutils.MockStorage)
		mockRepo.On("Put", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)

		h := NewHandler(svc, accSvc,
			profileSvc, nil,
			nil, mockRepo, nil)

		rr := httptest.NewRecorder()

		router := mux.NewRouter().SkipClean(true)
		h.Register(router)

		req, _ := http.NewRequest(http.MethodPost,
			"/kubes/import",
			bytes.NewBuffer(testCase.req))

		h.importKube(rr, req)

		if rr.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rr.Code)
		}
	}
}
