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

const StepName = "write_certificates"

type Task struct {
	runner runner.Runner
	script *template.Template
	output io.Writer
}

func New(script *template.Template,
	outStream io.Writer, cfg *ssh.Config) (*Task, error) {
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

func (t *Task) Run(ctx context.Context, config steps.Config) error {
	err := steps.RunTemplate(ctx, t.script, t.runner, t.output, config.CertificatesConfig)

	if err != nil {
		return errors.Wrap(err, "error running write certificates template as a command")
	}

	return nil
}
