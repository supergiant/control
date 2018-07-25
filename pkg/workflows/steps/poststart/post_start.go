package tiller

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "post_start"

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
	err := steps.RunTemplate(context.Background(), j.script, config.Runner, out, config.PostStartConfig)

	if err != nil {
		return errors.Wrap(err, "error running post start template as a command")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}
