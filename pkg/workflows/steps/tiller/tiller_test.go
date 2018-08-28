package tiller

import (
	"bytes"

	"context"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

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

func TestInstallTiller(t *testing.T) {
	helmVersion := "helm-v2.8.2"
	operatingSystem := "linux"
	arch := "amd64"
	r := &fakeRunner{}

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)

	j := &Step{
		tpl,
	}

	cfg := &steps.Config{
		TillerConfig: steps.TillerConfig{
			helmVersion,
			operatingSystem,
			arch,
		},
		Runner: r,
	}

	err = j.Run(context.Background(), output, cfg)

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

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	j := &Step{
		proxyTemplate,
	}

	cfg := &steps.Config{
		TillerConfig: steps.TillerConfig{},
		Runner:       r,
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
