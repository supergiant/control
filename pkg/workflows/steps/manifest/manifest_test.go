package manifest

import (
	"bytes"

	"context"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/templatemanager"
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

func TestWriteManifestMaster(t *testing.T) {
	var (
		kubernetesVersion   = "1.8.7"
		kubernetesConfigDir = "/kubernetes/conf/dir"
		providerString      = "aws"
		masterHost          = "127.0.0.1"
		masterPort          = "8080"

		r runner.Runner = &fakeRunner{}
	)
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	cfg := steps.Config{
		IsMaster:            true,
		ManifestConfig: steps.ManifestConfig{
			K8SVersion:          kubernetesVersion,
			KubernetesConfigDir: kubernetesConfigDir,
			RBACEnabled:         true,
			MasterHost:          masterHost,
			MasterPort:          masterPort,
			ProviderString:      providerString,
		},
		MasterNodes: map[string]*node.Node{
			"id": {
				PrivateIp: masterHost,
			},
		},
		Runner: r,
	}

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, &cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), kubernetesVersion) {
		t.Errorf("kubernetes version dir %s not found in %s", kubernetesVersion, output.String())
	}

	if !strings.Contains(output.String(), "NodeRestriction") {
		t.Errorf("NodeRestriction not found in %s", output.String())
	}

	if !strings.Contains(output.String(), masterHost) {
		t.Errorf("master host %s not found in %s", masterHost, output.String())
	}

	if !strings.Contains(output.String(), masterPort) {
		t.Errorf("master port %s not found in %s", masterPort, output.String())
	}

	if !strings.Contains(output.String(), providerString) {
		t.Errorf("provider string %s not found in %s", providerString, output.String())
	}
}

func TestWriteManifestNode(t *testing.T) {
	var (
		kubernetesVersion   = "1.8.7"
		kubernetesConfigDir = "/kubernetes/conf/dir"
		providerString      = "aws"
		masterHost          = "127.0.0.1"
		masterPort          = "8080"

		r runner.Runner = &fakeRunner{}
	)

	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)
	cfg := steps.Config{
		MasterNodes: map[string]*node.Node{
			"id": {
				PrivateIp: masterHost,
			},
		},
		ManifestConfig: steps.ManifestConfig{
			IsMaster:            false,
			K8SVersion:          kubernetesVersion,
			KubernetesConfigDir: kubernetesConfigDir,
			RBACEnabled:         false,
			MasterPort:          masterPort,
			ProviderString:      providerString,
		},
		Runner: r,
	}

	j := &Step{
		tpl,
	}

	err = j.Run(context.Background(), output, &cfg)

	if err != nil {
		t.Errorf("Unpexpected error while  provision node %v", err)
	}

	if !strings.Contains(output.String(), kubernetesConfigDir) {
		t.Errorf("kubernetes config dir %s not found in %s", kubernetesConfigDir, output.String())
	}

	if !strings.Contains(output.String(), kubernetesVersion) {
		t.Errorf("kubernetes version dir %s not found in %s", kubernetesVersion, output.String())
	}

	if !strings.Contains(output.String(), masterHost) {
		t.Errorf("master host %s not found in %s", masterHost, output.String())
	}

	if !strings.Contains(output.String(), masterPort) {
		t.Errorf("master port %s not found in %s", masterPort, output.String())
	}

	if strings.Contains(output.String(), "kube-apiserver.yaml") {
		t.Errorf("Unexpected section kube-apiserver.yaml in node manifest %s", output.String())
	}
}

func TestWriteManifestError(t *testing.T) {
	errMsg := "error has occurred"

	r := &fakeRunner{
		errMsg: errMsg,
	}

	proxyTemplate, err := template.New(StepName).Parse("")
	output := new(bytes.Buffer)
	cfg := steps.Config{
		ManifestConfig: steps.ManifestConfig{},
		MasterNodes: map[string]*node.Node{
			"id": {
				PrivateIp: "127.0.0.1",
			},
		},
		Runner: r,
	}

	j := &Step{
		proxyTemplate,
	}

	err = j.Run(context.Background(), output, &cfg)

	if err == nil {
		t.Errorf("Error must not be nil")
		return
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Error message expected to contain %s actual %s", errMsg, err.Error())
	}
}
