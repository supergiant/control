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

	tpl, err := template.New(StepName).Parse(script)

	if err != nil {
		t.Errorf("error parsing flannel test templatemanager %s", err.Error())
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

		output := &bytes.Buffer{}

		config := steps.Config{
			FlannelConfig: steps.FlannelConfig{
				testCase.version,
				testCase.arch,
				testCase.network,
				testCase.networkType,
			},
			Runner: r,
		}

		task := &Task{
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
	}
}
