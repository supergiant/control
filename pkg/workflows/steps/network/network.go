package network

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
)

const StepName = "network"

type Config struct {
	IsBootstrap     bool
	CIDR            string
	NetworkProvider string
}

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(tpl *template.Template) *Step {
	return &Step{
		script: tpl,
	}
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	if !config.IsBootstrap {
		return nil
	}

	err := steps.RunTemplate(context.Background(), t.script, config.Runner, out, toStepCfg(config))
	if err != nil {
		return errors.Wrap(err, "configure network step")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Description() string {
	return "Configure CNI plugin, that must happen during bootstrap node provisioning for HA cluster"
}

func (s *Step) Depends() []string {
	return []string{kubeadm.StepName}
}

func toStepCfg(c *steps.Config) Config {
	return Config{
		IsBootstrap:     c.IsBootstrap,
		CIDR:            c.Kube.Networking.CIDR,
		NetworkProvider: c.Kube.Networking.Provider,
	}
}
