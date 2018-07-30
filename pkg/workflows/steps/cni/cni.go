package cni

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "cni_tools"

type Step struct {
	script *template.Template
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (j *Step) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(ctx, j.script, config.Runner, out, nil)

	if err != nil {
		return errors.Wrap(err, "error running cni template as a command")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}
