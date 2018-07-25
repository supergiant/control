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
	etcdPort := "2379"
	proxyPort := "8080"

	var (
		r           runner.Runner = &fakeRunner{}
		proxyScript               = `        "master": "http://{{ .MasterPrivateIP }}:{{ .ProxyPort }} http://{{ .MasterPrivateIP }}:{{ .EtcdClientPort }}";sudo docker run --privileged=true --volume=/etc/ssl/cer:/usr/share/ca-certificates --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }} /hyperkube proxy --config /etc/kubernetes/config.json --master http://{{ .MasterPrivateIP }}`
	)

	proxyTemplate, err := template.New(StepName).Parse(proxyScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)

	cfg := steps.Config{
		KubeProxyConfig: steps.KubeProxyConfig{
			KubernetesVersion: k8sVersion,
			MasterPrivateIP:   masterIp,
			ProxyPort:         proxyPort,
			EtcdClientPort:    etcdPort,
		},
	}

	j := &Step{
		r,
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
	cfg := steps.Config{
		KubeProxyConfig: steps.KubeProxyConfig{},
	}

	j := &Step{
		r,
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
