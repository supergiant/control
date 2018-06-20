package digitalocean

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
	mustHave := []string{"config", "service", "kubeletscript", "proxy"}

	var (
		r              runner.Runner = &fakeRunner{}
		config                       = `config {{ .MasterPrivateIP }} {{ .KubernetesVersion }} `
		kubeletService               = `service {{ .KubernetesVersion }}`
		kubeletScript                = `kubeletscript '{{ .KubeletService }}' > /etc/systemd/system/kubelet.service`
		proxyScript                  = `proxy "{{ .ConfigFile }}" > /etc/kubernetes/config.json`
	)

	cfgTpl, err := template.New("config").Parse(config)

	if err != nil {
		t.Errorf("Error while parsing config template %v", err)
	}

	kubeletServiceTpl, err := template.New("service").Parse(kubeletService)

	if err != nil {
		t.Errorf("Error while parsing config template %v", err)
	}

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

		cfgTpl,
		kubeletServiceTpl,
		kubeletScriptTemplate,
		proxyTemplate,
		output,
		output,
	}

	err = j.ProvisionNode(k8sVersion, masterIp)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	for _, substr := range mustHave {
		if !strings.Contains(output.String(), substr) {
			t.Errorf("required string %s not found in %s", substr, output.String())
		}
	}

	if !strings.Contains(output.String(), masterIp) {
		t.Errorf("master ip %s not found in %s", masterIp, output.String())
	}

	if !strings.Contains(output.String(), k8sVersion) {
		t.Errorf("k8s version %s not found in %s", k8sVersion, output.String())
	}
}
