package drain

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/sgerrors"
)

const StepName = "drain"

type Step struct {
	script *template.Template
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	masterNode := config.GetMaster()

	if masterNode == nil {
		return errors.Wrapf(sgerrors.ErrNotFound, "master node not found")
	}

	cfg := ssh.Config{
		Host:    masterNode.PublicIp,
		Port:    config.SshConfig.Port,
		User:    config.SshConfig.User,
		Timeout: 10,
		Key:     []byte(config.SshConfig.BootstrapPrivateKey),
	}

	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return errors.Wrapf(err, "create ssh runner")
	}

	err = steps.RunTemplate(ctx, s.script, sshRunner, out, config.DrainConfig)

	if err != nil {
		return errors.Wrap(err, "drain step")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "drain resources from a node"
}

func (s *Step) Depends() []string {
	return nil
}
