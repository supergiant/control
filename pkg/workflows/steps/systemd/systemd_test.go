package tiller

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
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

func TestSystemdUpdate(t *testing.T) {
	var (
		r              runner.Runner = &fakeRunner{}
		k8sVersion                   = "1.8.7"
		k8sProvider                  = "provider"
		kubeletService               = "service"

		postStartScript = `    cat << EOF > /etc/systemd/system/{{ .KubeletService }}
[Unit]
Description=Kubernetes Kubelet Server
Documentation=https://github.com/kubernetes/kubernetes
Requires=docker.service network-online.target
After=docker.service network-online.target

[Service]
ExecStartPre=/bin/bash -c "/opt/bin/download-k8s-binary"
ExecStartPost=/bin/bash -c "/opt/bin/kube-post-start.sh"

ExecStart=/usr/bin/docker run \
        --net=host \
        --pid=host \
        --privileged \
        -v /dev:/dev \
        -v /sys:/sys:ro \
        -v /var/run:/var/run:rw \
        -v /var/lib/docker/:/var/lib/docker:rw \
        -v /var/lib/kubelet/:/var/lib/kubelet:shared \
        -v /var/log:/var/log:shared \
        -v /srv/kubernetes:/srv/kubernetes:ro \
        -v /etc/kubernetes:/etc/kubernetes:ro \
        gcr.io/google-containers/hyperkube:v{{ .KubernetesVersion }} \
        /hyperkube kubelet --allow-privileged=true \
        --cluster-dns=10.3.0.10 \
        --cluster_domain=cluster.local \
        --cadvisor-port=0 \
        --pod-manifest-path=/etc/kubernetes/manifests \
        --kubeconfig=/var/lib/kubelet/kubeconfig \
        --volume-plugin-dir=/etc/kubernetes/volumeplugins \
        {{- .KubernetesProvider }}
        --register-node=false
Restart=always
StartLimitInterval=0
RestartSec=10
KillMode=process

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ${KUBELET_SERVICE}
systemctl start ${KUBELET_SERVICE}`
	)

	proxyTemplate, err := template.New("systemd").Parse(postStartScript)

	if err != nil {
		t.Errorf("Error while parsing kubeproxy template %v", err)
	}

	output := new(bytes.Buffer)

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	cfg := Config{
		k8sVersion,
		kubeletService,
		k8sProvider,
	}

	err = j.UpdateSystemd(cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), k8sVersion) {
		t.Errorf("k8s version %s not found in %s", k8sVersion, output.String())
	}

	if !strings.Contains(output.String(), k8sProvider) {
		t.Errorf("k8s provider %s not found in %s", k8sProvider, output.String())
	}

	if !strings.Contains(output.String(), kubeletService) {
		t.Errorf("kubelet service %s not found in %s", kubeletService, output.String())
	}
}

func TestSystemdUpdateError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New("systemd").Parse("")
	output := new(bytes.Buffer)

	j := &Task{
		r,
		proxyTemplate,
		output,
	}

	cfg := Config{}
	err = j.UpdateSystemd(cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
