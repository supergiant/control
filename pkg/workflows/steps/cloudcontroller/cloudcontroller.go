package cloudcontroller

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "cloudcontroller"

type Config struct {
	Provider      string
	DOAccessToken string
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

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), s.script, config.Runner, out, toStepCfg(config))

	if err != nil {
		return errors.Wrap(err, "install cloud-controller-manager")
	}

	return nil
}

func (s *Step) Rollback(ctx context.Context, out io.Writer, config *steps.Config) error {
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "create cloud-controller-manager"
}

func (s *Step) Depends() []string {
	return nil
}

func toStepCfg(c *steps.Config) Config {
	return Config{
		Provider:      string(c.Kube.Provider),
		DOAccessToken: c.DigitalOceanConfig.AccessToken,
	}
}
