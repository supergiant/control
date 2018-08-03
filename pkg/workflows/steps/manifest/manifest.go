package manifest

import (
	"context"
	"text/template"

	"github.com/pkg/errors"

	"io"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "manifest"

type Step struct {
	script *template.Template
}

func init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (j *Step) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(ctx, j.script, config.Runner, out, config.ManifestConfig)

	if err != nil {
		return errors.Wrap(err, "write manifest step")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}
