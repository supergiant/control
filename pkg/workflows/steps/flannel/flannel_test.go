package flannel

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestFlannelJob_InstallFlannel(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	etcdHost := "127.0.0.1"

	testCases := []struct {
		version       string
		arch          string
		network       string
		networkType   string
		expectedError error
	}{
		{
			"0.9.0",
			"amd64",
			"10.0.2.0/24",
			"vxlan",
			nil,
		},
		{
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

		config := &steps.Config{
			FlannelConfig: steps.FlannelConfig{
				testCase.version,
				testCase.arch,
				etcdHost,
			},
			Runner: r,
		}

		task := &Step{
			scriptTemplate: tpl,
		}

		err := task.Run(context.Background(), output, config)

		if testCase.expectedError != errors.Cause(err) {
			t.Fatalf("wrong error expected %v actual %v", testCase.expectedError, err)
		}

		if !strings.Contains(output.String(), testCase.version) {
			t.Fatalf("Version %s not found in output %s", testCase.version, output.String())
		}

		if !strings.Contains(output.String(), testCase.arch) {
			t.Fatalf("architecture %s not found in output %s", testCase.arch, output.String())
		}

		if testCase.expectedError == nil && !strings.Contains(output.String(), etcdHost) {
			t.Fatalf("etcd host %s not found in output %s", etcdHost, output.String())
		}
	}
}
