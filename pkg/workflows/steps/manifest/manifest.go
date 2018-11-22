package manifest

import (
	"context"
	"io"
	"text/template"
	"fmt"

	"github.com/pkg/errors"

	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
)

const StepName = "manifest"

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
	return []string{certificates.StepName}
}
