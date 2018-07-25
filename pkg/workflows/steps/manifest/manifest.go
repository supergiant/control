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

const StepName = "manifest"

type Step struct {
	runner runner.Runner
	script *template.Template
	output io.Writer
}

func New(script *template.Template, outStream io.Writer, cfg *ssh.Config) (*Step, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Step{
		runner: sshRunner,
		script: script,
		output: outStream,
	}

	return t, nil
}

func (j *Step) Run(ctx context.Context, config steps.Config) error {
	err := steps.RunTemplate(ctx, j.script, j.runner, j.output, config.ManifestConfig)

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