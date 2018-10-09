package helm

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

var (
	fakeErr = errors.New("fake error")
)

type fakeService struct {
	repo      *helm.Repository
	repoList  []helm.Repository
	chrt      *helm.Chart
	chrtList  []helm.Chart
	chrtFiles *chart.Chart
	err       error
}

func (fs fakeService) CreateRepo(ctx context.Context, e *repo.Entry) (*helm.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) GetRepo(ctx context.Context, repoName string) (*helm.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) ListRepos(ctx context.Context) ([]helm.Repository, error) {
	return fs.repoList, fs.err
}
func (fs fakeService) DeleteRepo(ctx context.Context, repoName string) (*helm.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) GetChart(ctx context.Context, repoName, chartName string) (*helm.Chart, error) {
	return fs.chrt, fs.err
}
func (fs fakeService) ListCharts(ctx context.Context, repoName string) ([]helm.Chart, error) {
	return fs.chrtList, fs.err
}
func (fs fakeService) GetChartFiles(ctx context.Context, repoName, chartName, chartVersion string) (*chart.Chart, error) {
	return fs.chrtFiles, fs.err
}

func TestHandler_createRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc     *fakeService
		inpRepo []byte

		expectedStatus  int
		expectedRepo    *helm.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			inpRepo:         []byte("{name:invalidJSON,,}"),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.InvalidJSON,
		},
		{ // TC#2
			inpRepo: []byte(`{"name":"alreadyExists"}`),
			svc: &fakeService{
				err: sgerrors.ErrAlreadyExists,
			},
			expectedStatus:  http.StatusConflict,
			expectedErrCode: sgerrors.AlreadyExists,
		},
		{ // TC#3
			inpRepo: []byte(`{"name":"createError"}`),
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#4
			inpRepo: []byte(`{"name":"sgRepo"}`),
			svc: &fakeService{
				repo: &helm.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "sgRepo",
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		// prepare
		req, err := http.NewRequest("", "", bytes.NewReader(tc.inpRepo))
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		http.HandlerFunc(h.createRepo).ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			hrepo := &helm.Repository{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(hrepo), "TC#%d: decode repo", i+1)

			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check repo", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_getRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc      *fakeService
		repoName string

		expectedStatus  int
		expectedRepo    *helm.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			repoName: "notFound",
			svc: &fakeService{
				err: sgerrors.ErrNotFound,
			},
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#2
			repoName: "getError",
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			svc: &fakeService{
				repo: &helm.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "sgRepo",
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		// prepare
		req, err := http.NewRequest("", "/helm/"+tc.repoName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		router := mux.NewRouter()
		router.HandleFunc("/helm/{repoName}", h.getRepo)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			hrepo := &helm.Repository{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(hrepo), "TC#%d: decode repo", i+1)

			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check repo", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_listRepos(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc *fakeService

		expectedStatus  int
		expectedRepos   []helm.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#2
			svc: &fakeService{
				repoList: []helm.Repository{
					{
						Config: repo.Entry{
							Name: "sgRepo",
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepos: []helm.Repository{
				{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		// prepare
		req, err := http.NewRequest("", "", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		http.HandlerFunc(h.listRepos).ServeHTTP(w, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d", i+1)

		if w.Code == http.StatusOK {
			hrepos := []helm.Repository{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(&hrepos), "TC#%d: decode repos", i+1)

			require.Equalf(t, tc.expectedRepos, hrepos, "TC#%d: check repos", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_deleteRepo(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc      *fakeService
		repoName string

		expectedStatus  int
		expectedRepo    *helm.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			repoName: "notFound",
			svc: &fakeService{
				err: sgerrors.ErrNotFound,
			},
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#2
			repoName: "deleteError",
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			svc: &fakeService{
				repo: &helm.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &helm.Repository{
				Config: repo.Entry{
					Name: "sgRepo",
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		// prepare
		req, err := http.NewRequest("", "/helm/"+tc.repoName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		router := mux.NewRouter()
		router.HandleFunc("/helm/{repoName}", h.deleteRepo)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			hrepo := &helm.Repository{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(hrepo), "TC#%d: decode repo", i+1)

			require.Equalf(t, tc.expectedRepo, hrepo, "TC#%d: check repo", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_getChart(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc      *fakeService
		repoName string
		chrtName string

		expectedStatus  int
		expectedChart   *helm.Chart
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			repoName: "sgRepo",
			chrtName: "notFound",
			svc: &fakeService{
				err: sgerrors.ErrNotFound,
			},
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: sgerrors.NotFound,
		},
		{ // TC#2
			repoName: "sgRepo",
			chrtName: "getChartError",
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			chrtName: "sgChart",
			svc: &fakeService{
				chrt: &helm.Chart{
					Name: "sgChart",
					Repo: "sgRepo",
				},
			},
			expectedStatus: http.StatusOK,
			expectedChart: &helm.Chart{
				Name: "sgChart",
				Repo: "sgRepo",
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		router := mux.NewRouter()
		h.Register(router)

		// prepare
		req, err := http.NewRequest("", "/helm/repositories/"+tc.repoName+"/charts/"+tc.chrtName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		router.ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d: check status code", i+1)

		if w.Code == http.StatusOK {
			chrt := &helm.Chart{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(chrt), "TC#%d: decode repo", i+1)

			require.Equalf(t, tc.expectedChart, chrt, "TC#%d: check repo", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}

func TestHandler_listCharts(t *testing.T) {
	loggerWriter := logrus.StandardLogger().Out
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(loggerWriter)

	tcs := []struct {
		svc      *fakeService
		repoName string

		expectedStatus  int
		expectedCharts  []helm.Chart
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			repoName: "listChartError",
			svc: &fakeService{
				err: fakeErr,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#2
			repoName: "sgRepo",
			svc: &fakeService{
				chrtList: []helm.Chart{
					{
						Name: "sgChart",
						Repo: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedCharts: []helm.Chart{
				{
					Name: "sgChart",
					Repo: "sgRepo",
				},
			},
		},
	}

	for i, tc := range tcs {
		// setup handler
		h := &Handler{svc: tc.svc}

		// prepare
		req, err := http.NewRequest("", "", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		w := httptest.NewRecorder()

		// run
		http.HandlerFunc(h.listCharts).ServeHTTP(w, req)

		// check
		require.Equalf(t, tc.expectedStatus, w.Code, "TC#%d", i+1)

		if w.Code == http.StatusOK {
			charts := []helm.Chart{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(&charts), "TC#%d: decode repos", i+1)

			require.Equalf(t, tc.expectedCharts, charts, "TC#%d: check repos", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}
