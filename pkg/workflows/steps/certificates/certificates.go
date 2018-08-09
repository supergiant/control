package certificates

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "certificates"

type Step struct {
	script *template.Template
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
	config.CertificatesConfig.MasterPrivateIP = config.Node.PrivateIp

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