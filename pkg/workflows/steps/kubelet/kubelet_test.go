package kubelet

import (
	"bytes"

	"context"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

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

func TestStartKubelet(t *testing.T) {
	k8sVersion := "1.8.7"
	etcdPort := "2379"
	proxyPort := "8080"

	var (
		r             runner.Runner = &fakeRunner{}
		kubeletScript               = `echo 'gcr.io/google-containers/hyperkube:v{{ .K8SVersion }}' > /etc/systemd/system/kubelet.service;systemctl start kubelet`
	)

	kubeletScriptTemplate, err := template.New(StepName).Parse(kubeletScript)

	if err != nil {
		t.Errorf("Error while parsing kubelet script templatemanager %v", err)
	}

	output := new(bytes.Buffer)

	cfg := &steps.Config{
		KubeletConfig: steps.KubeletConfig{
			K8SVersion:     k8sVersion,
			ProxyPort:      proxyPort,
			EtcdClientPort: etcdPort,
		},
		Runner: r,
	}

	task := &Step{
		kubeletScriptTemplate,
	}

	err = task.Run(context.Background(), output, cfg)

	if !strings.Contains(output.String(), k8sVersion) {
		t.Errorf("k8s version %s not found in %s", k8sVersion, output.String())
	}
}

func TestStartKubeletError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	kubeletScriptTemplate, err := template.New(StepName).Parse("")

	output := new(bytes.Buffer)
	config := &steps.Config{
		KubeletConfig: steps.KubeletConfig{},
		Runner:        r,
	}

	j := &Step{
		kubeletScriptTemplate,
	}

	err = j.Run(context.Background(), output, config)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
