package helm

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
	"strings"
)

func TestService_Create(t *testing.T) {
	svc := &Service{}
	expectedError := sgerrors.ErrNotFound

	err := svc.Create(context.Background(), nil)

	if err != expectedError {
		t.Errorf("expected error %v actual %v", expectedError, err)
	}
}

func TestService_Get(t *testing.T) {
	mockStorage := &testutils.MockStorage{}
	err := errors.New("error")
	mockStorage.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(mock.Anything, err)

	svc := &Service{
		storage: mockStorage,
	}

	_, err2 := svc.Get(context.Background(), "repo_name")

	if errors.Cause(err2) != err {
		t.Errorf("expected error %v actual %v", err, err2)
	}

}

func TestService_GetAll(t *testing.T) {
	testCases := []struct {
		rawData    [][]byte
		storageErr error

		expectedResult []helm.Repository
		expectedErr    error
	}{
		{
			storageErr: errors.New("storage error"),
		},
		{
			rawData:     [][]byte{[]byte(`{`)},
			storageErr:  nil,
			expectedErr: errors.New("unmarshal"),
		},
		{
			rawData: [][]byte{[]byte(`{"name":"test", "url":"http://helm.com"}`)},
			expectedResult: []helm.Repository{{
				Name: "test",
				URL:  "http://helm.com",
			}},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		mockStorage := &testutils.MockStorage{}
		err := errors.New("error")
		mockStorage.On("GetAll", mock.Anything, mock.Anything).
			Return(testCase.rawData, testCase.storageErr)

		svc := &Service{
			storage: mockStorage,
		}

		repos, err := svc.GetAll(context.Background())

		if testCase.expectedErr != nil && err != nil &&
			!strings.Contains(err.Error(), testCase.expectedErr.Error()) {
			t.Errorf("wrong error expected  %v actual %v",
				testCase.expectedErr, err)
		}

		if len(repos) != len(testCase.expectedResult) {
			t.Errorf("Wrong counf of repos expected %d actual %d",
				len(testCase.expectedResult), len(repos))
		}

		if len(repos) > 0 && repos[0].Name != "test" && repos[0].URL != "http://helm.com" {
			t.Errorf("wrong repo name expected test url http://helm.com actual %v",
				repos[0])
		}
	}
}
