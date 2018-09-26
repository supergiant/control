package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"

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

func TestInstallDocker(t *testing.T) {
	dockerVersion := "17.05"
	releaseVersion := "1.29"
	arch := "amd64"
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := &bytes.Buffer{}
	r := &testutils.MockRunner{}

	config := steps.Config{
		DockerConfig: steps.DockerConfig{
			Version:        dockerVersion,
			ReleaseVersion: releaseVersion,
			Arch:           arch,
		},
		Runner: r,
	}

	task := &Step{
		scriptTemplate: tpl,
	}

	err = task.Run(context.Background(), output, &config)

	if !strings.Contains(output.String(), dockerVersion) {
		t.Fatalf("docker version %s not found in output %s", dockerVersion, output.String())
	}

	if !strings.Contains(output.String(), dockerVersion) {
		t.Fatalf("docker version %s not found in output %s", dockerVersion, output.String())
	}
}

func TestDockerError(t *testing.T) {
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

	if len(s.Depends()) != 0 {
		t.Errorf("Wrong dependency list %v expected %v", s.Depends(), []string{})
	}
}
