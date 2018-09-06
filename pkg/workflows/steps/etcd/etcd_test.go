package etcd

import (
	"bytes"
	"context"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"strings"
	"testing"
	"time"
)

func TestInstallEtcD(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
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
	config := steps.NewConfig("", "", "", profile.Profile{})
	config.EtcdConfig = steps.EtcdConfig{
		Host:           host,
		ServicePort:    servicePort,
		ManagementPort: managementPort,
		DataDir:        dataDir,
		Version:        version,
		Name:           name,
		DiscoveryUrl:   clusterToken,
		Timeout:        time.Second * 10,
		RestartTimeout: "5",
		StartTimeout:   "0",
	}
	config.IsMaster = true
	config.Runner = r

	task := &Step{
		scriptTemplate: tpl,
	}

	err = task.Run(context.Background(), output, config)

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
}

func TestInstallEtcdTimeout(t *testing.T) {
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	r := &testutils.MockRunner{
		Err: nil,
	}

	output := &bytes.Buffer{}
	config := steps.NewConfig("", "", "", profile.Profile{})
	config.EtcdConfig = steps.EtcdConfig{
		Timeout: time.Second * 0,
	}
	config.IsMaster = true
	config.Runner = r

	task := &Step{
		scriptTemplate: tpl,
	}

	err = task.Run(context.Background(), output, config)

	if err == nil {
		t.Error("Error must not be nil")
	}

	if !strings.Contains(err.Error(), "deadline") {
		t.Errorf("deadline not found in error message %s", err.Error())
	}
}
