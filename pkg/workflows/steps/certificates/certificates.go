package certificates

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/util"
)

const StepName = "certificates"

type Step struct {
	template *template.Template
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
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
		template: tpl,
	}
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	// TODO: why does these is set here, not on the building config step?
	config.CertificatesConfig.PrivateIP = config.Node.PrivateIp
	config.CertificatesConfig.PublicIP = config.Node.PublicIp
	config.CertificatesConfig.IsMaster = config.IsMaster

	if !config.IsMaster {
		master := config.GetMaster()
		config.CertificatesConfig.MasterHost = master.PrivateIp
		config.CertificatesConfig.MasterPort = "443"
		config.CertificatesConfig.NodeName = config.Node.Name
	}

	kubeDefaultSvcIp, err := util.GetKubernetesDefaultSvcIP(config.CertificatesConfig.ServicesCIDR)
	if err != nil {
		return errors.Wrapf(err, "get cluster dns ip from the %s subnet", config.CertificatesConfig.ServicesCIDR)
	}
	config.CertificatesConfig.KubernetesSvcIP = kubeDefaultSvcIp.String()

	err = steps.RunTemplate(ctx, s.template, config.Runner, out, config.CertificatesConfig)
	if err != nil {
		return errors.Wrap(err, "write certificates step")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return nil
}
