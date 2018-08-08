package certificates

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

func TestWriteCertificates(t *testing.T) {
	var (
		kubernetesConfigDir = "/etc/kubernetes"
		masterPrivateIP     = "10.20.30.40"
		userName            = "user"
		password            = "1234"

		r      runner.Runner = &fakeRunner{}
		script               = `KUBERNETES_SSL_DIR={{ .KubernetesConfigDir }}/ssl

mkdir -p ${KUBERNETES_SSL_DIR}
openssl genrsa -out /etc/kubernetes/ssl/ca-key.pem 2048
openssl req -x509 -new -nodes -key /etc/kubernetes/ssl/ca-key.pem -days 10000 -out /etc/kubernetes/ssl/ca.pem -subj "/CN=kube-ca"
sed -e "s/\{{ .MasterPrivateIP }}/curl ipinfo.io/ip/" < /etc/kubernetes/ssl/openssl.cnf.template > /etc/kubernetes/ssl/openssl.cnf.public
sed -e "s/\{{ .MasterPrivateIP }}/{{ .MasterPrivateIP }}/" < /etc/kubernetes/ssl/openssl.cnf.public > /etc/kubernetes/ssl/openssl.cnf
openssl genrsa -out /etc/kubernetes/ssl/apiserver-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/apiserver-key.pem -out /etc/kubernetes/ssl/apiserver.csr -subj "/CN=kube-apiserver" -config /etc/kubernetes/ssl/openssl.cnf
openssl x509 -req -in /etc/kubernetes/ssl/apiserver.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/apiserver.pem -days 365 -extensions v3_req -extfile /etc/kubernetes/ssl/openssl.cnf
openssl genrsa -out /etc/kubernetes/ssl/worker-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/worker-key.pem -out /etc/kubernetes/ssl/worker.csr -subj "/CN=kube-worker"
openssl x509 -req -in /etc/kubernetes/ssl/worker.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/worker.pem -days 365
openssl genrsa -out /etc/kubernetes/ssl/admin-key.pem 2048
openssl req -new -key /etc/kubernetes/ssl/admin-key.pem -out /etc/kubernetes/ssl/admin.csr -subj "/CN=kube-admin"
openssl x509 -req -in /etc/kubernetes/ssl/admin.csr -CA /etc/kubernetes/ssl/ca.pem -CAkey /etc/kubernetes/ssl/ca-key.pem -CAcreateserial -out /etc/kubernetes/ssl/admin.pem -days 365
chmod 600 /etc/kubernetes/ssl/*-key.pem
chown root:root /etc/kubernetes/ssl/*-key.pem

cat > /etc/kubernetes/ssl/basic_auth.csv <<EOF
{{ .Password }},{{ .Username }},admin
EOF

cat > /etc/kubernetes/ssl/known_tokens.csv <<EOF
{{ .Password }},kubelet,kubelet
{{ .Password }},kube_proxy,kube_proxy
{{ .Password }},system:scheduler,system:scheduler
{{ .Password }},system:controller_manager,system:controller_manager
{{ .Password }},system:logging,system:logging
{{ .Password }},system:monitoring,system:monitoring
{{ .Password }},system:dns,system:dns
EOF`
	)

	certificatesTemplate, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy templatemanager %v", err)
	}

	output := new(bytes.Buffer)

	cfg := steps.Config{
		CertificatesConfig: steps.CertificatesConfig{
			kubernetesConfigDir,
			masterPrivateIP,
			userName,
			password,
		},
		Runner: r,
	}

	task := &Step{
		certificatesTemplate,
	}

	err = task.Run(context.Background(), output, &cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), masterPrivateIP) {
		t.Errorf("master private up %s not found in %s", masterPrivateIP, output.String())
	}

	if !strings.Contains(output.String(), userName) {
		t.Errorf("username %s not found in %s", userName, output.String())
	}

	if !strings.Contains(output.String(), password) {
		t.Errorf("password %s not found in %s", password, output.String())
	}
}

func TestWriteCertificatesError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)

	task := &Step{
		proxyTemplate,
	}

	cfg := steps.Config{
		CertificatesConfig: steps.CertificatesConfig{},
		Runner:             r,
	}
	err = task.Run(context.Background(), output, &cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
