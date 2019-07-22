package tiller

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
)

const (
	StepName = "tiller"
)

type Config struct {
	HelmVersion     string
	RBACEnabled     bool
	OperatingSystem string
	Arch            string
}

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

func (j *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), j.script, config.Runner, out, toStepCfg(config))

	if err != nil {
		return errors.Wrap(err, "install tiller step")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{poststart.StepName}
}

func toStepCfg(c *steps.Config) Config {
	return Config{
		HelmVersion:     c.Kube.HelmVersion,
		OperatingSystem: c.Kube.OperatingSystem,
		Arch:            c.Kube.Arch,
		RBACEnabled:     c.Kube.RBACEnabled,
	}
}
