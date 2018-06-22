package jobs

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type FlannelJobConfig struct {
	Version     string
	Arch        string
	Network     string
	NetworkType string
}

type FlannelJob struct {
	scriptTemplate *template.Template
	runner         runner.Runner

	output io.Writer
}

func NewFlannelJob(tpl *template.Template, outStream io.Writer, cfg *ssh.Config) (*FlannelJob, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	return &FlannelJob{
		scriptTemplate: tpl,
		runner:         sshRunner,
		output:         outStream,
	}, nil
}

func (i *FlannelJob) InstallFlannel(config FlannelJobConfig) error {
	err := runTemplate(context.Background(), i.scriptTemplate, i.runner, i.output, config)
	if err != nil {
		return errors.Wrap(err, "error running template for flannel")
	}

	return nil
}
