package digitaloceansdk

import (
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"testing"
)

func TestNew(t *testing.T) {
	token := "test"
	sdk := New(token)

	if sdk == nil {
		t.Error("sdk must not be nil")
	}

	if sdk.accessToken != token {
		t.Errorf("Wrong access token expected %s actual %s", token, sdk.accessToken)
	}
}

func TestTokenSource_Token(t *testing.T) {
	testCases := []struct {
		ts          *TokenSource
		expectedErr error
	}{
		{
			ts: &TokenSource{
				AccessToken: "test",
			},
			expectedErr: nil,
		},
		{
			ts: &TokenSource{
				AccessToken: "",
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
	}

	for _, testCase := range testCases {
		token, err := testCase.ts.Token()

		if err != testCase.expectedErr {
			t.Errorf("expected error %v actual %v", testCase.expectedErr, err)
		}

		if testCase.expectedErr == nil && token.AccessToken != testCase.ts.AccessToken {
			t.Errorf("Wrong access token expected %s actual %s",
				testCase.ts.AccessToken, token.AccessToken)
		}
	}
}

func TestNewFromAccount(t *testing.T) {
	testCases := []struct {
		account     *model.CloudAccount
		expectedErr error
	}{
		{
			account: &model.CloudAccount{
				Credentials: map[string]string{},
			},
			expectedErr: sgerrors.ErrInvalidCredentials,
		},
		{
			account: &model.CloudAccount{
				Credentials: map[string]string{
					clouds.DigitalOceanAccessToken: "accessToken",
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		sdk, err := NewFromAccount(testCase.account)

		if err != testCase.expectedErr {
			t.Errorf("expected error %v actual %v", testCase.expectedErr, err)
			return
		}

		if err == nil && sdk == nil {
			t.Errorf("sdk must not be nil")
		}

		if sdk != nil && sdk.accessToken != testCase.account.Credentials[clouds.DigitalOceanAccessToken] {
			t.Errorf("Wrong access token expected %s actual %s",
				testCase.account.Credentials[clouds.DigitalOceanAccessToken], sdk.accessToken)
		}
	}
}

func TestSDKGetClient(t *testing.T) {
	sdk := SDK{
		accessToken: "test",
	}

	client := sdk.GetClient()

	if client == nil {
		t.Error("client must not be nil")
	}
}
