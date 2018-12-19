package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"fmt"
	sgjwt "github.com/supergiant/control/pkg/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		description  string
		err          error
		expectedCode int
		authHeader   string
		query        string
		userId       string
		issuer       func(string) (string, error)
	}{
		{
			description:  "token ok",
			err:          nil,
			expectedCode: http.StatusOK,
			authHeader:   "Bearer %s",
			userId:       "login",
			issuer: func(userId string) (string, error) {
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
			description:  "token from query",
			err:          nil,
			expectedCode: http.StatusOK,
			userId:       "login",
			query:        "/url?token=%s",
			issuer: func(userId string) (string, error) {
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
			description:  "token expired",
			err:          nil,
			expectedCode: http.StatusForbidden,
			authHeader:   "Bearer %s",
			userId:       "",
			issuer: func(userId string) (string, error) {
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
			description:  "user id empty",
			err:          nil,
			expectedCode: http.StatusForbidden,
			authHeader:   "Bearer %s",
			issuer: func(userId string) (string, error) {
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
			description:  "no token provided",
			err:          nil,
			expectedCode: http.StatusForbidden,
			userId:       "root",
			issuer: func(userId string) (string, error) {
				return "", nil
			},
		},
		{
			description:  "no enough data in token",
			err:          nil,
			expectedCode: http.StatusForbidden,
			userId:       "root",
			authHeader:   "%s",
			issuer: func(userId string) (string, error) {
				return "", nil
			},
		},
		{
			description:  "bad claims",
			err:          nil,
			expectedCode: http.StatusForbidden,
			userId:       "root",
			authHeader:   "Bearer %s",
			issuer: func(userId string) (string, error) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
					"user_id": "user",
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
		t.Log(testCase.description)
		tokenString, _ := testCase.issuer(testCase.userId)

		url := ""

		if testCase.query != "" {
			url = fmt.Sprintf(testCase.query, tokenString)
		}

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, url, nil)

		if testCase.authHeader != "" {
			req.Header.Set("Authorization",
				fmt.Sprintf(testCase.authHeader, tokenString))
		}

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
