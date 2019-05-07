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
	config.KubeadmConfig.AdvertiseAddress = config.Node.PrivateIp
	config.KubeadmConfig.NodeIp = config.Node.PrivateIp
	config.KubeadmConfig.PrivateIp = config.Node.PrivateIp
	config.KubeadmConfig.JoinAddress = config.InternalDNSName

	// NOTE(stgleb): GCE sends traffic from load balancer like it originates from outside world,
	// so if we do not listen to 0.0.0.0 kube-apiserver will refuse all connection outside of instance
	if config.Provider == clouds.GCE {
		config.KubeadmConfig.PrivateIp = "0.0.0.0"
		config.KubeadmConfig.AdvertiseAddress = "0.0.0.0"
		config.KubeadmConfig.NodeIp = config.Node.PublicIp

		if !config.IsBootstrap && config.IsMaster {
			master := config.GetMaster()

			if master == nil {
				return errors.New("master not found")
			}

			config.KubeadmConfig.JoinAddress = master.PublicIp
		}
	}

	logrus.Debugf("kubeadm step: %s cluster: isBootstrap=%t extDNS=%s intDNS=%s",
		config.ClusterID, config.KubeadmConfig.IsBootstrap, config.KubeadmConfig.ExternalDNSName,
		config.KubeadmConfig.InternalDNSName)

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
