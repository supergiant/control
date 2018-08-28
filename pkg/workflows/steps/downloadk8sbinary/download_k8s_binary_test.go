package downloadk8sbinary

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
