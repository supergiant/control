package digitalocean

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/supergiant/pkg/jobs"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type Job struct {
	runner runner.Runner

	configTemplate *template.Template
	kubeletService *template.Template
	kubeletScript  *template.Template
	proxyScript    *template.Template

	out io.Writer
	err io.Writer
}

type jobConfig struct {
	MasterPrivateIP   string
	KubernetesVersion string
	ConfigFile        string
	KubeletService    string
}

func NewJob(configFileName, kubeletConfigFileName, startKubeletFileName, startProxyFileName string,
	outStream, errStream io.Writer, cfg *ssh.Config) (*Job, error) {
	configTemplate, err := jobs.ReadTemplate(configFileName, "config")

	if err != nil {
		return nil, err
	}

	kubeletService, err := jobs.ReadTemplate(kubeletConfigFileName, "kubelet")

	if err != nil {
		return nil, err
	}

	kubeletScript, err := jobs.ReadTemplate(startKubeletFileName, "start_kubelet")

	if err != nil {
		return nil, err
	}

	kubeProxyScript, err := jobs.ReadTemplate(startProxyFileName, "start_kube_proxy")

	if err != nil {
		return nil, err
	}

	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, err
	}

	t := &Job{
		runner:         sshRunner,
		configTemplate: configTemplate,
		kubeletService: kubeletService,
		kubeletScript:  kubeletScript,
		proxyScript:    kubeProxyScript,

		out: outStream,
		err: errStream,
	}

	return t, nil
}

func (j *Job) ProvisionNode(k8sVersion, masterPrivateIp string) error {
	buffer := new(bytes.Buffer)
	cfg := jobConfig{
		MasterPrivateIP:   masterPrivateIp,
		KubernetesVersion: k8sVersion,
	}

	err := j.configTemplate.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cfg.ConfigFile = buffer.String()
	buffer.Reset()

	err = j.kubeletService.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cfg.KubeletService = buffer.String()
	buffer.Reset()

	err = j.runTemplate(context.Background(), j.kubeletScript, cfg)

	if err != nil {
		return err
	}

	j.runTemplate(context.Background(), j.proxyScript, cfg)

	if err != nil {
		return err
	}

	return nil
}

// TODO(stgleb): maybe it can be moved to util and not to be a method of job
func (j *Job) runTemplate(ctx context.Context, tpl *template.Template, cfg jobConfig) error {
	buffer := new(bytes.Buffer)
	err := tpl.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cmd := runner.NewCommand(context.Background(), buffer.String(), j.out, j.err)
	err = j.runner.Run(cmd)

	if err != nil {
		return err
	}

	return nil
}
