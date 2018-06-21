package node

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/supergiant/supergiant/pkg/runner"
)

type fakeRunner struct {
}

func (f *fakeRunner) Run(command *runner.Command) error {
	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}

func TestJob_ProvisionNode(t *testing.T) {
	masterIp := "20.30.40.50"
	k8sVersion := "1.8.7"
	etcdPort := "2379"
	proxyPort := "8080"

	var (
		r             runner.Runner = &fakeRunner{}
		kubeletScript               = `echo 'gcr.io/google-containers/hyperkube:v{{ .KubernetesVersion }}' > /etc/systemd/system/kubelet.service;systemctl start kubelet`
		proxyScript                 = `        "master": "http://{{ .MasterPrivateIP }}:{{ .ProxyPort }} http://{{ .MasterPrivateIP }}:{{ .EtcdPort }}";sudo docker run --privileged=true --volume=/etc/ssl/cer:/usr/share/ca-certificates --volume=/etc/kubernetes/worker-kubeconfig.yaml:/etc/kubernetes/worker-kubeconfig.yaml:ro --volume=/etc/kubernetes/ssl:/etc/kubernetes/ssl gcr.io/google_containers/hyperkube:v{{ .KubernetesVersion }} /hyperkube proxy --config /etc/kubernetes/config.json --master http://{{ .MasterPrivateIP }}`
	)

	kubeletScriptTemplate, err := template.New("kubelet").Parse(kubeletScript)

	if err != nil {
		t.Errorf("Error while parsing kubelet script template %v", err)
	}

	proxyTemplate, err := template.New("proxy").Parse(proxyScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)

	j := &Job{
		r,
		kubeletScriptTemplate,
		proxyTemplate,
		output,
		output,
	}

	cfg := JobConfig{
		KubernetesVersion: k8sVersion,
		MasterPrivateIP: masterIp,
		ProxyPort: proxyPort,
		EtcdPort: etcdPort,
	}

	err = j.ProvisionNode(cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), masterIp) {
		t.Errorf("master ip %s not found in %s", masterIp, output.String())
	}

	if !strings.Contains(output.String(), k8sVersion) {
		t.Errorf("k8s version %s not found in %s", k8sVersion, output.String())
	}
}
