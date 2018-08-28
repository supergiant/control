package poststart

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubelet"
	"time"
)

const StepName = "poststart"

type Step struct {
	script *template.Template
}

func Init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (j *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	ctx2, _ := context.WithTimeout(ctx, time.Duration(config.PostStartConfig.Timeout)*time.Second)
	err := steps.RunTemplate(ctx2, j.script, config.Runner, out, config.PostStartConfig)

	if err != nil {
		return errors.Wrap(err, "run post start script step")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{kubelet.StepName}
}
