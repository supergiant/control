package cloudcontroller

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
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

func TestKubeadm(t *testing.T) {
	r := &fakeRunner{}
	err := templatemanager.Init("../../../../templates")

	if err != nil {
		t.Fatal(err)
	}

	tpl, _ := templatemanager.GetTemplate(StepName)

	if tpl == nil {
		t.Fatal("template not found")
	}

	output := new(bytes.Buffer)

	cfg := &steps.Config{
		IsMaster:    true,
		IsBootstrap: true,
		Kube: model.Kube{
			Provider: clouds.DigitalOcean,
			Networking: model.Networking{
				CIDR: "10.0.0.0/24",
			},
			BootstrapToken:  "1234",
			ExternalDNSName: "external.dns.name",
			InternalDNSName: "internal.dns.name",
		},
		DigitalOceanConfig: steps.DOConfig{
			AccessToken: "accesstoken",
		},
		Runner: r,
		Node: model.Machine{
			PrivateIp: "10.20.30.40",
		},
	}

	task := &Step{
		tpl,
	}

	err = task.Run(context.Background(), output, cfg)
	require.Nil(t, err)

	if !strings.Contains(output.String(), cfg.DigitalOceanConfig.AccessToken) {
		t.Errorf("token %q not found in %s", cfg.DigitalOceanConfig.AccessToken, output.String())
	}
}
