package tiller

import (
	"context"
	"text/template"

	"github.com/pkg/errors"

	"io"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "install_tiller"

type Task struct {
	script *template.Template
}

func New(script *template.Template) (*Task, error) {

	t := &Task{
		script: script,
	}

	return t, nil
}

func (j *Task) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), j.script, config.Runner, out, config.TillerConfig)

	if err != nil {
		return errors.Wrap(err, "error running tiller template as a command")
	}

	return nil
}
