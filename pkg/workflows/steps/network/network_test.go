package network

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
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

func TestNetworkConfig(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl,_ := templatemanager.GetTemplate(StepName)

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
			script: tpl,
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

func TestNetworkErrors(t *testing.T) {
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
	defer func(){
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
