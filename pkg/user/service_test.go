package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
	"testing"
)

func TestNewService(t *testing.T) {
	prefix := "prefix"
	repo := &testutils.MockStorage{}
	svc := NewService(prefix, repo)

	if svc.storagePrefix != prefix {
		t.Errorf("Wrong prefix expected %s actual %s", prefix, svc.storagePrefix)
	}

	if svc.repository != repo {
		t.Errorf("Wrong repo expected %v actual %v", repo, svc.repository)
	}
}

func TestService_Create(t *testing.T) {
	err := errors.New("put error")
	testCases := []struct {
		user         *User
		getError     error
		putError     error
		serviceError error
	}{
		{
			user:         nil,
			getError:     sgerrors.ErrNilValue,
			serviceError: sgerrors.ErrNilValue,
		},
		{
			user: &User{
				Login:    "user",
				Password: "1234",
			},
			getError:     errors.New("test"),
			serviceError: sgerrors.ErrAlreadyExists,
		},
		{
			user: &User{
				Login:    "user",
				Password: "1234",
			},
			getError:     nil,
			putError:     err,
			serviceError: err,
		},
		{
			user: &User{
				Login:    "user",
				Password: "1234",
			},
		},
	}

	for _, testCase := range testCases {
		storage := &testutils.MockStorage{}
		storage.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(mock.Anything, testCase.getError)
		storage.On("Put", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(testCase.putError)
		service := &Service{
			storagePrefix: "prefix",
			repository:    storage,
		}

		err := service.Create(context.Background(), testCase.user)

		if err != nil && err != testCase.serviceError {
			t.Errorf("Service has returned wrong error expected %v actual %v",
				testCase.serviceError, err)
		}
	}
}

func TestService_Authenticate(t *testing.T) {
	err := errors.New("unknown error")
	testCases := []struct {
		user       *User
		userData   []byte
		user2      *User
		repoErr    error
		serviceErr error
	}{
		{
			user: &User{
				Login:    "",
				Password: "1234",
			},
			serviceErr: sgerrors.ErrInvalidCredentials,
		},
		{
			user: &User{
				Login:    "root",
				Password: "",
			},
			serviceErr: sgerrors.ErrInvalidCredentials,
		},
		{
			user: &User{
				Login:    "root",
				Password: "1234",
			},
			repoErr:    sgerrors.ErrNotFound,
			serviceErr: sgerrors.ErrNotFound,
		},
		{
			user: &User{
				Login:    "root",
				Password: "1234",
			},
			repoErr:    err,
			serviceErr: err,
		},
		{
			user: &User{
				Login:    "root",
				Password: "1234",
			},
			userData:   []byte(`{"login"}`),
			serviceErr: sgerrors.ErrInvalidJson,
		},
		{
			user: &User{
				Login:    "root",
				Password: "1234",
			},
			userData:   []byte(`{"login":"root","password":"1234"}`),
			serviceErr: sgerrors.ErrInvalidJson,
		},
		{
			user: &User{
				Login:    "root",
				Password: "1234",
			},
			user2:    &User{Login: "root", Password: "1234"},
			userData: []byte(`{"login":"root","encrypted_password":"$2a$10$JM.SQtGGoUwJ/.3NOHbL6e6Eb5lQi9vsCaaXkiqgu4ZHG95JY/Q.y"}`),

			serviceErr: nil,
		},
	}

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}

		if testCase.user2 != nil {
			testCase.user2.encryptPassword()

			data, _ := json.Marshal(testCase.user2)
			mockRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(data, nil)
		} else {
			mockRepo.On("Get", mock.Anything,
				mock.Anything, mock.Anything).Return(mock.Anything, testCase.repoErr)
		}

		svc := Service{
			"prefix",
			mockRepo,
		}

		err := svc.Authenticate(context.Background(),
			testCase.user.Login, testCase.user.Password)

		if err != testCase.serviceErr {
			t.Errorf("Wrong error expected %v actual %v",
				testCase.serviceErr, err)
		}
	}
}

func TestService_GetAll(t *testing.T) {
	testCases := []struct {
		repoData      [][]byte
		expectedErr   error
		expectedCount int
	}{
		{
			repoData:    [][]byte{[]byte(`{`)},
			expectedErr: sgerrors.ErrInvalidJson,
		},
		{
			repoData:      [][]byte{[]byte(`{}`), []byte(`{}`)},
			expectedCount: 2,
		},
	}

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}
		mockRepo.On("GetAll", mock.Anything, mock.Anything).
			Return(testCase.repoData, testCase.expectedErr)

		svc := Service{
			repository: mockRepo,
		}

		data, err := svc.GetAll(context.Background())

		if testCase.expectedErr != err {
			t.Errorf("Wrong error expected %s actual %v",
				testCase.expectedErr, err)
		}

		if testCase.expectedErr == nil && len(data) != testCase.expectedCount {
			t.Errorf("Wrong count expected %d actual %d",
				testCase.expectedCount, len(data))
		}
	}
}
