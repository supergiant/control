package storageclass

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "storageclass"

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

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s *Step) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	log.Infof("[%s] - applying default storage class", s.Name())

	err := steps.RunTemplate(ctx, s.script, cfg.Runner, w, cfg)
	if err != nil {
		return errors.Wrap(err, "apply default storage class step")
	}

	return nil
}

func (*Step) Name() string {
	return StepName
}

func (*Step) Description() string {
	return ""
}

func (*Step) Depends() []string {
	return nil
}

func (*Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
