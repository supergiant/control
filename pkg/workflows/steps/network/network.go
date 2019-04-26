package network

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "network"

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
	logrus.Debugf("cluster %s: network config: %+v", config.ClusterName, config.NetworkConfig)
	config.NetworkConfig.IsBootstrap = config.IsBootstrap

	err := steps.RunTemplate(context.Background(), t.script, config.Runner, out, config.NetworkConfig)
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
