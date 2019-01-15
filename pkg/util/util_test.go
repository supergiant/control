package util

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
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
		testName     string
		cloudAccount *model.CloudAccount
		err          error
	}{
		{
			testName: "digital ocean",
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.DigitalOcean,
				Credentials: map[string]string{
					"accessToken": "abcd",
					"publicKey":   "test-public-key",
				},
			},
			err: nil,
		},
		{
			testName: "aws",
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.AWS,
				Credentials: map[string]string{
					"access_key":  "1",
					"secret_key":  "secret-key",
					"keyPairName": "my-key-pair",
					"publicKey":   "test-public-key",
				},
			},
			err: nil,
		},
		{
			testName: "aws",
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: clouds.GCE,
				Credentials: map[string]string{
					"projectId":   "ordinal-case-222023",
					"privateKey":  "-----BEGIN PRIVATE KEY-----\n\n-----END PRIVATE KEY-----\n",
					"clientEmail": "myemail@gmail.comn",
					"tokenURI":    "https://oauth2.googleapis.com/token",
					"publicKey":   "ssh-rsa  myemail@gmail.com",
				},
			},
			err: nil,
		},
		{
			testName: "unknown provider",
			cloudAccount: &model.CloudAccount{
				Name:     "testName",
				Provider: "unknown",
			},
			err: sgerrors.ErrUnknownProvider,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.testName)
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
			if !strings.EqualFold(testCase.cloudAccount.Credentials["access_key"], config.AWSConfig.KeyID) {
				t.Errorf("Wrong key id expected %s actual %s",
					testCase.cloudAccount.Credentials["access_key"], config.AWSConfig.KeyID)
			}

			if !strings.EqualFold(testCase.cloudAccount.Credentials["secret_key"], config.AWSConfig.Secret) {
				t.Errorf("Wrong secret expected %s actual %s",
					testCase.cloudAccount.Credentials["secret_key"], config.AWSConfig.Secret)
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
	testCases := []struct {
		keyName        string
		isUser         bool
		expectedResult string
	}{
		{
			keyName:        "test",
			isUser:         true,
			expectedResult: fmt.Sprintf("%s-user", "test"),
		},
		{
			keyName:        "",
			isUser:         false,
			expectedResult: "-provision",
		},
		{
			keyName:        "test",
			isUser:         false,
			expectedResult: fmt.Sprintf("%s-provision", "test"),
		},
	}

	for _, testCase := range testCases {
		actual := MakeKeyName(testCase.keyName, testCase.isUser)

		if !strings.Contains(actual, testCase.expectedResult) {
			t.Errorf("Wrong key name expected %s actual %s",
				testCase.expectedResult, actual)
		}
	}
}

func TestMakeRole(t *testing.T) {
	testCases := []bool{true, false}

	for _, testCase := range testCases {
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
	testCases := []struct {
		input  map[string]string
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
			&struct {
				Key string `json:"key"`
			}{},
			"",
		},
		{
			map[string]string{
				"key": "value",
			},
			&struct {
				Key string `json:"key"`
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

	testCases := []struct {
		nodeMap      map[string]*node.Node
		expectedNode *node.Node
	}{
		{
			nodeMap: map[string]*node.Node{
				"node-1": n,
			},
			expectedNode: n,
		},
		{
			nodeMap:      map[string]*node.Node{},
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

func TestGetWriter(t *testing.T) {
	testCases := []struct {
		name   string
		hasErr bool
	}{
		{
			name:   "test.txt",
			hasErr: false,
		},
		{
			name:   "",
			hasErr: true,
		},
	}

	for _, testCase := range testCases {
		writer, err := GetWriter(testCase.name)

		if err == nil && testCase.hasErr {
			t.Errorf("error must not be nil")
		}

		if testCase.hasErr && writer == nil {
			t.Errorf("Writer must not be nil")
		}
	}
}

func TestLoadCloudSpecificDataFromKube(t *testing.T) {
	testCases := []struct {
		description string
		kube        *model.Kube
		provider    clouds.Name
		hasErr      bool
	}{
		{
			description: "digitalocean",
			kube: &model.Kube{
				BootstrapPrivateKey: []byte(`private-key`),
				BootstrapPublicKey:  []byte(`public-key`),
				SshPublicKey:        []byte(`public-key2`),

				Region: "fra-1",
			},
			provider: clouds.DigitalOcean,
		},
		{
			description: "gce",
			kube: &model.Kube{
				BootstrapPrivateKey: []byte(`private-key`),
				BootstrapPublicKey:  []byte(`public-key`),
				SshPublicKey:        []byte(`public-key2`),
			},
			provider: clouds.GCE,
		},
		{
			description: "aws",
			kube: &model.Kube{
				BootstrapPrivateKey: []byte(`private-key`),
				BootstrapPublicKey:  []byte(`public-key`),
				SshPublicKey:        []byte(`public-key2`),
				CloudSpec: map[string]string{
					clouds.AwsImageID:               "imageId",
					clouds.AwsVpcID:                 "vpcId",
					clouds.AwsNodeInstanceProfile:   "nodeProfile",
					clouds.AwsMasterInstanceProfile: "masterProfile",
					clouds.AwsInternetGateWayID:     "internetGWId",
					clouds.AwsNodesSecgroupID:       "nodesSecurityId",
					clouds.AwsMastersSecGroupID:     "masterSecGroup",
					clouds.AwsAZ:                    "az",
					clouds.AwsRouteTableID:          "routetableid",
					clouds.AWSAccessKeyID:           "accessKey",
					clouds.AWSSecretKey:             "secretKey",
					clouds.AwsKeyPairName:           "keyPairName",
				},
			},
			provider: clouds.AWS,
		},
		{
			description: "unsupported",
			kube: &model.Kube{
				BootstrapPrivateKey: []byte(`private-key`),
				BootstrapPublicKey:  []byte(`public-key`),
				SshPublicKey:        []byte(`public-key2`),
			},
			provider: clouds.Name("unsupported"),
			hasErr:   true,
		},
		{
			description: "nil value",
			hasErr:      true,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		config := &steps.Config{
			Provider: testCase.provider,
		}
		err := LoadCloudSpecificDataFromKube(testCase.kube, config)

		if testCase.hasErr && err == nil {
			t.Errorf("Error must not be nil")
		}

		if !testCase.hasErr && err != nil {
			t.Errorf("unexpected error %v", err)
		}
	}
}
