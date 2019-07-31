package downloadk8sbinary

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "download_kubernetes_binary"

type Config struct {
	K8SVersion      string
	Arch            string
	OperatingSystem string
}

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
	err := steps.RunTemplate(context.Background(), s.script, config.Runner, out, toStepCfg(config))
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
	return "Download kubectl"
}

func (s *Step) Depends() []string {
	return nil
}

func toStepCfg(c *steps.Config) Config {
	return Config{
		K8SVersion:      c.Kube.K8SVersion,
		Arch:            c.Kube.Arch,
		OperatingSystem: c.Kube.OperatingSystem,
	}
}
