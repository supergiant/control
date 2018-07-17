package node

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/jobs"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type Job struct {
	runner runner.Runner

	kubeletScript *template.Template
	proxyScript   *template.Template

	output io.Writer
}

type JobConfig struct {
	MasterPrivateIP   string
	ProxyPort         string
	EtcdClientPort    string
	KubernetesVersion string
}

func New(startKubeletTemplate, startProxyTemplate *template.Template,
	outStream io.Writer, cfg *ssh.Config) (*Job, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Job{
		runner:        sshRunner,
		kubeletScript: startKubeletTemplate,
		proxyScript:   startProxyTemplate,

		output: outStream,
	}

	return t, nil
}

func (j *Job) ProvisionNode(config JobConfig) error {
	err := jobs.RunTemplate(context.Background(), j.kubeletScript, j.runner, j.output, config)

	if err != nil {
		return errors.Wrap(err, "error running  kubelet template as a command")
	}

	err = jobs.RunTemplate(context.Background(), j.proxyScript, j.runner, j.output, config)

	if err != nil {
		return errors.Wrap(err, "error running proxy template as a command")
	}

	return nil
}
