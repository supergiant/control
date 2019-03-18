// +build dev

package ui

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

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
