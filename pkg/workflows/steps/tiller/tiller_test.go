package tiller

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
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
	operatingSystem := "linux"
	arch := "amd64"

	var (
		r            runner.Runner = &fakeRunner{}
		tillerScript               = `wget http://storage.googleapis.com/kubernetes-helm/{{ .HelmVersion }}-{{ .OperatingSystem }}-{{ .Arch }}.tar.gz --directory-prefix=/tmp/`
	)

	proxyTemplate, err := template.New("tiller").Parse(tillerScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	cfg := Config{
		helmVersion,
		operatingSystem,
		arch,
	}

	err = j.Run(cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
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

	proxyTemplate, err := template.New("tiller").Parse("")
	output := new(bytes.Buffer)

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	cfg := Config{}
	err = j.Run(cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
