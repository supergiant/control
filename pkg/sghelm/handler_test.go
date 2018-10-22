package sghelm

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
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

var (
	errFake = errors.New("fake error")
)

var _ Servicer = &fakeService{}

type fakeService struct {
	repo      *model.Repository
	repoList  []model.Repository
	chrt      *model.Chart
	chrtList  []model.Chart
	chrtFiles *chart.Chart
	err       error
}

func (fs fakeService) CreateRepo(ctx context.Context, e *repo.Entry) (*model.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) GetRepo(ctx context.Context, repoName string) (*model.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) ListRepos(ctx context.Context) ([]model.Repository, error) {
	return fs.repoList, fs.err
}
func (fs fakeService) DeleteRepo(ctx context.Context, repoName string) (*model.Repository, error) {
	return fs.repo, fs.err
}
func (fs fakeService) GetChartInfo(ctx context.Context, repoName, chartName string) (*model.Chart, error) {
	return fs.chrt, fs.err
}
func (fs fakeService) ListChartInfos(ctx context.Context, repoName string) ([]model.Chart, error) {
	return fs.chrtList, fs.err
}
func (fs fakeService) GetChart(ctx context.Context, repoName, chartName, chartVersion string) (*chart.Chart, error) {
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
		expectedRepo    *model.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			inpRepo:         []byte("{name:invalidJSON,,}"),
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{ // TC#2
			inpRepo: []byte(`{"name":"validationFailed"}`),
			svc: &fakeService{
				err: sgerrors.ErrAlreadyExists,
			},
			expectedStatus:  http.StatusBadRequest,
			expectedErrCode: sgerrors.ValidationFailed,
		},
		{ // TC#3
			inpRepo: []byte(`{"name":"alreadyExists","url":"url"}`),
			svc: &fakeService{
				err: sgerrors.ErrAlreadyExists,
			},
			expectedStatus:  http.StatusConflict,
			expectedErrCode: sgerrors.AlreadyExists,
		},
		{ // TC#4
			inpRepo: []byte(`{"name":"createError","url":"url"}`),
			svc: &fakeService{
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#5
			inpRepo: []byte(`{"name":"sgRepo","url":"url"}`),
			svc: &fakeService{
				repo: &model.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &model.Repository{
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
			hrepo := &model.Repository{}
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
		expectedRepo    *model.Repository
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
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			svc: &fakeService{
				repo: &model.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &model.Repository{
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
			hrepo := &model.Repository{}
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
		expectedRepos   []model.Repository
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			svc: &fakeService{
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#2
			svc: &fakeService{
				repoList: []model.Repository{
					{
						Config: repo.Entry{
							Name: "sgRepo",
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepos: []model.Repository{
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
			hrepos := []model.Repository{}
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
		expectedRepo    *model.Repository
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
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			svc: &fakeService{
				repo: &model.Repository{
					Config: repo.Entry{
						Name: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedRepo: &model.Repository{
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
			hrepo := &model.Repository{}
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
		expectedChart   *model.Chart
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
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#3
			repoName: "sgRepo",
			chrtName: "sgChart",
			svc: &fakeService{
				chrt: &model.Chart{
					Name: "sgChart",
					Repo: "sgRepo",
				},
			},
			expectedStatus: http.StatusOK,
			expectedChart: &model.Chart{
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
			chrt := &model.Chart{}
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
		expectedCharts  []model.Chart
		expectedErrCode sgerrors.ErrorCode
	}{
		{ // TC#1
			repoName: "listChartError",
			svc: &fakeService{
				err: errFake,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: sgerrors.UnknownError,
		},
		{ // TC#2
			repoName: "sgRepo",
			svc: &fakeService{
				chrtList: []model.Chart{
					{
						Name: "sgChart",
						Repo: "sgRepo",
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedCharts: []model.Chart{
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
			charts := []model.Chart{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(&charts), "TC#%d: decode repos", i+1)

			require.Equalf(t, tc.expectedCharts, charts, "TC#%d: check repos", i+1)
		} else {
			apiErr := &message.Message{}
			require.Nilf(t, json.NewDecoder(w.Body).Decode(apiErr), "TC#%d: decode message", i+1)

			require.Equalf(t, tc.expectedErrCode, apiErr.ErrorCode, "TC#%d: check error code", i+1)
		}
	}
}
