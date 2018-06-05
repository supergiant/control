package repository

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/mux"

	"github.com/supergiant/supergiant/internal/testutils"
	"github.com/supergiant/supergiant/pkg/model/helm"
)

func TestHandler_Create(t *testing.T) {
	tcs := []struct {
		repo    *helm.Repository
		rawRepo []byte

		storageError error

		expectedStatus int
	}{
		{ // TC#1
			rawRepo: []byte("{name:invalid_json,,}"),
			repo: &helm.Repository{
				Name: "invalid_json",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{ // TC#2
			repo: &helm.Repository{
				Name: "fail_validation",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{ // TC#3
			storageError: errors.New("error"),
			repo: &helm.Repository{
				Name: "fail_to_put",
				URL:  "test",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#4
			repo: &helm.Repository{
				Name: "success",
				URL:  "test",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(storage)

		// prepare
		if tc.rawRepo == nil {
			raw, err := json.Marshal(tc.repo)
			require.Equalf(t, nil, err, "TC#%d: %v", i+1, err)
			tc.rawRepo = raw
		}

		req, err := http.NewRequest("", "", bytes.NewReader(tc.rawRepo))
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		storage.On(testutils.StoragePut, mock.Anything, mock.Anything, tc.repo.Name, mock.Anything).Return(tc.storageError)
		rr := httptest.NewRecorder()

		// run
		http.HandlerFunc(h.Create).ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		// TODO: check error message
	}
}

func TestHandler_Get(t *testing.T) {
	tcs := []struct {
		repoName string

		storageRepo    *helm.Repository
		storageRepoRaw []byte

		expectedStatus int
	}{
		{ // TC#1
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			repoName:       "invalid_json",
			storageRepoRaw: []byte("{name;}"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#3
			repoName:       "not_found",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#4
			repoName: "stable",
			storageRepo: &helm.Repository{
				Name: "stable",
				URL:  "stable",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(storage)

		// prepare
		req, err := http.NewRequest("", "/helm/"+tc.repoName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		if tc.storageRepo != nil {
			raw, err := json.Marshal(tc.storageRepo)
			require.Equalf(t, nil, err, "TC#%d: %v", i+1, err)
			tc.storageRepoRaw = raw
		}
		storage.On(testutils.StorageGet, mock.Anything, mock.Anything, tc.repoName).Return(tc.storageRepoRaw, nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/helm/{repoName}", h.Get)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.storageRepo != nil {
			repo := new(helm.Repository)
			require.Nil(t, json.NewDecoder(rr.Body).Decode(repo))

			require.Equalf(t, tc.storageRepo, repo, "TC#%d", i+1)
		}
	}
}

func TestHandler_ListAll(t *testing.T) {
	tcs := []struct {
		storageRepo  *helm.Repository
		storageError error

		expectedStatus int
	}{
		{ // TC#1
			storageError:   errors.New("storage error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			storageRepo: &helm.Repository{
				Name: "stable",
				URL:  "stable",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(storage)

		// prepare
		req, err := http.NewRequest("", "", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		raw, err := json.Marshal(tc.storageRepo)
		require.Nil(t, err)

		storage.On(testutils.StorageGetAll, mock.Anything, mock.Anything).Return([][]byte{raw}, tc.storageError)
		rr := httptest.NewRecorder()

		// run
		http.HandlerFunc(h.ListAll).ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.storageRepo != nil {
			repos := make([]helm.Repository, 1)
			if err = json.NewDecoder(rr.Body).Decode(&repos); err != nil {
				t.Errorf("TC#%d: decode body: %v", i+1, err)
			}

			require.Equalf(t, []helm.Repository{*tc.storageRepo}, repos, "TC#%d", i+1)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	tcs := []struct {
		repoName string

		storageError error

		expectedStatus int
	}{
		{ // TC#1
			repoName:       "not_found",
			storageError:   errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			repoName:       "delete",
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(storage)

		// prepare
		req, err := http.NewRequest("", "/helm/"+tc.repoName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		storage.On(testutils.StorageDelete, mock.Anything, mock.Anything, mock.Anything).Return(tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/helm/{repoName}", h.Delete)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)
	}
}
