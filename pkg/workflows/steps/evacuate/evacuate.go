package evacuate

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"text/template"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "evacuate"

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
	if !config.IsMaster {
		err := steps.RunTemplate(ctx, s.script, config.Runner, out, struct{PrivateIP string}{config.Node.PrivateIp})

		if err != nil {
			return errors.Wrap(err, "evacuate step has failed")
		}
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "Evacuate k8s node"
}

func (s *Step) Depends() []string {
	return nil
}
