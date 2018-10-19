package manifest

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/certificates"
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

func TestWriteManifestMaster(t *testing.T) {
	var (
		kubernetesVersion   = "1.8.7"
		kubernetesConfigDir = "/kubernetes/conf/dir"
		providerString      = "aws"
		masterHost          = "10.20.30.40"
		masterPort          = "8080"

		r runner.Runner = &fakeRunner{}
	)
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	p := profile.Profile{
		K8SVersion:  kubernetesVersion,
		RBACEnabled: true,
	}
	cfg := steps.NewConfig("", "", "", p)
	cfg.Node = node.Node{
		PrivateIp: masterHost,
	}

	cfg.IsMaster = true
	cfg.ManifestConfig.MasterPort = masterPort
	cfg.ManifestConfig.ProviderString = providerString
	cfg.ManifestConfig.KubernetesConfigDir = kubernetesConfigDir
	cfg.Runner = r

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), kubernetesVersion) {
		t.Errorf("kubernetes version dir %s not found in %s", kubernetesVersion, output.String())
	}

	if !strings.Contains(output.String(), "NodeRestriction") {
		t.Errorf("NodeRestriction not found in %s", output.String())
	}

	if !strings.Contains(output.String(), masterHost) {
		t.Errorf("master host %s not found in %s", masterHost, output.String())
	}

	if !strings.Contains(output.String(), masterPort) {
		t.Errorf("master port %s not found in %s", masterPort, output.String())
	}

	if !strings.Contains(output.String(), providerString) {
		t.Errorf("provider string %s not found in %s", providerString, output.String())
	}
}

func TestWriteManifestNode(t *testing.T) {
	var (
		kubernetesVersion   = "1.8.7"
		kubernetesConfigDir = "/kubernetes/conf/dir"
		providerString      = "aws"
		masterHost          = "127.0.0.1"

		r runner.Runner = &fakeRunner{}
	)

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	p := profile.Profile{
		K8SVersion:  kubernetesVersion,
		RBACEnabled: false,
	}
	cfg := steps.NewConfig("", "", "", p)
	cfg.AddMaster(&node.Node{
		State:     node.StateActive,
		PrivateIp: masterHost,
	})
	cfg.Runner = r
	cfg.ManifestConfig.KubernetesConfigDir = kubernetesConfigDir
	cfg.ManifestConfig.ProviderString = providerString
	cfg.ManifestConfig.IsMaster = false

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), kubernetesVersion) {
		t.Errorf("kubernetes version dir %s not found in %s", kubernetesVersion, output.String())
	}

	if !strings.Contains(output.String(), masterHost) {
		t.Errorf("master host %s not found in %s", masterHost, output.String())
	}

	if strings.Contains(output.String(), "kube-apiserver.yaml") {
		t.Errorf("Unexpected section kube-apiserver.yaml in node manifest %s", output.String())
	}
}

func TestWriteManifestError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)
	cfg := steps.NewConfig("", "", "", profile.Profile{})
	cfg.AddMaster(&node.Node{
		State:     node.StateActive,
		PrivateIp: "127.0.0.1",
	})
	cfg.Runner = r

	j := &Step{
		proxyTemplate,
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

	if len(s.Depends()) != 1 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{certificates.StepName})
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
	Init()

	s := steps.GetStep(StepName)

	if s == nil {
		t.Error("Step not found")
	}
}
