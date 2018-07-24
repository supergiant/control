package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
)

func TestEndpoint_Authenticate(t *testing.T) {
	ts := jwt.NewTokenService(60, []byte("secret"))
	testCases := []struct {
		user         *User
		ar           AuthRequest
		expectedCode int
	}{
		{
			&User{
				Login:    "user1",
				Password: "1234",
			},
			AuthRequest{
				Login:    "user1",
				Password: "1234",
			},
			http.StatusOK,
		},
		{
			&User{
				Login:    "user2",
				Password: "1234",
			},
			AuthRequest{
				Login:    "user1",
				Password: "12345",
			},
			http.StatusForbidden,
		},
	}
	storage := new(testutils.MockStorage)

	userEndpoint := NewHandler(NewService(DefaultStoragePrefix, storage), ts)
	handler := http.HandlerFunc(userEndpoint.Authenticate)

	for _, testCase := range testCases {
		err := testCase.user.encryptPassword()
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte(""))
		json.NewEncoder(buf).Encode(testCase.ar)

		req, err := http.NewRequest("", "", buf)
		require.NoError(t, err)

		storage.On("Get", mock.Anything,
			mock.Anything,
			mock.Anything).Return(userToJSON(testCase.user), nil)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, testCase.expectedCode, rec.Code)
	}
}

func userToJSON(user *User) (data []byte) {
	data, _ = json.Marshal(user)
	return
}

func TestEndpoint_Create(t *testing.T) {
	tt := []struct {
		user         *User
		expectedCode int
	}{
		{
			user: &User{
				Login:    "login",
				Password: "password",
			},
			expectedCode: http.StatusOK,
		},
		{
			user: &User{
				Login:    "",
				Password: "",
			},
			expectedCode: http.StatusBadRequest,
		},
	}
	storage := new(testutils.MockStorage)
	userEndpoint := NewHandler(NewService(DefaultStoragePrefix, storage), jwt.NewTokenService(64, []byte("secret")))
	handler := http.HandlerFunc(userEndpoint.Create)

	storage.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, sgerrors.ErrNotFound)
	storage.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	for _, testCase := range tt {
		req, err := http.NewRequest("", "", bytes.NewReader(userToJSON(testCase.user)))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, testCase.expectedCode, rec.Code)
	}
}
