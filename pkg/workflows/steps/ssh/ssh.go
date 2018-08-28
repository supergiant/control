package ssh

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "ssh"

type Step struct{}

func Init() {
	steps.RegisterStep(StepName, &Step{})
}

func (s *Step) Run(ctx context.Context, writer io.Writer, config *steps.Config) error {
	var err error

	cfg := ssh.Config{
		Host:    config.Node.PublicIp,
		Port:    config.SshConfig.Port,
		User:    config.SshConfig.User,
		Timeout: config.SshConfig.Timeout,
		// TODO(stgleb): Pass ssh key id instead of key itself
		Key: config.SshConfig.PrivateKey,
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

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{"node"}
}
