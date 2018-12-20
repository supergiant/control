package util

import (
	"testing"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
)

func TestValidateCredentials(t *testing.T) {
	testCases := []struct {
		description   string
		cloudAccount  *model.CloudAccount
		getCreds      func(map[string]string) error
		expectedError error
	}{
		{
			description: "unsupported provider",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    "unknown",
				Credentials: map[string]string{},
			},
			expectedError: sgerrors.ErrUnsupportedProvider,
		},
		{
			description: "digitalocean",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			description: "digitalOCean wrong creds",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.DigitalOcean,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return sgerrors.ErrInvalidCredentials
			},
			expectedError: sgerrors.ErrInvalidCredentials,
		},
		{
			description: "aws",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			description: "gce",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.GCE,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			description: "gce invalid creds",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.GCE,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return sgerrors.ErrInvalidCredentials
			},
			expectedError: sgerrors.ErrInvalidCredentials,
		},
		{
			description: "aws invalid creads",
			cloudAccount: &model.CloudAccount{
				Name:        "test",
				Provider:    clouds.AWS,
				Credentials: map[string]string{},
			},
			getCreds: func(map[string]string) error {
				return sgerrors.ErrInvalidCredentials
			},
			expectedError: sgerrors.ErrInvalidCredentials,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		validator := CloudAccountValidatorImpl{
			digitalOcean: testCase.getCreds,
			aws:          testCase.getCreds,
			gce:          testCase.getCreds,
		}

		err := validator.ValidateCredentials(testCase.cloudAccount)

		if err != testCase.expectedError {
			t.Errorf("Expected error %v actual %v", testCase.expectedError, err)
		}
	}
}

func TestNewCloudAccountValidator(t *testing.T) {
	validator := NewCloudAccountValidator()

	if validator == nil {
		t.Errorf("validator must not be nil")
	}

	if validator.aws == nil {
		t.Errorf("aws must not be nil")
	}

	if validator.gce == nil {
		t.Errorf("gce must not be nil")
	}

	if validator.digitalOcean == nil {
		t.Errorf("digitalocean must not be nil")
	}
}
