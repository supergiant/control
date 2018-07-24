package flannel

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type Config struct {
	Version     string
	Arch        string
	Network     string
	NetworkType string
}

type Task struct {
	scriptTemplate *template.Template
	runner         runner.Runner

	output io.Writer
}

func New(tpl *template.Template, outStream io.Writer, cfg *ssh.Config) (*Task, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	return &Task{
		scriptTemplate: tpl,
		runner:         sshRunner,
		output:         outStream,
	}, nil
}

func (t *Task) Run(ctx context.Context, config workflows.Config) error {
	err := steps.RunTemplate(context.Background(), t.scriptTemplate, t.runner, t.output, config.FlannelConfig)
	if err != nil {
		return errors.Wrap(err, "Run template has failed for Install flannel job")
	}

	return nil
}
