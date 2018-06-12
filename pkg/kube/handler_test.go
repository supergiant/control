package kube

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/testutils"
)


func TestHandler_createKube(t *testing.T) {
	tcs := []struct {
		kube    *Kube
		rawKube []byte

		storageError error

		expectedStatus int
	}{
		{ // TC#1
			rawKube: []byte("{name:invalid_json,,}"),
			kube: &Kube{
				Name: "invalid_json",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{ // TC#2
			kube: &Kube{
				Name: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{ // TC#3
			storageError: errors.New("error"),
			kube: &Kube{
				Name: "fail_to_put",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#4
			kube: &Kube{
				Name: "success",
			},
			expectedStatus: http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(NewService(DefaultStoragePrefix, storage))

		// prepare
		if tc.rawKube == nil {
			raw, err := json.Marshal(tc.kube)
			require.Equalf(t, nil, err, "TC#%d: %v", i+1, err)
			tc.rawKube = raw
		}

		req, err := http.NewRequest(http.MethodPost, "/kubes", bytes.NewReader(tc.rawKube))
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		storage.On(testutils.StoragePut, mock.Anything, mock.Anything, tc.kube.Name, mock.Anything).Return(tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		// TODO: check error message
	}
}

func TestHandler_getKube(t *testing.T) {
	tcs := []struct {
		kubeName string

		storageKube    *Kube
		storageKubeRaw []byte
		storageError error

		expectedStatus int
	}{
		{ // TC#1
			kubeName:       "not_found",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#2
			kubeName: "storage_error",
			storageError: errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			kubeName:       "invalid_json",
			storageKubeRaw: []byte("{name;}"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#3
			kubeName:       "not_found",
			expectedStatus: http.StatusNotFound,
		},
		{ // TC#4
			kubeName: "stable",
			storageKube: &Kube{
				Name: "stable",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(NewService(DefaultStoragePrefix, storage))

		// prepare
		req, err := http.NewRequest(http.MethodGet, "/kubes/"+tc.kubeName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		if tc.storageKube != nil {
			raw, err := json.Marshal(tc.storageKube)
			require.Equalf(t, nil, err, "TC#%d: %v", i+1, err)
			tc.storageKubeRaw = raw
		}
		storage.On(testutils.StorageGet, mock.Anything, mock.Anything, tc.kubeName).Return(tc.storageKubeRaw, tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.storageKube != nil {
			repo := new(Kube)
			require.Nil(t, json.NewDecoder(rr.Body).Decode(repo))

			require.Equalf(t, tc.storageKube, repo, "TC#%d", i+1)
		}
	}
}

func TestHandler_listKubes(t *testing.T) {
	tcs := []struct {
		storageKube  *Kube
		storageError error

		expectedStatus int
	}{
		{ // TC#1
			storageError: errors.New("storage error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			storageKube: &Kube{
				Name: "stable",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(NewService(DefaultStoragePrefix, storage))

		// prepare
		req, err := http.NewRequest(http.MethodGet, "/kubes", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		raw, err := json.Marshal(tc.storageKube)
		require.Nil(t, err)

		storage.On(testutils.StorageGetAll, mock.Anything, mock.Anything).Return([][]byte{raw}, tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.storageKube != nil {
			repos := make([]Kube, 1)
			if err = json.NewDecoder(rr.Body).Decode(&repos); err != nil {
				t.Errorf("TC#%d: decode body: %v", i+1, err)
			}

			require.Equalf(t, []Kube{*tc.storageKube}, repos, "TC#%d", i+1)
		}
	}
}

func TestHandler_deleteKube(t *testing.T) {
	tcs := []struct {
		kubeName string

		storageError error

		expectedStatus int
	}{
		{ // TC#1
			kubeName: "not_found",
			storageError: errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			kubeName: "delete",
			expectedStatus: http.StatusAccepted,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(NewService(DefaultStoragePrefix, storage))

		// prepare
		req, err := http.NewRequest(http.MethodDelete, "/kubes/"+tc.kubeName, nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		storage.On(testutils.StorageDelete, mock.Anything, mock.Anything, mock.Anything).Return(tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)
	}
}

func TestHandler_listResources(t *testing.T) {
	tcs := []struct {
		storageKube  *Kube
		storageError error

		expectedStatus int
	}{
		{ // TC#1
			storageError: errors.New("storage error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{ // TC#2
			storageKube: &Kube{
				Name: "stable",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for i, tc := range tcs {
		// setup handler
		storage := new(testutils.MockStorage)
		h := NewHandler(NewService(DefaultStoragePrefix, storage))

		// prepare
		req, err := http.NewRequest(http.MethodGet, "/kubes", nil)
		require.Equalf(t, nil, err, "TC#%d: create request: %v", i+1, err)

		raw, err := json.Marshal(tc.storageKube)
		require.Nil(t, err)

		storage.On(testutils.StorageGetAll, mock.Anything, mock.Anything).Return([][]byte{raw}, tc.storageError)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		h.Register(router)

		// run
		router.ServeHTTP(rr, req)

		// check
		// TODO: check error message
		require.Equalf(t, tc.expectedStatus, rr.Code, "TC#%d", i+1)

		if tc.storageKube != nil {
			repos := make([]Kube, 1)
			if err = json.NewDecoder(rr.Body).Decode(&repos); err != nil {
				t.Errorf("TC#%d: decode body: %v", i+1, err)
			}

			require.Equalf(t, []Kube{*tc.storageKube}, repos, "TC#%d", i+1)
		}
	}
}
