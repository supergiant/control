package workflows

import (
	"context"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type mockCloudAccountService struct {
	cloudAccount *model.CloudAccount
	err          error
}

func (m *mockCloudAccountService) Get(ctx context.Context, name string) (*model.CloudAccount, error) {
	return m.cloudAccount, m.err
}

func TestBindParams(t *testing.T) {
	obj := &struct {
		ParamA string `json:"a"`
		ParamB string `json:"b"`
	}{}

	expectedA := "hello"
	expectedB := "world"

	params := map[string]string{
		"a": expectedA,
		"b": expectedB,
	}

	util.BindParams(params, obj)

	if !strings.EqualFold(obj.ParamA, params["a"]) {
		t.Errorf("Wrong value for paramA expected %s actual %s", params["a"], obj.ParamA)
	}

	if !strings.EqualFold(obj.ParamB, params["b"]) {
		t.Errorf("Wrong value for paramB expected %s actual %s", params["b"], obj.ParamB)
	}
}

// TODO(stgleb): extend for other types of cloud providers
func TestFillCloudAccountCredentials(t *testing.T) {
	testCases := []struct {
		cloudAccount *model.CloudAccount
		err          error
	}{
		{
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"name":         "hello_world",
					"k8sVersion":   "",
					"region":       "",
					"size":         "",
					"role":         "",
					"image":        "",
					"fingerprints": "fingerprint",
					"accessToken":  "abcd",
				},
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		mock := &mockCloudAccountService{
			testCase.cloudAccount,
			testCase.err,
		}

		config := &steps.Config{
			CloudAccountName: testCase.cloudAccount.Name,
		}

		FillCloudAccountCredentials(context.Background(), mock, config)

		if !strings.EqualFold(testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken) {
			t.Errorf("Wrong access token expected %s actual %s",
				testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken)
		}

		if !strings.EqualFold(testCase.cloudAccount.Credentials["name"], config.DigitalOceanConfig.Name) {
			t.Errorf("Wrong cloud account name expected %s actual %s",
				testCase.cloudAccount.Credentials["name"], config.DigitalOceanConfig.Name)
		}
	}
}
