package kubeadm

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
)

const (
	StepName = "kubeadm"
)

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

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.KubeadmConfig.Provider = string(config.Provider)
	config.KubeadmConfig.IsBootstrap = config.IsBootstrap
	config.KubeadmConfig.IsMaster = config.IsMaster

	// NOTE(stgleb): Kubeadm accepts only ipv4 or ipv6 addresses as advertise address
	if config.IsBootstrap {
		config.KubeadmConfig.AdvertiseAddress = config.Node.PrivateIp

		// TODO(stgleb): Remove that when all providers support Load Balancers
		if config.Provider == clouds.AWS || config.Provider == clouds.DigitalOcean {
			config.KubeadmConfig.ExternalDNSName = config.Node.PublicIp
			config.KubeadmConfig.InternalDNSName = config.Node.PublicIp
		}
	}

	// TODO(stgleb): Remove that when all providers support Load Balancers
	if config.Provider == clouds.AWS || config.Provider == clouds.DigitalOcean {
		config.KubeadmConfig.InternalDNSName = config.InternalDNSName
		config.KubeadmConfig.ExternalDNSName = config.ExternalDNSName
	} else {
		if master := config.GetMaster(); master != nil {
			config.KubeadmConfig.InternalDNSName = master.PrivateIp
			config.KubeadmConfig.ExternalDNSName = master.PublicIp
		}
	}

	err := steps.RunTemplate(ctx, t.script, config.Runner, out, config.KubeadmConfig)

	if err != nil {
		return errors.Wrap(err, "kubeadm step")
	}

	return nil
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return "run kubeadm"
}

func (s *Step) Depends() []string {
	return []string{docker.StepName}
}
