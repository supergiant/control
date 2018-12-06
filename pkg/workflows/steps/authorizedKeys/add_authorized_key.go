package authorizedKeys

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type Step struct {
	script *template.Template
}

const StepName = "add_authorized_keys"

func Init() {
	tpl, err := templatemanager.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}
	steps.RegisterStep(StepName, NewAddAuthorizedKeys(tpl))
}

func NewAddAuthorizedKeys(script *template.Template) *Step {
	return &Step{
		script: script,
	}
}

func (s *Step) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	if cfg.SshConfig.PublicKey != "" {
		err := steps.RunTemplate(ctx, s.script, cfg.Runner, w, cfg.SshConfig, cfg.DryRun)
		if err != nil {
			return errors.Wrap(err, "add authorized key step")
		}
	}
	return nil
}

func (*Step) Name() string {
	return StepName
}

func (*Step) Description() string {
	return "adds ssh public key to the authorized keys file"
}

func (*Step) Depends() []string {
	return nil
}

func (*Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
