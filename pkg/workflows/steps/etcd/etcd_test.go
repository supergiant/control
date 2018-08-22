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
	script := `mkdir -p /tmp/etcd-data
cat > /etc/systemd/system/etcd.service <<EOF
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd

[Service]
Restart=always
RestartSec={{ .RestartTimeout }}s
LimitNOFILE=40000
TimeoutStartSec={{ .StartTimeout }}s

ExecStart=/usr/bin/docker run \
            -p {{ .ServicePort }}:{{ .ServicePort }} \
            -p {{ .ManagementPort }}:{{ .ManagementPort }} \
            --volume={{ .DataDir }}:/etcd-data \
            --name {{ .Name }} \
            gcr.io/etcd-development/etcd:v{{ .Version }} \
            /usr/local/bin/etcd \
            --name {{ .Name }} \
            --data-dir /etcd-data \
            --listen-client-urls http://{{ .Host }}:{{ .ServicePort }} \
            --advertise-client-urls http://{{ .Host }}:{{ .ServicePort }} \
            --listen-peer-urls http://{{ .Host }}:{{ .ManagementPort }} \
            --initial-advertise-peer-urls http://{{ .Host }}:{{ .ManagementPort }} \
            --discovery {{ .DiscoveryUrl }} \
            --listen-peer-urls http://{{ .Host }}:2380 --listen-client-urls http://{{ .Host }}:2379

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable etcd.service
systemctl start etcd.service
`
	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing etcd test templatemanager %s", err.Error())
		return
	}

	host := "10.20.30.40"
	servicePort := "2379"
	managementPort := "2380"
	dataDir := "/var/data"
	version := "3.3.9"
	name := "etcd0"
	clusterToken := "tkn"

	r := &testutils.MockRunner{
		Err: nil,
	}

	output := &bytes.Buffer{}
	config := steps.Config{
		EtcdConfig: steps.EtcdConfig{
			Host:           host,
			ServicePort:    servicePort,
			ManagementPort: managementPort,
			DataDir:        dataDir,
			Version:        version,
			Name:           name,
			DiscoveryUrl:   clusterToken,
			RestartTimeout: "5",
			StartTimeout:   "0",
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

	if !strings.Contains(output.String(), host) {
		t.Errorf("Master private ip %s not found in %s", host, output.String())
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
}
