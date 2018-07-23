package kubelet

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type Config struct {
	MasterPrivateIP   string
	ProxyPort         string
	EtcdClientPort    string
	KubernetesVersion string
}

type Task struct {
	runner runner.Runner
	config Config
	script *template.Template
	output io.Writer
}

func New(script *template.Template, config Config,
	outStream io.Writer, cfg *ssh.Config) (*Task, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Task{
		runner: sshRunner,
		config: config,
		script: script,
		output: outStream,
	}

	return t, nil
}

func (t *Task) Run(ctx context.Context) error {
	err := steps.RunTemplate(ctx, t.script, t.runner, t.output, t.config)

	if err != nil {
		return errors.Wrap(err, "error running  kubelet template as a command")
	}

	return nil
}
