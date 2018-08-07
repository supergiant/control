package etcd

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestInstallEtcD(t *testing.T) {
	script := `cat > /etc/systemd/system/etcd.service <<EOF
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd

[Service]
Type=notify
Restart=always
RestartSec={{ .RestartTimeout }}s
LimitNOFILE=40000
TimeoutStartSec={{ .StartTimeout }}s

ExecStart=  docker run \
            -p {{ .ServicePort }}:{{ .ServicePort }} \
            -p {{ .ManagementPort }}:{{ .ManagementPort }} \
            --mount type=bind,source=/tmp/etcd-data.tmp,destination={{ .DataDir }} \
            --name etcd-gcr-v{{ .Version }} \
            gcr.io/etcd-development/etcd:v{{ .Version }} \
            /usr/local/bin/etcd \
            --name {{ .Name }} \
            --data-dir {{ .DataDir }} \
            --listen-client-urls http://{{ .MasterPrivateIP }}:{{ .ServicePort }} \
            --advertise-client-urls http://{{ .MasterPrivateIP }}:{{ .ServicePort }} \
            --listen-peer-urls http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --initial-advertise-peer-urls http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --initial-cluster s1=http://{{ .MasterPrivateIP }}:{{ .ManagementPort }} \
            --initial-cluster-token {{ .ClusterToken }} \
            --initial-cluster-state new

[Install]
WantedBy=multi-user.target
EOF
systemctl start etcd
`
	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing etcd test templatemanager %s", err.Error())
		return
	}

	masterPrivateIp := "10.20.30.40"
	servicePort := "2379"
	managementPort := "2380"
	dataDir := "/var/data"
	version := "3.3.9"
	name := "etcd0"
	clusterToken := "tkn"

	r := &testutils.FakeRunner{
		Err: nil,
	}

	output := &bytes.Buffer{}
	config := steps.Config{
		EtcdConfig: steps.EtcdConfig{
			MasterPrivateIP: masterPrivateIp,
			ServicePort:     servicePort,
			ManagementPort:  managementPort,
			DataDir:         dataDir,
			Version:         version,
			Name:            name,
			ClusterToken:    clusterToken,
			RestartTimeout:  "5",
			StartTimeout:    "0",
		},
		Runner: r,
	}

	task := &Step{
		scriptTemplate: tpl,
	}

	err = task.Run(context.Background(), output, &config)

	if err != nil {
		t.Errorf("Unpexpected error %s", err.Error())
	}

	if !strings.Contains(output.String(), masterPrivateIp) {
		t.Errorf("Master private ip %s not found in %s", masterPrivateIp, output.String())
	}

	if !strings.Contains(output.String(), servicePort) {
		t.Errorf("Service port %s not found in %s", servicePort, output.String())
	}

	if !strings.Contains(output.String(), managementPort) {
		t.Errorf("Management port %s not found in %s", managementPort, output.String())
	}

	if !strings.Contains(output.String(), dataDir) {
		t.Errorf("data dir %s not found in %s", dataDir, output.String())
	}

	if !strings.Contains(output.String(), version) {
		t.Errorf("version %s not found in %s", version, output.String())
	}

	if !strings.Contains(output.String(), name) {
		t.Errorf("name %s not found in %s", name, output.String())
	}

	if !strings.Contains(output.String(), clusterToken) {
		t.Errorf("cluster token %s not found in %s", clusterToken, output.String())
	}
}
