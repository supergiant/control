package network

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
)

const StepName = "network"

type Step struct {
	scriptTemplate *template.Template
}

func Init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(tpl *template.Template) *Step {
	return &Step{
		scriptTemplate: tpl,
	}
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), t.scriptTemplate,
		config.Runner, out, config.NetworkConfig)
	if err != nil {
		return errors.Wrap(err, "configure network step")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{etcd.StepName}
}
