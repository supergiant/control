package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
)

type mockTokenIssuer struct {
	mock.Mock
}

func (m *mockTokenIssuer) Issue(userId string) (string, error) {
	args := m.Called(userId)
	val, ok := args.Get(0).(string)
	if !ok {
		return "", args.Error(1)
	}
	return val, args.Error(1)
}

func TestEndpoint_Authenticate(t *testing.T) {
	testCases := []struct {
		user            *User
		ar              []byte
		tokenIssueError error
		expectedCode    int
	}{
		{
			&User{
				Login:    "user1",
				Password: "1234",
			},
			[]byte(`{"login":"user1","password":"1234"}`),
			nil,
			http.StatusOK,
		},
		{
			&User{
				Login:    "user2",
				Password: "1234",
			},
			[]byte(`{"login":"user1","password":"12345"}`),
			nil,
			http.StatusForbidden,
		},
		{
			&User{
				Login:    "user1",
				Password: "1234",
			},
			[]byte(`{"login":"user1","password":"1234"}`),
			errors.New("test"),
			http.StatusInternalServerError,
		},
		{
			&User{
				Login:    "user1",
				Password: "1234",
			},
			[]byte(`"login":"user1","password":"1234"}`),
			nil,
			http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		storage := new(testutils.MockStorage)

		ts := &mockTokenIssuer{}
		ts.On("Issue", mock.Anything).
			Return("test", testCase.tokenIssueError)
		userEndpoint := NewHandler(NewService(DefaultStoragePrefix, storage), ts)
		handler := http.HandlerFunc(userEndpoint.Authenticate)

		err := testCase.user.encryptPassword()
		require.NoError(t, err)

		req, err := http.NewRequest("", "", bytes.NewReader(testCase.ar))
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
		user         []byte
		storageError error
		expectedCode int
	}{
		{
			user:         []byte(`{"login":"login","password": "password"}`),
			storageError: nil,
			expectedCode: http.StatusOK,
		},
		{
			user:         []byte(`{"login":"","password": ""}`),
			storageError: nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			user:         []byte(`{"login":"login","password": "password"}`),
			storageError: sgerrors.ErrAlreadyExists,
			expectedCode: http.StatusBadRequest,
		},
		{
			user:         []byte(`{"login":"login","password": "password"}`),
			storageError: errors.New("unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			user:         []byte(`"login":"login","password": "password"}`),
			storageError: errors.New("unknown error"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range tt {
		storage := new(testutils.MockStorage)
		userEndpoint := NewHandler(NewService(DefaultStoragePrefix, storage),
			jwt.NewTokenService(64, []byte("secret")))
		handler := http.HandlerFunc(userEndpoint.Create)

		storage.On("Get", mock.Anything,
			mock.Anything, mock.Anything).
			Return(nil, sgerrors.ErrNotFound)
		storage.On("Put", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(testCase.storageError)

		req, err := http.NewRequest("", "", bytes.NewReader(testCase.user))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, testCase.expectedCode, rec.Code)
	}
}
