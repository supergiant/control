package flannel

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"context"
	"github.com/supergiant/supergiant/pkg/testutils"
)

func TestFlannelJob_InstallFlannel(t *testing.T) {
	script := `#!/bin/bash
wget -P /usr/bin/ https://github.com/coreos/flannel/releases/download/v{{ .Version }}/flanneld-{{ .Arch }}
mv /usr/bin/flanneld-{{ .Arch }} /usr/bin/flanneld

	echo "[Unit]
	Description=Networking service
	Requires=etcd-member.service
	[Service]
	Environment=FLANNEL_IMAGE_TAG=v{{ .Version }}
	ExecStartPre=/usr/bin/etcdctl set /coreos.com/network/config '{"Network":"{{ .Network }}", "Backend": {"Type": "{{ .NetworkType }}"}}'" > \
/etc/systemd/system/flanneld.service
systemctl enable flanneld.service
systemctl restart flanneld.service
`

	tpl, err := template.New("").Parse(script)

	if err != nil {
		t.Errorf("error parsing flannel test template %s", err.Error())
		return
	}

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
		r := &testutils.FakeRunner{
			Err: testCase.expectedError,
		}

		buffer := &bytes.Buffer{}

		config := Config{
			testCase.version,
			testCase.arch,
			testCase.network,
			testCase.networkType,
		}

		job := &Task{
			scriptTemplate: tpl,
			runner:         r,
			output:         buffer,
			config:         config,
		}

		err := job.Run(context.Background())

		if testCase.expectedError != errors.Cause(err) {
			t.Fatalf("wrong error expected %v actual %v", testCase.expectedError, err)
		}

		if !strings.Contains(buffer.String(), testCase.version) {
			t.Fatalf("Version %s not found in output %s", testCase.version, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.arch) {
			t.Fatalf("architecture %s not found in output %s", testCase.arch, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.network) {
			t.Fatalf("network %s not found in output %s", testCase.network, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.networkType) {
			t.Fatalf("network type %s not found in output %s", testCase.networkType, buffer.String())
		}
	}
}
