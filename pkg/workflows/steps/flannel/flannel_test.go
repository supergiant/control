package flannel

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func TestFlannelJob_InstallFlannel(t *testing.T) {
	script := `#!/bin/bash
wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld
chmod 755 /usr/bin/flanneld

# install etcdctl
GITHUB_URL=https://github.com/coreos/etcd/releases/download
ETCD_VER=v3.3.9
curl -L ${GITHUB_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /usr/bin --strip-components=1

ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_INITIAL_ADVERTISE_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCD_ADVERTISE_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_CLIENT_URLS=http://{{ .EtcdHost }}:2379
ETCD_LISTEN_PEER_URLS=http://{{ .EtcdHost }}:2380
ETCDCTL_API=3 /usr/bin/etcdctl version

/usr/bin/etcdctl set /coreos.com/network/config '{"Network":"{{ .Network }}", "Backend": {"Type": "{{ .NetworkType }}"}}'
/usr/bin/etcdctl get /coreos.com/network/config

cat << EOF > /etc/systemd/system/flanneld.service
[Unit]
Description=Networking service

[Service]
Restart=always

Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
Environment="ETCDCTL_API=3"
ExecStart=/usr/bin/flanneld --etcd-endpoints=http://{{ .EtcdHost }}:2379

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable flanneld.service
systemctl start flanneld.service
`

	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing flannel test templatemanager %s", err.Error())
		return
	}

	etcdHost := "127.0.0.1"

	testCases := []struct {
		version       string
		arch          string
		network       string
		networkType   string
		expectedError error
	}{
		{
			"0.9.0",
			"amd64",
			"10.0.2.0/24",
			"vxlan",
			nil,
		},
		{
			"",
			"",
			"",
			"",
			errors.New("error has occurred"),
		},
	}

	for _, testCase := range testCases {
		r := &testutils.MockRunner{
			Err: testCase.expectedError,
		}

		output := &bytes.Buffer{}

		config := &steps.Config{
			FlannelConfig: steps.FlannelConfig{
				testCase.version,
				testCase.arch,
				testCase.network,
				etcdHost,
				testCase.networkType,
			},
			Runner: r,
		}

		task := &Step{
			scriptTemplate: tpl,
		}

		err := task.Run(context.Background(), output, config)

		if testCase.expectedError != errors.Cause(err) {
			t.Fatalf("wrong error expected %v actual %v", testCase.expectedError, err)
		}

		if !strings.Contains(output.String(), testCase.version) {
			t.Fatalf("Version %s not found in output %s", testCase.version, output.String())
		}

		if !strings.Contains(output.String(), testCase.arch) {
			t.Fatalf("architecture %s not found in output %s", testCase.arch, output.String())
		}

		if !strings.Contains(output.String(), testCase.network) {
			t.Fatalf("network %s not found in output %s", testCase.network, output.String())
		}

		if !strings.Contains(output.String(), testCase.networkType) {
			t.Fatalf("network type %s not found in output %s", testCase.networkType, output.String())
		}

		if testCase.expectedError == nil && !strings.Contains(output.String(), etcdHost) {
			t.Fatalf("etcd host %s not found in output %s", etcdHost, output.String())
		}
	}
}
