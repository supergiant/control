package docker

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

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
