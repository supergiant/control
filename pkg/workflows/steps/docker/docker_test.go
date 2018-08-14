package docker

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestInstallDocker(t *testing.T) {
	script := `
#!/bin/sh

# https://download.docker.com/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_17.06.0~ce-0~ubuntu_amd64.deb

DOCKER_VERSION={{ .DockerVersion }}
UBUNTU_RELEASE={{ .ReleaseVersion }}
ARCH={{ .Arch }}
OUT_DIR=/tmp
URL="https://download.docker.com/linux/ubuntu/dists/${UBUNTU_RELEASE}/pool/stable/${ARCH}/docker-ce_${DOCKER_VERSION}~ce-0~ubuntu_${ARCH}.deb"

wget -O $OUT_DIR/$(basename $URL) $URL
sudo apt install -y $OUT_DIR/$(basename $URL)
rm $OUT_DIR/$(basename $URL)
`
	dockerVersion := "17.05"
	releaseVersion := "1.29"
	arch := "amd64"

	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing flannel test templatemanager %s", err.Error())
		return
	}
	output := &bytes.Buffer{}
	r := &testutils.FakeRunner{}

	config := steps.Config{
		DockerConfig: steps.DockerConfig{
			DockerVersion:  dockerVersion,
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
