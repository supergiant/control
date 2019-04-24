package ssh

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/runner/dry"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/azure"
)

const StepName = "ssh"

type Step struct{}

func Init() {
	steps.RegisterStep(StepName, &Step{})
}

func (s *Step) Run(ctx context.Context, writer io.Writer, config *steps.Config) error {
	var err error

	if config.DryRun {
		if config.Runner == nil {
			config.Runner = dry.NewDryRunner()
		}
		return nil
	}

	if config.Provider == clouds.AWS {
		//on aws default user name on ubuntu images are not root but ubuntu
		//https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html
		// TODO: this should be set by provisioner
		config.Kube.SSHConfig.User = "ubuntu"
	} else if config.Provider == clouds.Azure {
		config.Kube.SSHConfig.User = azure.OSUser
	}

	cfg := ssh.Config{
		Host:    config.Node.PublicIp,
		Port:    config.Kube.SSHConfig.Port,
		User:    config.Kube.SSHConfig.User,
		Timeout: config.Kube.SSHConfig.Timeout,
		// TODO(stgleb): Use secure storage for private keys instead carrying them in plain text
		Key: []byte(config.Kube.SSHConfig.BootstrapPrivateKey),
	}

	config.Runner, err = ssh.NewRunner(cfg)
	if err != nil {
		return errors.Wrap(err, "ssh config step")
	}
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{"node"}
}
