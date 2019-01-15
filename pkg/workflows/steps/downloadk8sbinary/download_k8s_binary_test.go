package downloadk8sbinary

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/testutils"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func TestDownloadK8SBinary(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl, _ := templatemanager.GetTemplate(StepName)

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
			script: tpl,
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

func TestStepName(t *testing.T) {
	s := Step{}

	if s.Name() != StepName {
		t.Errorf("Unexpected step name expected %s actual %s", StepName, s.Name())
	}
}

func TestDepends(t *testing.T) {
	s := Step{}

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}

func TestStep_Rollback(t *testing.T) {
	s := Step{}
	err := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}

func TestNew(t *testing.T) {
	tpl := template.New("test")
	s := New(tpl)

	if s.script != tpl {
		t.Errorf("Wrong template expected %v actual %v", tpl, s.script)
	}
}

func TestInit(t *testing.T) {
	templatemanager.SetTemplate(StepName, &template.Template{})
	Init()
	templatemanager.DeleteTemplate(StepName)

	s := steps.GetStep(StepName)

	if s == nil {
		t.Error("Step not found")
	}
}

func TestInitPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("recover output must not be nil")
		}
	}()

	Init()

	s := steps.GetStep("not_found.sh.tpl")

	if s == nil {
		t.Error("Step not found")
	}
}

func TestStep_Description(t *testing.T) {
	s := &Step{}

	if desc := s.Description(); desc != "Download kubectl" {
		t.Errorf("Wrong desription expected %s actual %s",
			"Download kubectl", desc)
	}
}
