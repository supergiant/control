package helm

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
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

func TestInstallTiller(t *testing.T) {
	helmVersion := "helm-v2.8.2"
	rbacEnabled := false
	operatingSystem := "linux"
	arch := "amd64"
	r := &fakeRunner{}

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl, _ := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)

	j := &Step{
		tpl,
	}

	cfg := &steps.Config{
		Kube: model.Kube{
			HelmVersion:     helmVersion,
			RBACEnabled:     rbacEnabled,
			OperatingSystem: operatingSystem,
			Arch:            arch,
		},
		Runner: r,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node: %v", err)
	}

	if !strings.Contains(output.String(), helmVersion) {
		t.Errorf("helm version %s not found in %s", helmVersion, output.String())
	}

	if !strings.Contains(output.String(), helmVersion) {
		t.Errorf("OS %s not found in %s", operatingSystem, output.String())
	}

	if !strings.Contains(output.String(), helmVersion) {
		t.Errorf("arch %s not found in %s", arch, output.String())
	}
}

func TestInstallTillerError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	j := &Step{
		proxyTemplate,
	}

	cfg := &steps.Config{
		Runner: r,
	}
	err = j.Run(context.Background(), output, cfg)

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

	if len(s.Depends()) != 1 && s.Depends()[0] != poststart.StepName {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{poststart.StepName})
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
