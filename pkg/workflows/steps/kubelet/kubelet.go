package kubelet

import (
	"context"
	"fmt"
	"io"
	"net"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/util"
)

const (
	StepName = "kubelet"

	// LabelNodeRole specifies the role of a node
	LabelNodeRole = "kubernetes.io/role"
)

type Config struct {
	IsMaster     bool   `json:"isMaster"`
	ServicesCIDR string `json:"servicesCIDR"`
	PublicIP     string `json:"publicIp"`
	PrivateIP    string `json:"privateIp"`

	LoadBalancerHost string `json:"loadBalancerHost"`
	APIServerPort    int64  `json:"apiserverPort"`
	NodeName         string `json:"nodeName"`
	UserName         string `json:"userName"`

	// TODO: this shouldn't be a part of SANs
	// https://kubernetes.io/docs/setup/certificates/#all-certificates
	KubernetesSvcIP string `json:"kubernetesSvcIp"`

	AdminCert string `json:"adminCert"`
	AdminKey  string `json:"adminKey"`
	CACert    string `json:"caCert"`
	CAKey     string `json:"caKey"`
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
	c, err := toStepCfg(config)
	if err != nil {
		return errors.Wrap(err, "build step config")
	}

	err = steps.RunTemplate(ctx, t.script, config.Runner, out, c)
	if err != nil {
		return errors.Wrap(err, "install kubelet step")
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
	return "Run kubelet"
}

func (s *Step) Depends() []string {
	return []string{docker.StepName}
}

func toStepCfg(c *steps.Config) (Config, error) {
	var svcIP net.IP
	var err error
	if len(c.Kube.ServicesCIDR) > 0 {
		svcIP, err = util.GetKubernetesDefaultSvcIP(c.Kube.ServicesCIDR)
		if err != nil {
			return Config{}, errors.Wrapf(err, "get cluster dns ip from the %s subnet", c.Kube.ServicesCIDR)
		}
	}

	return Config{
		IsMaster:         c.IsMaster,
		LoadBalancerHost: c.Kube.InternalDNSName,
		NodeName:         c.Node.Name,
		PrivateIP:        c.Node.PrivateIp,
		PublicIP:         c.Node.PublicIp,
		CACert:           c.Kube.Auth.CACert,
		CAKey:            c.Kube.Auth.CAKey,
		AdminCert:        c.Kube.Auth.AdminCert,
		AdminKey:         c.Kube.Auth.AdminKey,
		APIServerPort:    c.Kube.APIServerPort,
		UserName:         c.Kube.SSHConfig.User,
		ServicesCIDR:     c.Kube.ServicesCIDR,
		KubernetesSvcIP:  svcIP.String(),
	}, nil
}
