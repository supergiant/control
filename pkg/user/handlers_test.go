package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

type MockUserRepository struct {
	getUser func() (*User, error)

	storage map[string]*User
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]User, error) {
	return nil, nil
}

func (m *MockUserRepository) Get(ctx context.Context, userId string) (*User, error) {
	return m.getUser()
}

func (m *MockUserRepository) Create(context.Context, *User) error {
	return nil
}

func TestAuthHandler(t *testing.T) {
	ts := sgjwt.NewTokenService(60, []byte("secret"))

	testCases := []struct {
		user         *User
		ar           authRequest
		expectedCode int
	}{
		{
			&User{
				Login:    "user1",
				Password: "1234",
			},
			authRequest{
				UserName: "user1",
				Password: "1234",
			},
			http.StatusOK,
		},
		{
			&User{
				Login:    "user2",
				Password: "1234",
			},
			authRequest{
				UserName: "user1",
				Password: "12345",
			},
			http.StatusForbidden,
		},
	}

	for _, testCase := range testCases {
		err := testCase.user.encryptPassword()

		if err != nil {
			t.Error(err)
		}

		buf := bytes.NewBuffer([]byte(""))
		json.NewEncoder(buf).Encode(testCase.ar)

		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest("", "", buf)

		if err != nil {
			t.Error(err)
		}

		rec := httptest.NewRecorder()

		authHandler := &AuthHandler{
			tokenService: ts,
			userService: Service{
				Repository: &MockUserRepository{
					getUser: func() (*User, error) {
						return testCase.user, nil
					},
				},
			},
		}

		authHandler.ServeHTTP(rec, req)
	}
}
