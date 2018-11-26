package downloadk8sbinary

import (
	"context"
	"io"
	"text/template"
	"fmt"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "download_kubernetes_binary"

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("get template error %v %s", err, StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(tpl *template.Template) *Step {
	return &Step{
		script: tpl,
	}
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), s.script,
		config.Runner, out, config.DownloadK8sBinary)
	if err != nil {
		return errors.Wrap(err, "download k8s binary step")
	}

	return nil
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return nil
}
