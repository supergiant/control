package jobs

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type NodeJob struct {
	runner runner.Runner

	kubeletScript *template.Template
	proxyScript   *template.Template

	out io.Writer
	err io.Writer
}

type JobConfig struct {
	MasterPrivateIP   string
	ProxyPort         string
	EtcdClientPort    string
	KubernetesVersion string
}

func NewJob(startKubeletTemplate, startProxyTemplate *template.Template,
	outStream, errStream io.Writer, cfg *ssh.Config) (*NodeJob, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &NodeJob{
		runner:        sshRunner,
		kubeletScript: startKubeletTemplate,
		proxyScript:   startProxyTemplate,

		out: outStream,
		err: errStream,
	}

	return t, nil
}

func (j *NodeJob) ProvisionNode(config JobConfig) error {
	err := runTemplate(context.Background(), j.kubeletScript, j.runner, j.out, j.err, config)

	if err != nil {
		return errors.Wrap(err, "error running  kubelet template as a command")
	}

	err = runTemplate(context.Background(), j.proxyScript, j.runner, j.out, j.err, config)

	if err != nil {
		return errors.Wrap(err, "error running proxy template as a command")
	}

	return nil
}
