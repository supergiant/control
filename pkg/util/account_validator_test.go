package util

import (
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"testing"
)

func TestValidateCredentials(t *testing.T) {
	testCases := []struct {
		cloudAccount  *model.CloudAccount
		digitalOcean  func(map[string]string) error
		aws           func(map[string]string) error
		expectedError error
	}{
		{
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    "unknown",
				Credentials: map[string]string{},
			},
			expectedError: sgerrors.ErrUnsupportedProvider,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			digitalOcean: func(map[string]string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			digitalOcean: func(map[string]string) error {
				return sgerrors.ErrInvalidCredentials
			},
			expectedError: sgerrors.ErrInvalidCredentials,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			aws: func(map[string]string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			aws: func(map[string]string) error {
				return sgerrors.ErrInvalidCredentials
			},
			expectedError: sgerrors.ErrInvalidCredentials,
		},
	}

	for _, testCase := range testCases {
		validator := CloudAccountValidatorImpl{
			digitalOcean: testCase.digitalOcean,
			aws:          testCase.aws,
		}

		err := validator.ValidateCredentials(testCase.cloudAccount)

		if err != testCase.expectedError {
			t.Errorf("Expected error %v actual %v", testCase.expectedError, err)
		}
	}
}
