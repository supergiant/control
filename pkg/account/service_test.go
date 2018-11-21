package account

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/testutils"
)

func TestNewService(t *testing.T) {
	mockRepo := &testutils.MockStorage{}
	prefix := "prefix"

	svc := NewService(prefix, mockRepo)

	if svc == nil {
		t.Error("service must not me nil")
	}

	if svc.repository != mockRepo {
		t.Errorf("expected repo %v actual %v", mockRepo, svc.repository)
	}

	if prefix != svc.storagePrefix {
		t.Errorf("expected storage prefix %s actual %s", prefix, svc.storagePrefix)
	}
}

func TestServiceCreate(t *testing.T) {
	testCases := []struct {
		getResponse []byte
		getError    error
		account     *model.CloudAccount
		expectedErr error
	}{
		{
			getResponse: []byte(`{}`),
			getError:    sgerrors.ErrAlreadyExists,
			expectedErr: sgerrors.ErrAlreadyExists,
			account: &model.CloudAccount{
				Name:     "test",
				Provider: "nonameprovider",
			},
		},
		{
			getResponse: nil,
			getError:    sgerrors.ErrNotFound,
			account: &model.CloudAccount{
				Name:     "test",
				Provider: "nonameprovider",
			},
			expectedErr: ErrUnsupportedProvider,
		},
		{
			getResponse: nil,
			getError:    sgerrors.ErrNotFound,
			account: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
		{
			getResponse: nil,
			getError:    sgerrors.ErrNotFound,
			account: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
	}

	for _, testCase := range testCases {
		mockRepo := &testutils.MockStorage{}
		mockRepo.On("Get", mock.Anything,
			mock.Anything, mock.Anything).
			Return(testCase.getResponse, testCase.getError)

		svc := &Service{
			repository: mockRepo,
		}
		err := svc.Create(context.Background(), testCase.account)

		if err != nil {
			if errors.Cause(err).Error() != testCase.expectedErr.Error() {
				t.Errorf("expected error %v actual %v",
					testCase.expectedErr, err)
			}
		}

		if err == nil && testCase.expectedErr != nil {
			t.Error("error mut not be nil")
		}
	}
}
