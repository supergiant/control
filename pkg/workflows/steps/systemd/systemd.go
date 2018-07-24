package tiller

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type Task struct {
	runner runner.Runner
	script *template.Template
	output io.Writer
}

func New(script *template.Template, outStream io.Writer, cfg *ssh.Config) (*Task, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Task{
		runner: sshRunner,
		script: script,
		output: outStream,
	}

	return t, nil
}

func (j *Task) Run(ctx context.Context, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), j.script, j.runner, j.output, config.KubeletSystemdServiceConfig)

	if err != nil {
		return errors.Wrap(err, "error running post start template as a command")
	}

	return nil
}
