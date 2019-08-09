package certificates

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "certificates"

type Config struct {
	IsBootstrap bool
	CACert      string
	CAKey       string
}

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
	err := steps.RunTemplate(ctx, s.template, config.Runner, out, toStepCfg(config))
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

func toStepCfg(c *steps.Config) Config {
	return Config{
		IsBootstrap: c.IsBootstrap,
		CACert:      c.Kube.Auth.CACert,
		CAKey:       c.Kube.Auth.CAKey,
	}
}
