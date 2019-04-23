package configmap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/runner/dry"
	"github.com/supergiant/control/pkg/storage/memory"
	"github.com/supergiant/control/pkg/workflows"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "configmap"

type bufferCloser struct {
	bytes.Buffer
	err error
}

func (b *bufferCloser) Close() error {
	return b.err
}

type Step struct {
	script *template.Template
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
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

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	dryRunner := dry.NewDryRunner()
	repository := memory.NewInMemoryRepository()

	task, err := workflows.NewTask(config, workflows.NodeTask, repository)

	if err != nil {
		return err
	}

	dryConfig := *config
	dryConfig.Node.PublicIp = "{{ .PublicIp }}"
	dryConfig.Node.PrivateIp = "{{ .PrivateIp }}"
	dryConfig.Runner = dryRunner

	resultChan := task.Run(ctx, dryConfig, &bufferCloser{})

	if err := <-resultChan; err != nil {
		return errors.Wrapf(err, "Run")
	}

	err = steps.RunTemplate(ctx, s.script, config.Runner, out, struct {
		Data      string
		Namespace string
	}{
		Data:      dryRunner.GetOutput(),
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
