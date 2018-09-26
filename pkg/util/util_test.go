package util

import (
	"context"
	"strings"
	"testing"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"bytes"
	"fmt"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/node"
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
					"publicKey": "test-public-key",
				},
			},
			err: nil,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.AWS,
				Credentials: map[string]string{
					"keyID": "1",
					"secret": "secret-key",
					"keyPairName": "my-key-pair",
					"publicKey": "test-public-key",
				},
			},
			err: nil,
		},
		{
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: "unknown",
			},
			err: sgerrors.ErrUnknownProvider,
		},
	}

	for _, testCase := range testCases {
		config := &steps.Config{
			CloudAccountName: testCase.cloudAccount.Name,
		}

		err := FillCloudAccountCredentials(context.Background(), testCase.cloudAccount, config)

		if testCase.cloudAccount.Provider == clouds.DigitalOcean {
			if !strings.EqualFold(testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken) {
				t.Errorf("Wrong access token expected %s actual %s",
					testCase.cloudAccount.Credentials["accessToken"], config.DigitalOceanConfig.AccessToken)
			}
		}

		if testCase.cloudAccount.Provider == clouds.AWS {
			if !strings.EqualFold(testCase.cloudAccount.Credentials["keyID"], config.AWSConfig.KeyID) {
				t.Errorf("Wrong key id expected %s actual %s",
					testCase.cloudAccount.Credentials["keyID"], config.AWSConfig.KeyID)
			}

			if !strings.EqualFold(testCase.cloudAccount.Credentials["secret"], config.AWSConfig.Secret) {
				t.Errorf("Wrong secret expected %s actual %s",
					testCase.cloudAccount.Credentials["secret"], config.AWSConfig.Secret)
			}

			if !strings.EqualFold(testCase.cloudAccount.Credentials["keyPairName"], config.AWSConfig.KeyPairName) {
				t.Errorf("Wrong keyPairName expected %s actual %s",
					testCase.cloudAccount.Credentials["keyPairName"], config.AWSConfig.KeyPairName)
			}
		}

		if config.SshConfig.PublicKey != testCase.cloudAccount.Credentials["publicKey"] {
			t.Errorf("PublicKey %s not found in credentials %v",
				testCase.cloudAccount.Credentials["publicKey"], config.SshConfig.PublicKey)
		}

		if err != testCase.err {
			t.Errorf("expected error %v actual %v", testCase.err, err)
		}
	}
}

func TestGetLogger(t *testing.T) {
	writer := &bytes.Buffer{}
	logger := GetLogger(writer)

	if logger.Out != writer {
		t.Errorf("Wrong output writer expected %v actual %v",
			writer, logger.Out)
	}
}

func TestMakeFileName(t *testing.T) {
	taskId := "1234abcd"
	fileName := MakeFileName(taskId)

	if !strings.Contains(fileName, taskId) {
		t.Errorf("file name %s must contain %s", fileName, taskId)
	}
}

func TestMakeKeyName(t *testing.T) {
	testCases := []struct{
		keyName string
		isUser bool
		expectedResult string
	}{
		{
			keyName: "test",
			isUser: true,
			expectedResult: fmt.Sprintf("%s-user", "test"),
		},
		{
			keyName: "test",
			isUser: false,
			expectedResult: fmt.Sprintf("%s-provision", "test"),
		},
	}

	for _, testCase := range testCases {
		actual := MakeKeyName(testCase.keyName, testCase.isUser)

		if !strings.EqualFold(actual, testCase.expectedResult) {
			t.Errorf("Wrong key name expected %s actual %s",
				testCase.expectedResult, actual)
		}
	}
}

func TestMakeRole(t *testing.T) {
	testCases := []bool{true, false}

	for _, testCase :=  range testCases {
		role := MakeRole(testCase)

		if testCase && !strings.EqualFold(role, string(node.RoleMaster)) {
			t.Errorf("Wrong role expected %s actual %s", node.RoleMaster, role)
		}

		if !testCase && !strings.EqualFold(role, string(node.RoleNode)) {
			t.Errorf("Wrong role expected %s actual %s", node.RoleNode, role)
		}
	}
}

func TestBindParams(t *testing.T) {
	testCases := []struct{
		input map[string]string
		output interface{}
		errMsg string
	}{
		{
			nil,
			nil,
			"Unmarshal",
		},
		{
			map[string]string{
				"key": "value",
			},
			nil,
			"Unmarshal",
		},
		{
			nil,
			&struct{
				Key string`json:"key"`
			}{},
			"",
		},
		{
			map[string]string{
				"key": "value",
			},
			&struct{
				Key string`json:"key"`
			}{},
			"",
		},
	}

	for _, testCase := range testCases {
		err := BindParams(testCase.input, testCase.output)

		if len(testCase.errMsg) == 0 && err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if len(testCase.errMsg) > 0 && !strings.Contains(err.Error(), testCase.errMsg) {
			t.Errorf("expected error message to have %s actual %s", testCase.errMsg, err.Error())
		}
	}
}

func TestGetRandomNode(t *testing.T) {
	n := &node.Node{}

	testCases := []struct{
		nodeMap map[string]*node.Node
		expectedNode *node.Node
	}{
		{
			nodeMap: map[string]*node.Node{
				"node-1": n,
			},
			expectedNode: n,
		},
		{
			nodeMap: map[string]*node.Node{},
			expectedNode: nil,
		},
	}
	
	for _, testCase := range testCases {
		actual := GetRandomNode(testCase.nodeMap)

		if actual != testCase.expectedNode {
			t.Errorf("expected node %v actual %v",
				testCase.expectedNode, actual)
		}
	}
}