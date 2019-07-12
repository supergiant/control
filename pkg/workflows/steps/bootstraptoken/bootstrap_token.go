package bootstraptoken

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "bootstrap_token"

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
	if config.IsBootstrap {
		token, err := GenerateBootstrapToken()
		config.BootstrapToken = token

		if err != nil {
			return errors.Wrapf(err, "generate bootstrap token")
		}

		logrus.Debugf("Create bootstrap token %s", config.BootstrapToken)
		// NOTE(stgleb): Reuse KubeadmConfig.Token field to avoid
		err = steps.RunTemplate(ctx, s.script, config.Runner, out, struct {
			IsBootstrap    bool
			Token          string
			CertificateKey string
			IsImport       bool
		}{
			IsBootstrap:    config.IsBootstrap,
			Token:          config.BootstrapToken,
			CertificateKey: config.KubeadmConfig.CertificateKey,
			IsImport:       config.IsImport,
		})

		if err != nil {
			return errors.Wrap(err, "create bootstrap token")
		}
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "create bootstrap token for cluster"
}

func (s *Step) Depends() []string {
	return nil
}
