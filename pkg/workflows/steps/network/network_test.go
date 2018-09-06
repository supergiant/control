package network

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestNetworkConfig(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	testCases := []struct {
		etcdRepositoryUrl string
		etcdVersion       string
		etcdHost          string
		arch              string
		operatingSystem   string
		network           string
		networkType       string
		expectedError     error
	}{
		{
			"https://github.com/coreos/etcd/releases/download",
			"0.9.0",
			"10.20.30.40",
			"amd64",
			"linux",
			"10.0.2.0/24",
			"vxlan",
			nil,
		},
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			errors.New("error has occurred"),
		},
	}

	for _, testCase := range testCases {
		r := &testutils.MockRunner{
			Err: testCase.expectedError,
		}

		output := &bytes.Buffer{}

		config := steps.NewConfig("", "", "", profile.Profile{})
		config.NetworkConfig = steps.NetworkConfig{
			EtcdHost:          testCase.etcdHost,
			EtcdVersion:       testCase.etcdVersion,
			EtcdRepositoryUrl: testCase.etcdRepositoryUrl,

			Arch:            testCase.arch,
			OperatingSystem: testCase.operatingSystem,

			Network:     testCase.network,
			NetworkType: testCase.networkType,
		}
		config.Runner = r
		config.IsMaster = true
		// Mark as done, we assume that etcd has been already deployed

		task := &Step{
			scriptTemplate: tpl,
		}

		err := task.Run(context.Background(), output, config)

		if testCase.expectedError != errors.Cause(err) {
			t.Fatalf("wrong error expected %v actual %v", testCase.expectedError, err)
		}

		if !strings.Contains(output.String(), testCase.etcdRepositoryUrl) {
			t.Fatalf("Etcd repository url Version %s not found in output %s", testCase.etcdRepositoryUrl, output.String())
		}

		if !strings.Contains(output.String(), testCase.etcdVersion) {
			t.Fatalf("Etcd Version %s not found in output %s", testCase.etcdVersion, output.String())
		}

		if !strings.Contains(output.String(), testCase.arch) {
			t.Fatalf("architecture %s not found in output %s", testCase.arch, output.String())
		}

		if !strings.Contains(output.String(), testCase.network) {
			t.Fatalf("network %s not found in output %s", testCase.network, output.String())
		}

		if !strings.Contains(output.String(), testCase.networkType) {
			t.Fatalf("network type %s not found in output %s", testCase.networkType, output.String())
		}

		if !strings.Contains(output.String(), testCase.arch) {
			t.Fatalf("arch %s not found in output %s", testCase.arch, output.String())
		}

		if !strings.Contains(output.String(), testCase.operatingSystem) {
			t.Fatalf("operating system %s not found in output %s", testCase.operatingSystem, output.String())
		}

		if testCase.expectedError == nil && !strings.Contains(output.String(), testCase.etcdHost) {
			t.Fatalf("etcd host %s not found in output %s", testCase.etcdHost, output.String())
		}
	}
}
