package manifest

import (
	"context"
	"text/template"

	"io"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "manifest"

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
	// NOTE(stgleb): This is needed for master node to put advertise address for kube api server.
	config.ManifestConfig.IsMaster = config.IsMaster
	config.ManifestConfig.MasterHost = config.GetMaster().PrivateIp

	err := steps.RunTemplate(ctx, j.script, config.Runner, out, config.ManifestConfig)

	if err != nil {
		return errors.Wrap(err, "write manifest step")
	}

	return nil
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return nil
}
