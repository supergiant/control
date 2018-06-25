package flannel

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
	"github.com/supergiant/supergiant/pkg/jobs"
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
		expectedError string
	}{
		{
			"0.9.0",
			"amd64",
			"10.0.2.0/24",
			"vxlan",
			"",
		},
		{
			"",
			"",
			"",
			"",
			"error has occurred",
		},
	}

	for _, testCase := range testCases {
		r := &jobs.FakeRunner{
			ErrMsg: testCase.expectedError,
		}

		buffer := &bytes.Buffer{}

		job := &Job{
			scriptTemplate: tpl,
			runner:         r,
			output:         buffer,
		}

		config := JobConfig{
			testCase.version,
			testCase.arch,
			testCase.network,
			testCase.networkType,
		}

		err := job.InstallFlannel(config)

		if len(testCase.expectedError) > 0 {
			if err == nil {
				t.Error("error must not be nil")
			}

			if !strings.Contains(err.Error(), testCase.expectedError) {
				t.Errorf("error message %s must contain substring %s", err.Error(), testCase.expectedError)
			}

			return
		}

		if !strings.Contains(buffer.String(), testCase.version) {
			t.Errorf("Version %s not found in output %s", testCase.version, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.arch) {
			t.Errorf("architecture %s not found in output %s", testCase.arch, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.network) {
			t.Errorf("network %s not found in output %s", testCase.network, buffer.String())
		}

		if !strings.Contains(buffer.String(), testCase.networkType) {
			t.Errorf("network type %s not found in output %s", testCase.networkType, buffer.String())
		}
	}
}
