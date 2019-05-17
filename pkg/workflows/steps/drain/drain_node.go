package drain

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/sgerrors"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "drain"

type Step struct {
	script    *template.Template
	getRunner func(string, *steps.Config) (runner.Runner, error)
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
		getRunner: func(masterIp string, config *steps.Config) (runner.Runner, error) {
			if config.Provider == clouds.AWS {
				//on aws default user name on ubuntu images are not root but ubuntu
				//https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html
				config.Kube.SSHConfig.User = "ubuntu"
			}

			cfg := ssh.Config{
				Host:    masterIp,
				Port:    config.Kube.SSHConfig.Port,
				User:    config.Kube.SSHConfig.User,
				Timeout: 10,
				Key:     []byte(config.Kube.SSHConfig.BootstrapPrivateKey),
			}

			sshRunner, err := ssh.NewRunner(cfg)

			if err != nil {
				return nil, errors.Wrapf(err, "create ssh runner")
			}

			return sshRunner, err
		},
	}

	return t
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	masterNode := config.GetMaster()

	if masterNode == nil {
		return errors.Wrapf(sgerrors.ErrNotFound, "master node not found")
	}

	r, err := s.getRunner(masterNode.PublicIp, config)

	if err != nil {
		return errors.Wrapf(err, "get runner")
	}

	err = steps.RunTemplate(ctx, s.script, r, out, config.DrainConfig)

	if err != nil {
		logrus.Errorf("drain step has failed with %v", err)
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
