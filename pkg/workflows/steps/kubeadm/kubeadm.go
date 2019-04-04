package kubeadm

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/sgerrors"
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
		config.KubeadmConfig.ExternalDNSName = config.Node.PublicIp
		config.KubeadmConfig.InternalDNSName = config.Node.PublicIp
	} else {

		if !config.IsMaster {
			master := config.GetMaster()
			if master == nil {
				return errors.Wrapf(sgerrors.ErrRawError, "no masters in the %s cluster", config.ClusterID)
			}
			config.KubeadmConfig.MasterPrivateIP = master.PrivateIp
			config.KubeadmConfig.InternalDNSName = master.PrivateIp
			config.KubeadmConfig.ExternalDNSName = master.PublicIp
		} else {
			config.KubeadmConfig.MasterPrivateIP = config.Node.PrivateIp
			config.KubeadmConfig.InternalDNSName = config.Node.PrivateIp
			config.KubeadmConfig.ExternalDNSName = config.Node.PublicIp
		}
	}

	// TODO(stgleb): Remove that when all providers support Load Balancers
	if config.Provider == clouds.AWS || config.Provider == clouds.DigitalOcean || config.Provider == clouds.GCE {
		config.KubeadmConfig.InternalDNSName = config.InternalDNSName
		config.KubeadmConfig.ExternalDNSName = config.ExternalDNSName
	}

	// TODO: needs more validation
	switch {
	case config.KubeadmConfig.ExternalDNSName == "":
		return errors.Wrap(sgerrors.ErrRawError, "external dns name should be set")
	case config.KubeadmConfig.InternalDNSName == "":
		return errors.Wrap(sgerrors.ErrRawError, "internal dns name should be set")
	case !config.KubeadmConfig.IsBootstrap && config.KubeadmConfig.MasterPrivateIP == "":
		return errors.Wrap(sgerrors.ErrRawError, "master address should be set")
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
