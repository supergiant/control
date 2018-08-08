package kubeproxy

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"context"

	"github.com/supergiant/supergiant/pkg/runner"
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

func TestStartKubeProxy(t *testing.T) {
	masterIp := "20.30.40.50"
	k8sVersion := "1.8.7"

	var (
		r           runner.Runner = &fakeRunner{}
		proxyScript               = `#!/bin/bash
mkdir -p  /etc/kubernetes
sudo docker run -d --privileged=true --volume=/etc/ssl/certs:/usr/share/ca-certificates \
    --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro \
    --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl \
    --volume=/etc/kubernetes/config.json:/etc/kubernetes/config.json \
    gcr.io/google_containers/hyperkube:v{{ .K8SVersion }} /hyperkube proxy \
    --master http://{{ .MasterPrivateIP }}	`
	)

	proxyTemplate, err := template.New(StepName).Parse(proxyScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy templatemanager %v", err)
	}

	output := new(bytes.Buffer)

	cfg := &steps.Config{
		KubeProxyConfig: steps.KubeProxyConfig{
			K8SVersion:      k8sVersion,
			MasterPrivateIP: masterIp,
		},
		Runner: r,
	}

	j := &Step{
		proxyTemplate,
	}

	err = j.Run(context.Background(), output, cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), masterIp) {
		t.Errorf("master ip %s not found in %s", masterIp, output.String())
	}
}

func TestStartKubeProxyError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)
	cfg := &steps.Config{
		KubeProxyConfig: steps.KubeProxyConfig{},
		Runner:          r,
	}

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
