package tiller

import (
	"context"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
)

const StepName = "install_tiller"

type Task struct {
	runner runner.Runner
	script *template.Template
}

func New(script *template.Template, cfg *ssh.Config) (*Task, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Task{
		runner: sshRunner,
		script: script,
	}

	return t, nil
}

func (j *Task) Run(ctx context.Context,out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), j.script, j.runner, out, config.TillerConfig)

	if err != nil {
		return errors.Wrap(err, "error running tiller template as a command")
	}

	return nil
}
