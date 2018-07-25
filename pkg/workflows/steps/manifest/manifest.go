package tiller

import (
	"context"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
)

const StepName = "manifest"

type Step struct {
	script *template.Template
}

func New(script *template.Template) (*Step, error) {
	t := &Step{
		script: script,
	}

	return t, nil
}

func (j *Step) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(ctx, j.script, config.Runner, out, config.ManifestConfig)

	if err != nil {
		return errors.Wrap(err, "error running write certificates template as a command")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}
