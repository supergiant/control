package flannel

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/jobs"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type JobConfig struct {
	Version     string
	Arch        string
	Network     string
	NetworkType string
}

type Job struct {
	scriptTemplate *template.Template
	runner         runner.Runner

	output io.Writer
}

func New(tpl *template.Template, outStream io.Writer, cfg *ssh.Config) (*Job, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	return &Job{
		scriptTemplate: tpl,
		runner:         sshRunner,
		output:         outStream,
	}, nil
}

func (i *Job) InstallFlannel(config JobConfig) error {
	err := jobs.RunTemplate(context.Background(), i.scriptTemplate, i.runner, i.output, config)
	if err != nil {
		return err
	}

	return nil
}
