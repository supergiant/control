package downloadk8sbinary

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestFlannelJob_InstallFlannel(t *testing.T) {
	script := `#!/bin/bash
source /etc/environment
mkdir -p /opt/bin
curl -sSL -o /opt/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v{{ .K8SVersion }}/bin/{{ .OperatingSystem }}/{{ .Arch }}/kubectl
chmod +x /opt/bin/$FILE
chmod +x /opt/bin/kubectl
`

	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing flannel test templatemanager %s", err.Error())
		return
	}

	testCases := []struct {
		version         string
		arch            string
		operatingSystem string
		expectedError   error
	}{
		{
			"1.11",
			"amd64",
			"linux",
			nil,
		},
		{
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

		config := steps.Config{
			DownloadK8sBinary: steps.DownloadK8sBinary{
				testCase.version,
				testCase.arch,
				testCase.operatingSystem,
			},
			Runner: r,
		}

		task := &Step{
			scriptTemplate: tpl,
		}

		err := task.Run(context.Background(), output, &config)

		if testCase.expectedError != errors.Cause(err) {
			t.Fatalf("wrong error expected %v actual %v", testCase.expectedError, err)
		}

		if !strings.Contains(output.String(), testCase.version) {
			t.Fatalf("k8sVersion %s not found in output %s", testCase.version, output.String())
		}

		if !strings.Contains(output.String(), testCase.arch) {
			t.Fatalf("architecture %s not found in output %s", testCase.arch, output.String())
		}

		if !strings.Contains(output.String(), testCase.operatingSystem) {
			t.Fatalf("operating system type %s not found in output %s", testCase.operatingSystem, output.String())
		}
	}
}
