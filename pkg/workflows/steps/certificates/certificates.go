package certificates

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"fmt"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
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

	err := steps.RunTemplate(ctx, s.template,
		config.Runner, out, config.CertificatesConfig)

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
