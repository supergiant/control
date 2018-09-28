package controlplane

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
