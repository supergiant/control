package certificates

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

const StepName = "certificates"

type Step struct {
	script *template.Template
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func Init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	if master := config.GetMaster(); master != nil {
		config.CertificatesConfig.PrivateIP = master.PrivateIp
		config.CertificatesConfig.PublicIP = master.PublicIp
	} else {
		return sgerrors.ErrNotFound
	}
	err := steps.RunTemplate(ctx, s.script,
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
