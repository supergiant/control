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
		proxyScript               = `# proxy
cat << EOF > ${KUBERNETES_MANIFESTS_DIR}/kube-proxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-proxy
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/hyperkube:v{{ .K8SVersion }}
    command:
    - /hyperkube
    - proxy
    - --v=2
    - --master=http://{{ .MasterPrivateIP }}:{{ .MasterPort }}
    - --proxy-mode=iptables
    securityContext:
      privileged: true
    volumeMounts:
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
`
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
