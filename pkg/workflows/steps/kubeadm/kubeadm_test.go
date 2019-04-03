package kubeadm

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
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

func TestKubeadm(t *testing.T) {
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

	cfg := &steps.Config{
		Provider: clouds.AWS,
		IsMaster: true,
		KubeadmConfig: steps.KubeadmConfig{
			IsMaster:    true,
			IsBootstrap: true,
			CIDR:        "10.0.0.0/24",
			Token:       "1234",
		},
		ExternalDNSName: "external.dns.name",
		InternalDNSName: "internal.dns.name",
		Runner:          r,
	}

	task := &Step{
		tpl,
	}

	err = task.Run(context.Background(), output, cfg)
	require.Nil(t, err)

	if !strings.Contains(output.String(), cfg.KubeadmConfig.CIDR) {
		t.Errorf("CIDR %s not found in %s", cfg.KubeadmConfig.CIDR, output.String())
	}

	if !strings.Contains(output.String(), cfg.KubeadmConfig.Token) {
		t.Errorf("Token %s not found in %s", cfg.KubeadmConfig.Token, output.String())
	}

	if !strings.Contains(output.String(), cfg.KubeadmConfig.InternalDNSName) {
		t.Errorf("LoadBalancerHost %s not found in %s", cfg.KubeadmConfig.InternalDNSName, output.String())
	}
}

func TestStartKubeadmError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	kubeletScriptTemplate, err := template.New(StepName).Parse("")

	output := new(bytes.Buffer)
	config := &steps.Config{
		KubeadmConfig: steps.KubeadmConfig{},
		Runner:        r,
	}

	j := &Step{
		kubeletScriptTemplate,
	}

	err = j.Run(context.Background(), output, config)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if errors.Cause(err) != sgerrors.ErrRawError {
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

	if len(s.Depends()) != 1 && s.Depends()[0] != docker.StepName {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{docker.StepName})
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

	if desc := s.Description(); desc != "run kubeadm" {
		t.Errorf("Wrong desription expected %s actual %s",
			"Run kubelet", desc)
	}
}
