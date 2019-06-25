package upgrade

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "upgrade"

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
	err := steps.RunTemplate(ctx, s.script, config.Runner, out, struct {
		K8SVersion  string
		IsBootstrap bool
		IsMaster    bool
	}{
		K8SVersion:  config.K8sVersion,
		IsBootstrap: config.IsBootstrap,
		IsMaster:    config.IsMaster,
	})

	if err != nil {
		return errors.Wrap(err, "upgrade step has failed")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "Upgrade k8s node"
}

func (s *Step) Depends() []string {
	return nil
}
