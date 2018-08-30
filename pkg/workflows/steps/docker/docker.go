package docker

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "docker"

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
		config.Runner, out, config.DockerConfig)
	if err != nil {
		return errors.Wrap(err, "install docker step")
	}

	return nil
}

func (t *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return nil
}
