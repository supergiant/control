package flannel

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
	"text/template"
	"io/ioutil"
)

type fakeRunner struct {
	errMsg string
}

func (f *fakeRunner) Run(command *runner.Command) error {
	if len(f.errMsg) > 0 {
		return errors.New(f.errMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}

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

func TestFlannelErrors(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	task := &Step{
		proxyTemplate,
	}

	cfg := steps.NewConfig("", "", "", profile.Profile{})
	cfg.Runner = r
	err = task.Run(context.Background(), output, cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
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

	if len(s.Depends()) != 1 && s.Depends()[0] != etcd.StepName {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{etcd.StepName})
	}
}


func TestStep_Rollback(t *testing.T) {
	s := Step{}
	err  := s.Rollback(context.Background(), ioutil.Discard, &steps.Config{})

	if err != nil {
		t.Errorf("unexpected error while rollback %v", err)
	}
}