package util

import (
	"context"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestRandomStringLen(t *testing.T) {
	testCases := []int{4, 8, 16}

	for _, testCase := range testCases {
		rndString := RandomString(testCase)

		if len(rndString) != testCase {
			t.Errorf("Wrong random string size expected %d actual %d", testCase, len(rndString))
		}
	}
}

func TestRandomStringUnique(t *testing.T) {
	m := make(map[string]struct{})
	count := 1000
	size := 8

	for i := 0; i < count; i++ {
		s := RandomString(size)

		if _, ok := m[s]; ok {
			t.Errorf("Duplicate string")
			return
		}
	}
}

func TestMakeNodeName(t *testing.T) {
	testCases := []struct {
		role        bool
		clusterName string
		taskId      string
		expected    string
	}{
		{
			true,
			"hello",
			"5678",
			"hello-master-5678",
		},
		{
			false,
			"world",
			"1234",
			"world-node-1234",
		},
	}

	for _, testCase := range testCases {
		nodeName := MakeNodeName(testCase.clusterName, testCase.taskId, testCase.role)

		if !strings.EqualFold(nodeName, testCase.expected) {
			t.Errorf("Wrong node clusterName expected %s actual %s",
				testCase.expected, nodeName[:len(nodeName)-5])
		}
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
					"accessToken": "abcd",
				},
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		config := &steps.Config{
			CloudAccountName: testCase.cloudAccount.Name,
		}

		FillCloudAccountCredentials(context.Background(), testCase.cloudAccount, config)

		if !strings.EqualFold(testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken) {
			t.Errorf("Wrong access token expected %s actual %s",
				testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken)
		}
	}
}
