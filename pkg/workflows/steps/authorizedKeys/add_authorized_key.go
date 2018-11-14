package authorizedKeys

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type Step struct {
	script *template.Template
}

const AddAuthorizedKeyStepName = "add_authorized_keys"

func Init() {
	steps.RegisterStep(AddAuthorizedKeyStepName, NewAddAuthorizedKeys(templatemanager.GetTemplate(AddAuthorizedKeyStepName)))
}

func NewAddAuthorizedKeys(script *template.Template) *Step {
	return &Step{
		script: script,
	}
}

func (s *Step) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	log.Infof("[%s] - adding user's public key to the node", s.Name())
	if cfg.SshConfig.PublicKey != "" {
		err := steps.RunTemplate(ctx, s.script, cfg.Runner, w, cfg.SshConfig)
		if err != nil {
			return errors.Wrap(err, "add authorized key step")
		}
	} else {
		log.Infof("[%s] - no public key provided, skipping...", s.Name())
	}

	return nil
}

func (*Step) Name() string {
	return AddAuthorizedKeyStepName
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
