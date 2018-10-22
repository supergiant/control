package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		description  string
		err          error
		expectedCode int
		authHeader   string
		userId       string
		issuer       func(string) (string, error)
	}{
		{
			"token ok",
			nil,
			http.StatusOK,
			"Bearer %s",
			"login",
			func(userId string) (string, error) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
					"accesses":   []string{"edit", "view"},
					"user_id":    userId,
					"issued_at":  time.Now().Unix(),
					"expires_at": time.Now().Unix() + 600,
				})

				tokenString, err := token.SignedString([]byte("secret"))

				if err != nil {
					return "", err
				}

				return tokenString, nil
			},
		},
		{
			"token expired",
			nil,
			http.StatusForbidden,
			"Bearer %s",
			"",
			func(userId string) (string, error) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
					"accesses":   []string{"edit", "view"},
					"user_id":    userId,
					"issued_at":  time.Now().Unix() - 2,
					"expires_at": time.Now().Unix() - 1,
				})

				tokenString, err := token.SignedString([]byte("secret"))

				if err != nil {
					return "", err
				}

				return tokenString, nil
			},
		},
		{
			"user id empty",
			nil,
			http.StatusForbidden,
			"Bearer %s",
			"",
			func(userId string) (string, error) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
					"accesses":   []string{"edit", "view"},
					"user_id":    userId,
					"issued_at":  time.Now().Unix(),
					"expires_at": time.Now().Unix() + 60,
				})

				tokenString, err := token.SignedString([]byte("secret"))

				if err != nil {
					return "", err
				}

				return tokenString, nil
			},
		},
		{
			"empty header",
			nil,
			http.StatusForbidden,
			"",
			"root",
			func(userId string) (string, error) {
				return "", nil
			},
		},
		{
			"bad claims",
			nil,
			http.StatusForbidden,
			"",
			"root",
			func(userId string) (string, error) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
					"user_id": 42,
				})

				tokenString, err := token.SignedString([]byte("secret"))

				if err != nil {
					return "", err
				}

				return tokenString, nil
			},
		},
	}

	ts := sgjwt.NewTokenService(60, []byte("secret"))

	for _, testCase := range testCases {
		tokenString, _ := ts.Issue(testCase.userId)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest("", "", nil)
		req.Header.Set("Authorization", fmt.Sprintf(testCase.authHeader, tokenString))

		if err != nil {
			t.Error(err)
		}

		md := Middleware{
			TokenService: ts,
		}

		md.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rec, req)

		if rec.Code != testCase.expectedCode {
			t.Errorf("Wrong response code expected %d actual %d", testCase.expectedCode, rec.Code)
		}
	}
}

type testHandler struct {
	called bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
}

func TestContentTypeJSON(t *testing.T) {
	router := mux.NewRouter()
	h := &testHandler{}
	router.Handle("/", ContentTypeJSON(h))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(rec, req)

	if h := rec.Header().Get("Content-Type"); !strings.EqualFold(h, "application/json") {
		t.Errorf("Wrong content type value expected %s actual %s", "application/json", h)
	}

	if h.called != true {
		t.Error("json middleware was not called")
	}
}
