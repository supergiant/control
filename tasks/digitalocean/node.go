package digitalocean

import (
	"os"
	"text/template"

	"bytes"
	"context"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/command"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/tasks"
)

type Job struct {
	runner runner.Runner

	configTemplate *template.Template
	kubeletService *template.Template
	kubeletScript  *template.Template
	proxyScript    *template.Template
}

type jobConfig struct {
	MasterPrivateIP   string
	KubernetesVersion string
	ConfigFile        string
	KubeletService    string
}

func NewJob(cfg *ssh.Config) (*Job, error) {
	configTemplate, err := tasks.ReadTemplate("config.json.tpl", "config")

	if err != nil {
		return nil, err
	}

	kubeletService, err := tasks.ReadTemplate("kubelet_service.conf", "kubelet")

	if err != nil {
		return nil, err
	}

	kubeletScript, err := tasks.ReadTemplate("start_kubelet.sh", "start_kubelet")

	if err != nil {
		return nil, err
	}

	kubeProxyScript, err := tasks.ReadTemplate("start_kube_proxy.sh", "start_kube_proxy")

	if err != nil {
		return nil, err
	}

	sshRunner, err := ssh.NewRunner(os.Stdout, os.Stderr, cfg)

	if err != nil {
		return nil, err
	}

	t := &Job{
		runner:         sshRunner,
		configTemplate: configTemplate,
		kubeletService: kubeletService,
		kubeletScript:  kubeletScript,
		proxyScript:    kubeProxyScript,
	}

	return t, nil
}

func (t *Job) ProvisionNode(k8sVersion, masterPrivateIp string) error {
	buffer := new(bytes.Buffer)
	cfg := jobConfig{
		MasterPrivateIP:   masterPrivateIp,
		KubernetesVersion: k8sVersion,
	}

	err := t.configTemplate.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cfg.ConfigFile = buffer.String()
	buffer.Reset()

	err = t.kubeletService.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cfg.KubeletService = buffer.String()
	buffer.Reset()

	err = t.runTemplate(context.Background(), t.kubeletScript, cfg)

	if err != nil {
		return err
	}

	t.runTemplate(context.Background(), t.proxyScript, cfg)

	if err != nil {
		return err
	}

	return nil
}

func (t *Job) runTemplate(ctx context.Context, tpl *template.Template, cfg jobConfig) error {
	buffer := new(bytes.Buffer)
	err := t.kubeletScript.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	for {
		c, err := buffer.ReadString('\n')

		if err != nil {
			break
		}

		// TODO(stgleb): pass writers for out and err from outside
		cmd := command.NewCommand(context.Background(), c, nil, os.Stdout, os.Stderr)
		err = t.runner.Run(cmd)

		if err != nil {
			return err
		}
	}

	return nil
}
