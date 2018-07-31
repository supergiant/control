package tiller

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/workflows/steps"

	"text/template"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
)

const (
	StepName    = "install_tiller"
	TemplateDir = "scripts"
)

type Step struct {
	script *template.Template
}

func init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (j *Step) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), j.script, config.Runner, out, config.TillerConfig)

	if err != nil {
		return errors.Wrap(err, "error running tiller templatemanager as a command")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return ""
}
