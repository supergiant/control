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
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

var (
	fakeErr = errors.New("fake error")

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
func (m *kubeServiceMock) InstallRelease(ctx context.Context, kname string, rls *ReleaseInput) (*release.Release, error) {
	return m.rls, m.rlsErr
}
func (m *kubeServiceMock) ListReleases(ctx context.Context, kname, ns, offset string, limit int) ([]*model.ReleaseInfo, error) {
	return m.rlsInfoList, m.rlsErr
}
func (m *kubeServiceMock) DeleteRelease(ctx context.Context, kname, rlsName string, purge bool) (*model.ReleaseInfo, error) {
	return m.rlsInfo, m.rlsErr
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
				Name:        "test",
				AccountName: "test",
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			kubeName:        "delete kube error",
			getAccountError: nil,
			accountName:     "test",
			account: &model.CloudAccount{
				Name:     "test",
				Provider: clouds.DigitalOcean,
			},
			getKubeError: nil,
			kube: &model.Kube{
				Name:        "test",
				AccountName: "test",
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
				Name:        "test",
				AccountName: "test",
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

		accSvc.On(serviceGet, mock.Anything, tc.accountName).Return(tc.account, tc.getAccountError)
		mockRepo := new(testutils.MockStorage)
		mockRepo.On("Put", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Delete", mock.Anything,
			mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("GetAll", mock.Anything,
			mock.Anything).Return([][]byte{}, nil)

		workflows.Init()
		workflows.RegisterWorkFlow(workflows.DigitalOceanDeleteCluster, []steps.Step{})

		rr := httptest.NewRecorder()

		h := NewHandler(svc, accSvc, nil, mockRepo)

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
				Masters: map[string]*node.Node{
					"": {},
				},
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
			"method not allowed",
			"test",
			"test",
			&model.Kube{
				Masters: map[string]*node.Node{
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
				Nodes: map[string]*node.Node{
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
				Nodes: map[string]*node.Node{
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
			"success",
			"test",
			"test",
			&model.Kube{
				AccountName: "test",
				Nodes: map[string]*node.Node{
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
	workflows.RegisterWorkFlow(workflows.DigitalOceanDeleteNode, []steps.Step{})

	for _, testCase := range testCases {
		t.Log(testCase.testName)
		svc := new(kubeServiceMock)
		svc.On(serviceGet, mock.Anything, testCase.kubeName).
			Return(testCase.kube, testCase.kubeServiceErr)
		svc.On(serviceCreate, mock.Anything, testCase.kube).Return(mock.Anything)

		accService := new(accServiceMock)
		accService.On("Get", mock.Anything, testCase.accountName).
			Return(testCase.account, testCase.accountErr)

		mockRepo := new(testutils.MockStorage)
		mockRepo.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		mockRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		handler := Handler{
			svc:            svc,
			accountService: accService,
			workflowMap: map[clouds.Name]workflows.WorkflowSet{
				clouds.DigitalOcean: {
					DeleteNode: workflows.DigitalOceanDeleteNode},
			},
			getWriter: testCase.getWriter,
			repo:      mockRepo,
		}

		router := mux.NewRouter()
		router.HandleFunc("/{kname}/nodes/{nodename}", handler.deleteNode).Methods(http.MethodDelete)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/nodes/%s", testCase.kubeName, testCase.nodeName), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d", testCase.expectedCode, rec.Code)
		}
	}
}

func TestHander_installRelease(t *testing.T) {
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
				rlsErr: fakeErr,
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.InvalidJSON,
		},
		{
			rlsInp: "{}",
			kubeSvc: &kubeServiceMock{
				rlsErr: fakeErr,
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{
			rlsInp: deployedReleaseInput,
			kubeSvc: &kubeServiceMock{
				rlsErr: fakeErr,
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

func TestHander_listReleases(t *testing.T) {
	tcs := []struct {
		kubeSvc *kubeServiceMock

		expectedRlsInfoList []*model.ReleaseInfo
		expectedStatus      int
		expectedErrCode     sgerrors.ErrorCode
	}{
		{
			kubeSvc: &kubeServiceMock{
				rlsErr: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{
			kubeSvc: &kubeServiceMock{
				rlsInfoList: []*model.ReleaseInfo{deployedReleaseInfo},
			},
			expectedStatus:      http.StatusOK,
			expectedRlsInfoList: []*model.ReleaseInfo{deployedReleaseInfo},
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

func TestHander_deleteRelease(t *testing.T) {
	tcs := []struct {
		kubeSvc *kubeServiceMock

		expectedRlsInfo *model.ReleaseInfo
		expectedStatus  int
		expectedErrCode sgerrors.ErrorCode
	}{
		{
			kubeSvc: &kubeServiceMock{
				rlsErr: fakeErr,
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
