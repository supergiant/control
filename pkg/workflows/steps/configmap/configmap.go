package configmap

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"text/template"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "configmap"

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(ctx, s.script, config.Runner, out, struct {
		Data string
		Namespace string
	}{
		Data: "",
		Namespace: "default",
	})

	if err != nil {
		return errors.Wrap(err, StepName)
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "create configmap for capacity service"
}

func (s *Step) Depends() []string {
	return nil
}
