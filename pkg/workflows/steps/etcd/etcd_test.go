package etcd

import (
	"bytes"
	"context"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/testutils"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"strings"
	"testing"
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
	clusterSizes := []int{1, 3}

	r := &testutils.MockRunner{
		Err: nil,
	}

	for _, clusterSize := range clusterSizes {
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
				ClusterSize:    clusterSize,
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

		if clusterSize > 1 {
			if !strings.Contains(output.String(), "discovery") {
				t.Errorf("discovery url not found")
			}
		} else {
			if !strings.Contains(output.String(), "--initial-cluster") {
				t.Error("initial-cluster not found")
			}
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
}
