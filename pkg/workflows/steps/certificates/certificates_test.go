package certificates

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
	"github.com/supergiant/supergiant/pkg/pki"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
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

func TestWriteCertificates(t *testing.T) {
	var (
		kubernetesConfigDir = "/etc/kubernetes"
		masterPrivateIP     = "10.20.30.40"
		userName            = "user"
		password            = "1234"

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

	caPair, err := pki.NewCAPair(nil)

	if err != nil {
		t.Errorf("unexpected error creating PKI bundle %v", err)
	}

	cfg := steps.NewConfig("", "", "", profile.Profile{})
	cfg.CertificatesConfig = steps.CertificatesConfig{
		KubernetesConfigDir: kubernetesConfigDir,
		MasterHost:          masterPrivateIP,
		Username:            userName,
		Password:            password,
		CAKey:               string(caPair.CA.Key),
		CACert:              string(caPair.CA.Cert),
	}
	cfg.Runner = r
	cfg.AddMaster(&node.Node{
		State:     node.StateActive,
		PrivateIp: "10.20.30.40",
	})
	task := &Step{
		tpl,
	}

	err = task.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), userName) {
		t.Errorf("username %s not found in %s", userName, output.String())
	}

	if !strings.Contains(output.String(), password) {
		t.Errorf("password %s not found in %s", password, output.String())
	}

	if !strings.Contains(output.String(), string(caPair.CA.Key)) {
		t.Errorf("CA key not found in %s", output.String())
	}

	if !strings.Contains(output.String(), string(caPair.CA.Cert)) {
		t.Errorf("CA cert not found in %s", output.String())
	}
}

func TestWriteCertificatesError(t *testing.T) {
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
	cfg.AddMaster(&node.Node{
		State:     node.StateActive,
		PrivateIp: "10.20.30.40",
	})
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
