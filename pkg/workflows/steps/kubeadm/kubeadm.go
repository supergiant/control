package kubeadm

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
)

const (
	StepName = "kubeadm"
)

type Config struct {
	K8SVersion      string
	KubeadmVersion  string
	IsBootstrap     bool
	IsMaster        bool
	InternalDNSName string
	ExternalDNSName string
	Token           string
	CACertHash      string
	CertificateKey  string
	CIDR            string
	ServiceCIDR     string
	UserName        string
	Provider        string
	APIServerPort   int
	NodeIp          string
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

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	logrus.Debugf("kubeadm step: %s cluster: isBootstrap=%t extDNS=%s intDNS=%s",
		config.Kube.ID, config.IsBootstrap, config.Kube.ExternalDNSName,
		config.Kube.InternalDNSName)

	err := steps.RunTemplate(ctx, t.script, config.Runner, out, toStepCfg(config))

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

func toStepCfg(c *steps.Config) Config {
	return Config{
		KubeadmVersion:  "1.15.1", // TODO(stgleb): get it from available versions once we have them
		K8SVersion:      c.Kube.K8SVersion,
		IsBootstrap:     c.IsBootstrap,
		IsMaster:        c.IsMaster,
		InternalDNSName: c.Kube.InternalDNSName,
		ExternalDNSName: c.Kube.ExternalDNSName,
		Token:           c.Kube.BootstrapToken,
		CACertHash:      c.Kube.Auth.CACertHash,
		CertificateKey:  c.Kube.Auth.CertificateKey,
		CIDR:            c.Kube.Networking.CIDR,
		ServiceCIDR:     c.Kube.ServicesCIDR,
		UserName:        clouds.OSUser,
		Provider:        toCloudProviderOpt(c.Kube.Provider),
		APIServerPort:   c.Kube.APIServerPort,
		NodeIp:          c.Node.PrivateIp,
	}
}
