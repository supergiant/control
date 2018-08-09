package kubelet

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
	"github.com/supergiant/supergiant/pkg/workflows/steps/manifest"
)

const StepName = "kubelet"

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

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.KubeletConfig.MasterPrivateIP = config.Node.PrivateIp
	err := steps.RunTemplate(ctx, t.script, config.Runner, out, config.KubeletConfig)

	if err != nil {
		return errors.Wrap(err, "install kubelet step")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{docker.StepName, manifest.StepName}
}