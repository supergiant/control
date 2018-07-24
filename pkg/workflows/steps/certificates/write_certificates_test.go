package tiller

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

func TestInstallTiller(t *testing.T) {
	var (
		kubernetesConfigDir   = "/etc/kubernetes"
		CACert                = "abc"
		CACertName            = "xyz"
		CAKeyCert             = "123"
		CAKeyName             = "456"
		APIServerCert         = "qwe"
		APIServerCertName     = "asd"
		APIServerKey          = "zxc"
		APIServerKeyName      = "tyu"
		kubeletClientCert     = "678"
		kubeletClientCertName = "234"
		kubeletClientKey      = "jkl"
		kubeletClientKeyName  = "iop"

		r            runner.Runner = &fakeRunner{}
		tillerScript               = `KUBERNETES_SSL_DIR={{ .KubernetesConfigDir }}/ssl

mkdir -p ${KUBERNETES_SSL_DIR}
echo "{{ .CACert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CACertName }}'
echo "{{ .CAKeyCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .CAKeyName }}'
echo "{{ .APIServerCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerCertName }}'
echo "{{ .APIServerKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .APIServerKeyName }}'
echo "{{ .KubeletClientCert }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientCertName }}'
echo "{{ .KubeletClientKey }}" > ${KUBERNETES_SSL_DIR}/'{{ .KubeletClientKeyName }}'`
	)

	proxyTemplate, err := template.New("tiller").Parse(tillerScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)

	cfg := steps.Config{
		CertificatesConfig: steps.CertificatesConfig{
			kubernetesConfigDir,
			CACert,
			CACertName,
			CAKeyCert,
			CAKeyName,
			APIServerCert,
			APIServerCertName,
			APIServerKey,
			APIServerKeyName,
			kubeletClientCert,
			kubeletClientCertName,
			kubeletClientKey,
			kubeletClientKeyName,
		},
	}

	task := &Task{
		r,
		proxyTemplate,
		output,
	}

	err = task.Run(context.Background(), cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), CAKeyName) {
		t.Errorf("CA key name %s not found in %s", CAKeyName, output.String())
	}

	if !strings.Contains(output.String(), CAKeyCert) {
		t.Errorf("CA key cert %s not found in %s", CAKeyCert, output.String())
	}

	if !strings.Contains(output.String(), CACertName) {
		t.Errorf("CA cert name %s not found in %s", CACertName, output.String())
	}

	if !strings.Contains(output.String(), APIServerKeyName) {
		t.Errorf("API server key name %s not found in %s", APIServerKeyName, output.String())
	}

	if !strings.Contains(output.String(), APIServerKey) {
		t.Errorf("API server key name %s not found in %s", APIServerKey, output.String())
	}

	if !strings.Contains(output.String(), APIServerKey) {
		t.Errorf("API server key name %s not found in %s", APIServerKey, output.String())
	}

	if !strings.Contains(output.String(), kubeletClientKeyName) {
		t.Errorf("kube client key name %s not found in %s", kubeletClientKeyName, output.String())
	}

	if !strings.Contains(output.String(), kubeletClientKey) {
		t.Errorf("kube client key %s not found in %s", kubeletClientKey, output.String())
	}

	if !strings.Contains(output.String(), kubeletClientCertName) {
		t.Errorf("kubelet client cert name %s not found in %s", kubeletClientCertName, output.String())
	}

	if !strings.Contains(output.String(), kubeletClientCert) {
		t.Errorf("kubelet client cert %s not found in %s", kubeletClientCert, output.String())
	}
}

func TestInstallTillerError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New("tiller").Parse("")
	output := new(bytes.Buffer)

	task := &Task{
		r,
		proxyTemplate,
		output,
	}

	cfg := steps.Config{
		CertificatesConfig: steps.CertificatesConfig{},
	}
	err = task.Run(context.Background(), cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
