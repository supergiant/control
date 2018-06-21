package node

import (
	"bytes"
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

	out io.Writer
	err io.Writer
}

type JobConfig struct {
	MasterPrivateIP   string
	ProxyPort string
	EtcdPort string
	KubernetesVersion string
}

func NewJob(startKubeletFileName, startProxyFileName string,
	outStream, errStream io.Writer, cfg *ssh.Config) (*Job, error) {
	kubeletScript, err := jobs.ReadTemplate(startKubeletFileName, "start_kubelet")

	if err != nil {
		return nil, errors.Wrap(err, "error reading kubelet template")
	}

	kubeProxyScript, err := jobs.ReadTemplate(startProxyFileName, "start_kube_proxy")

	if err != nil {
		return nil, errors.Wrap(err, "error reading proxy template")
	}

	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	t := &Job{
		runner:        sshRunner,
		kubeletScript: kubeletScript,
		proxyScript:   kubeProxyScript,

		out: outStream,
		err: errStream,
	}

	return t, nil
}

func (j *Job) ProvisionNode(config JobConfig) error {
	err := j.runTemplate(context.Background(), j.kubeletScript, config)

	if err != nil {
		return errors.Wrap(err, "error running  kubelet template as a command")
	}

	j.runTemplate(context.Background(), j.proxyScript, config)

	if err != nil {
		return errors.Wrap(err, "error running proxy template as a command")
	}

	return nil
}

// TODO(stgleb): maybe it can be moved to util and not to be a method of job
func (j *Job) runTemplate(ctx context.Context, tpl *template.Template, cfg JobConfig) error {
	buffer := new(bytes.Buffer)
	err := tpl.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cmd := runner.NewCommand(ctx, buffer.String(), j.out, j.err)
	err = j.runner.Run(cmd)

	if err != nil {
		return err
	}

	return nil
}
