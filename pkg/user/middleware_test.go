package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		err          error
		expectedCode int
		userId       string
	}{
		{
			nil,
			http.StatusOK,
			"login",
		},
		{
			nil,
			http.StatusForbidden,
			"",
		},
	}

	ts := sgjwt.NewTokenService(60, []byte("secret"))

	for _, testCase := range testCases {
		tokenString, _ := ts.Issue(testCase.userId)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest("", "", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		if err != nil {
			t.Error(err)
		}

		md := AuthMiddleware(ts, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		md.ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d", testCase.expectedCode, rec.Code)
		}
	}
}
