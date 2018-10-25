package flannel

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
	"github.com/supergiant/supergiant/pkg/workflows/steps/network"
)

const StepName = "flannel"

type Step struct {
	script *template.Template
}

func Init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(tpl *template.Template) *Step {
	return &Step{
		script: tpl,
	}
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), t.script,
		config.Runner, out, config.FlannelConfig)
	if err != nil {
		return errors.Wrap(err, "install flannel step")
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

func (t *Step) Depends() []string {
	return []string{etcd.StepName, network.StepName}
}
