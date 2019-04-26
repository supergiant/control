package kubeadm

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/supergiant/control/pkg/clouds"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
	config.KubeadmConfig.Provider = toCloudProviderOpt(config.Provider)
	config.KubeadmConfig.IsBootstrap = config.IsBootstrap
	config.KubeadmConfig.IsMaster = config.IsMaster
	config.KubeadmConfig.InternalDNSName = config.InternalDNSName
	config.KubeadmConfig.ExternalDNSName = config.ExternalDNSName
	config.KubeadmConfig.Token = config.BootstrapToken

	// NOTE(stgleb): Kubeadm accepts only ipv4 or ipv6 addresses as advertise address
	if config.IsBootstrap {
		config.KubeadmConfig.AdvertiseAddress = config.Node.PrivateIp
	}

	// TODO: fix azure load balancer
	if config.Provider == clouds.Azure {
		config.KubeadmConfig.InternalDNSName = config.Node.PrivateIp
		if config.IsBootstrap {
			config.KubeadmConfig.InternalDNSName = config.Node.PrivateIp
		}
	}

	logrus.Debugf("kubeadm step: %s cluster: isBootstrap=%t extDNS=%s intDNS=%s",
		config.ClusterID, config.KubeadmConfig.IsBootstrap, config.KubeadmConfig.ExternalDNSName,
		config.KubeadmConfig.InternalDNSName)

	config.KubeadmConfig.IsMaster = config.IsMaster
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

// TODO: cloud profiles is deprecated by kubernetes, use controller-managers
func toCloudProviderOpt(cloudName clouds.Name) string {
	switch cloudName {
	case clouds.AWS:
		return "aws"
	case clouds.GCE:
		return "gce"
	}
	return ""
}
