package kube

import (
	"bytes"
	"strings"
	"testing"

	clientcmddapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
)

func TestIp2Host(t *testing.T) {
	testCases := []struct {
		ip   string
		host string
	}{
		{
			ip:   "172.16.10.2",
			host: "ip-172-16-10-2",
		},
		{
			ip:   "",
			host: "ip-",
		},
	}

	for _, testCase := range testCases {
		host := ip2Host(testCase.ip)

		if host != testCase.host {
			t.Errorf("Wrong host expected %s actual %s",
				testCase.host, host)
		}
	}
}

func TestProcessAWSMetrics(t *testing.T) {
	masters := map[string]*model.Machine{
		"master-1": {
			Name:      "Master-1",
			PrivateIp: "10.20.30.40",
		},
	}

	nodes := map[string]*model.Machine{
		"node-1": {
			Name:      "node-1",
			PrivateIp: "172.16.0.1",
		},
		"node-2": {
			Name:      "Node-2",
			PrivateIp: "172.16.0.2",
		},
	}

	k := &model.Kube{
		Provider: clouds.AWS,
		Masters:  masters,
		Nodes:    nodes,
	}

	metrics := map[string]map[string]interface{}{
		"ip-10-20-30-40": {},
		"ip-172-16-0-1":  {},
		"ip-172-16-0-2":  {},
	}

	processAWSMetrics(k, metrics)

	for _, masterNode := range masters {
		if _, ok := metrics[strings.ToLower(masterNode.Name)]; !ok {
			t.Errorf("Node %s not found in %v", masterNode.Name, metrics)
		}
	}

	for _, workerNode := range nodes {
		if _, ok := metrics[strings.ToLower(workerNode.Name)]; !ok {
			t.Errorf("Node %s not found in %v", workerNode.Name, metrics)
		}
	}
}

func TestKubeFromKubeConfig(t *testing.T) {
	testCases := []struct {
		description string
		kubeConfig  clientcmddapi.Config
		expectedErr string
	}{
		{
			description: "current context not found",
			kubeConfig: clientcmddapi.Config{
				Contexts:       map[string]*clientcmddapi.Context{},
				CurrentContext: "notFound",
			},
			expectedErr: "current context",
		},
		{
			description: "auth info not found",
			kubeConfig: clientcmddapi.Config{
				Contexts: map[string]*clientcmddapi.Context{
					"kubernetes": {
						AuthInfo: "not_found",
					},
				},
				AuthInfos:      map[string]*clientcmddapi.AuthInfo{},
				CurrentContext: "kubernetes",
			},
			expectedErr: "authInfo",
		},
		{
			description: "cluster not found",
			kubeConfig: clientcmddapi.Config{
				Contexts: map[string]*clientcmddapi.Context{
					"kubernetes": {
						AuthInfo: "kubernetes",
						Cluster:  "not_found",
					},
				},
				Clusters: map[string]*clientcmddapi.Cluster{},
				AuthInfos: map[string]*clientcmddapi.AuthInfo{
					"kubernetes": {
						ClientCertificateData: []byte(`client cert`),
						ClientKeyData:         []byte(`client key`),
					},
				},
				CurrentContext: "kubernetes",
			},
			expectedErr: "cluster",
		},
		{
			description: "success",
			kubeConfig: clientcmddapi.Config{
				Contexts: map[string]*clientcmddapi.Context{
					"admin@kubernetes": {
						AuthInfo: "kubernetes",
						Cluster:  "kubernetes",
					},
				},
				Clusters: map[string]*clientcmddapi.Cluster{
					"kubernetes": {
						CertificateAuthorityData: []byte(`ca cert`),
					},
				},
				AuthInfos: map[string]*clientcmddapi.AuthInfo{
					"kubernetes": {
						ClientCertificateData: []byte(`client cert`),
						ClientKeyData:         []byte(`client key`),
					},
				},
				CurrentContext: "admin@kubernetes",
			},
		},
	}

	for _, testCase := range testCases {
		kube, err := kubeFromKubeConfig(testCase.kubeConfig)

		if err == nil && testCase.expectedErr != "" {
			t.Error("Error must not be nil")
			continue
		}

		if err != nil && testCase.expectedErr != "" && !strings.Contains(err.Error(), testCase.expectedErr) {
			t.Errorf("Error message %s does not contain expectedc text %s", err.Error(), testCase.expectedErr)
			continue
		}

		if testCase.expectedErr == "" && kube == nil {
			t.Errorf("kube must not be nil")
			continue
		}

		if kube != nil && bytes.Compare(testCase.kubeConfig.AuthInfos["kubernetes"].ClientKeyData, []byte(kube.Auth.AdminKey)) != 0 {
			t.Errorf("Admin key does not match")
		}

		if kube != nil && bytes.Compare(testCase.kubeConfig.AuthInfos["kubernetes"].ClientCertificateData, []byte(kube.Auth.AdminCert)) != 0 {
			t.Errorf("Admin cert does not match")
		}

		if kube != nil && bytes.Compare(testCase.kubeConfig.Clusters["kubernetes"].CertificateAuthorityData, []byte(kube.Auth.CACert)) != 0 {
			t.Errorf("CA cert does not match")
		}
	}
}


func TestFindNextK8SVersion(t *testing.T) {
	testCases := []struct{
		description string
		current string
		version []string
		expected string
	}{
		{
			"success",
			"1.13.7",
			[]string{"1.11.5", "1.12.7", "1.13.7", "1.14.3"},
			"1.14.3",
		},
		{
			"version is too low",
			"1.9.7",
			[]string{"1.11.5", "1.12.7", "1.13.7", "1.14.3"},
			"",
		},
		{
			"empty current version",
			"",
			[]string{"1.11.5", "1.12.7", "1.13.7", "1.14.3"},
			"",
		},
		{
			"malformed current version",
			"1.1",
			[]string{"1.11.5", "1.12.7", "1.13.7", "1.14.3"},
			"",
		},
		{
			"empty versions list",
			"1.12.7",
			[]string{},
			"",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		actual := findNextMinorVersion(testCase.current, testCase.version)

		if !strings.EqualFold(actual, testCase.expected) {
			t.Errorf("Test %s has failed Wrong version expected %s actual %s",
				testCase.description, testCase.expected, actual)
		}
	}
}