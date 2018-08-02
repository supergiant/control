package flannel

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/runner"
)

const StepName = "flannel"

type Task struct {
	scriptTemplate *template.Template
}

func init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(tpl *template.Template) *Step {
	return &Task{
		scriptTemplate: tpl,
	}
}

func (t *Task) Run(ctx context.Context, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), t.scriptTemplate, t.runner, t.output, config.FlannelConfig)
	if err != nil {
		return errors.Wrap(err, "install flannel step")
	}

	return nil
}

func (t *Task) Name() string {
	return StepName
}

func (t *Task) Description() string {
	return ""
}
