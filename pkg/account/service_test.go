package account

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/testutils"
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
		account     *model.CloudAccount
		expectedErr error
	}{
		{
			account: &model.CloudAccount{
				Provider: "nonameprovider",
			},
			expectedErr: ErrUnsupportedProvider,
		},
		{
			account: &model.CloudAccount{
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
		{
			account: &model.CloudAccount{
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
	}

	for _, testCase := range testCases {
		svc := &Service{}
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
