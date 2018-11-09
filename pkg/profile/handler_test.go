package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func TestKubeProfileEndpointCreateProfileErr(t *testing.T) {
	endpoint := &Handler{}
	kubeProfile := &Profile{
		ID:          "",
		K8SVersion:  "1.11.1",
		RBACEnabled: false,
	}

	data, _ := json.Marshal(kubeProfile)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/kubeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d", http.StatusBadRequest, rr.Code)
	}
}

func TestKubeProfileEndpointCreateProfileSuccess(t *testing.T) {
	kubeProfile := &Profile{
		ID: "key",

		MasterProfiles: []NodeProfile{
			{"hello": "world"},
		},
		NodesProfiles: []NodeProfile{
			{"hello1": "world1"},
			{"hello2": "world2"},
		},

		K8SVersion:            "1.11.1",
		Provider:              clouds.AWS,
		Region:                "fra1",
		Arch:                  "amd64",
		OperatingSystem:       "linux",
		UbuntuVersion:         "xenial",
		DockerVersion:         "1.18.1",
		FlannelVersion:        "0.9.0",
		NetworkType:           "vxlan",
		CIDR:                  "10.0.0.1/24",
		HelmVersion:           "0.11.1",
		RBACEnabled:           false,
		CloudSpecificSettings: map[string]string{},
	}

	mockRepo := &testutils.MockStorage{}
	data, _ := json.Marshal(kubeProfile)
	mockRepo.On("Put", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(nil)
	svc := NewService("prefix", mockRepo)
	endpoint := &Handler{
		service: svc,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost,
		"/kubeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Wrong response code, expected %d actual %d",
			http.StatusCreated, rr.Code)
	}
}

func TestKubeProfileEndpointCreateProfileErrJson(t *testing.T) {
	mockRepo := &testutils.MockStorage{}
	data := []byte(`{`)
	mockRepo.On("Put", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(nil)
	svc := NewService("prefix", mockRepo)
	endpoint := &Handler{
		service: svc,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost,
		"/kubeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Wrong response code, expected %d actual %d",
			http.StatusBadRequest, rr.Code)
	}
}

func TestKubeProfileEndpointCreateProfileInternalError(t *testing.T) {
	kubeProfile := &Profile{
		ID: "key",

		MasterProfiles: []NodeProfile{
			{"hello": "world"},
		},
		NodesProfiles: []NodeProfile{
			{"hello1": "world1"},
			{"hello2": "world2"},
		},

		K8SVersion:      "1.11.1",
		Provider:        clouds.AWS,
		Region:          "fra1",
		Arch:            "amd64",
		OperatingSystem: "linux",
		UbuntuVersion:   "xenial",
		DockerVersion:   "1.18.1",
		FlannelVersion:  "0.9.0",
		NetworkType:     "vxlan",
		CIDR:            "10.0.0.1/24",
		HelmVersion:     "0.11.1",
		RBACEnabled:     false,
	}

	mockRepo := &testutils.MockStorage{}
	data, _ := json.Marshal(kubeProfile)
	mockRepo.On("Put", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(sgerrors.ErrAlreadyExists)
	svc := NewService("prefix", mockRepo)
	endpoint := &Handler{
		service: svc,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost,
		"/kubeprofile", bytes.NewReader(data))

	handler := http.HandlerFunc(endpoint.CreateProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Wrong response code, expected %d actual %d",
			http.StatusInternalServerError, rr.Code)
		t.Log(rr.Body)
	}
}

func TestNewKubeProfileHandler(t *testing.T) {
	svc := &Service{}
	h := NewHandler(svc)

	if h.service != svc {
		t.Errorf("Wrong service expected %v actual %v", svc, h.service)
	}
}

func TestHandler_Register(t *testing.T) {
	r := mux.NewRouter()
	h := Handler{}
	h.Register(r)
	expectedRouteCount := 3
	routes := []*mux.Route{}

	walkFn := func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		routes = append(routes, route)
		return nil
	}

	err := r.Walk(walkFn)

	if err != nil {
		t.Errorf("unexpected walk error %v", err)
	}

	if len(routes) != expectedRouteCount {
		t.Errorf("Wrong routes clount expected %d actual %d",
			expectedRouteCount, len(routes))
	}
}

func TestGetProfile(t *testing.T) {
	kubeProfile := &Profile{
		ID: "key",

		MasterProfiles: []NodeProfile{
			{"hello": "world"},
		},
		NodesProfiles: []NodeProfile{
			{"hello1": "world1"},
			{"hello2": "world2"},
		},

		K8SVersion:      "1.11.1",
		Provider:        clouds.AWS,
		Region:          "fra1",
		Arch:            "amd64",
		OperatingSystem: "linux",
		UbuntuVersion:   "xenial",
		DockerVersion:   "1.18.1",
		FlannelVersion:  "0.9.0",
		NetworkType:     "vxlan",
		CIDR:            "10.0.0.1/24",
		HelmVersion:     "0.11.1",
		RBACEnabled:     false,
	}
	data, _ := json.Marshal(kubeProfile)

	testCases := []struct {
		profileId     string
		profileData   []byte
		getProfileErr error
		expectedCode  int
	}{
		{
			expectedCode: http.StatusNotFound,
		},
		{
			profileData:   []byte(`{}`),
			getProfileErr: sgerrors.ErrNotFound,
			expectedCode:  http.StatusNotFound,
		},
		{
			profileId:    "profileId",
			profileData:  []byte(`{`),
			expectedCode: http.StatusInternalServerError,
		},
		{
			profileId:    "profileId",
			profileData:  []byte(`{`),
			expectedCode: http.StatusInternalServerError,
		},
		{
			profileId:    "profileId",
			profileData:  data,
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}
		mockRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(testCase.profileData, testCase.getProfileErr)
		svc := &Service{
			prefix:             "prefix",
			kubeProfileStorage: mockRepo,
		}
		h := Handler{
			service: svc,
		}

		router := mux.NewRouter()
		router.HandleFunc("/profile/{id}", h.GetProfile)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/profile/%s",
			testCase.profileId), nil)
		router.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestService_GetAll(t *testing.T) {
	testCases := []struct {
		repoErr      error
		getAllData   [][]byte
		expectedCode int
	}{
		{
			repoErr:      errors.New("unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			getAllData:   [][]byte{[]byte(`{`)},
			expectedCode: http.StatusInternalServerError,
		},
		{
			getAllData:   [][]byte{[]byte(`{}`)},
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}
		mockRepo.On("GetAll", mock.Anything, mock.Anything).
			Return(testCase.getAllData, testCase.repoErr)
		svc := &Service{
			prefix:             "prefix",
			kubeProfileStorage: mockRepo,
		}
		h := Handler{
			service: svc,
		}

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		h.GetProfiles(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}
