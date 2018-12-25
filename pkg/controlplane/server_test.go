package controlplane

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"time"
	"strings"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		cfg          *Config
		headers      map[string]string
		method       string
		expectedCode int
	}{
		{
			cfg: &Config{},
			headers: map[string]string{
				"Access-Control-Request-Headers": "something",
				"Access-Control-Request-Method":  "something",
				"Authorization":                  "Bearer token",
				"Origin":                         "localhost",
			},
			method:       http.MethodOptions,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			cfg: &Config{},
			headers: map[string]string{
				"Authorization": "Bearer token",
				"Origin":        "localhost",
			},
			method:       http.MethodOptions,
			expectedCode: http.StatusBadRequest,
		},
		{
			cfg: &Config{},
			headers: map[string]string{
				"Access-Control-Request-Headers": "something",
				"Authorization":                  "Bearer token",
			},
			method:       http.MethodDelete,
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		router := mux.NewRouter()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		})

		server := NewServer(router, testCase.cfg)
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(testCase.method, "/", nil)

		if err != nil {
			t.Errorf("create request %v", err)
		}
		for k, v := range testCase.headers {
			req.Header.Set(k, v)
		}

		// Allow localhost as an origin
		origins := handlers.AllowedOrigins([]string{"*"})
		server.server.Handler = handlers.CORS(origins)(server.server.Handler)

		server.server.Handler.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("unexpected response code expected %d actual %d",
				testCase.expectedCode, rec.Code)
		}
	}
}

func TestTrimPrefix(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "/hello",
			output: "/",
		},
		{
			input:  "/static/vendor.js",
			output: "static/vendor.js",
		},
	}

	for _, testCase := range testCases {
		called := false
		actualURL := ""

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			actualURL = r.URL.Path
		})

		h2 := trimPrefix(h)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.input, nil)
		h2.ServeHTTP(rec, req)

		if !called {
			t.Error("Handler has not been called")
		}

		if actualURL != testCase.output {
			t.Errorf("url must be empty after trimming prefix actual %s", actualURL)
		}
	}
}

func TestConfigureApp(t *testing.T) {
	config := &Config{
		PprofListenStr: ":9090",
		TemplatesDir:   "../../templates",
		UiDir:          "../../cmd/ui",
		SpawnInterval: time.Second * 5,
	}

	router, err := configureApplication(config)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if router == nil {
		t.Errorf("router must not be nil")
	}
}

func TestNewVersionHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/version", nil)
	version := "2.0.0"

	h := NewVersionHandler(version)

	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Wrong response code expected %d actual %d",
			http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), version) {
		t.Errorf("Version %s not found in response body %s",
			rec.Body.String(), version)
	}
}