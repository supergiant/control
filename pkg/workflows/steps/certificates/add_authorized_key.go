package certificates

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type AddAuthorizedKeyStep struct {
	script *template.Template
}

const AddAuthorizedKeyStepName = "add_authorized_keys"

func InitAddAuthorizedKeys() {
	steps.RegisterStep(AddAuthorizedKeyStepName, NewAddAuthorizedKeys(templatemanager.GetTemplate(AddAuthorizedKeyStepName)))
}

func NewAddAuthorizedKeys(script *template.Template) *AddAuthorizedKeyStep {
	return &AddAuthorizedKeyStep{
		script: script,
	}
}

func (s *AddAuthorizedKeyStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
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

func (*AddAuthorizedKeyStep) Name() string {
	return AddAuthorizedKeyStepName
}

func (*AddAuthorizedKeyStep) Description() string {
	return "adds ssh public key to the authorized keys file"
}

func (*AddAuthorizedKeyStep) Depends() []string {
	return nil
}

func (*AddAuthorizedKeyStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
