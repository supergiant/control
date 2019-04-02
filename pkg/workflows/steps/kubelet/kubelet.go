package kubelet

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/supergiant/control/pkg/clouds"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
)

const (
	StepName = "kubelet"

	// LabelNodeRole specifies the role of a node
	LabelNodeRole = "kubernetes.io/role"
)

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

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(ctx, t.script, config.Runner, out, struct{ Provider clouds.Name }{config.Provider})

	if err != nil {
		return errors.Wrap(err, "install kubelet step")
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
	return "Run kubelet"
}

func (s *Step) Depends() []string {
	return []string{docker.StepName}
}

func getNodeLables(role string) string {
	return labels.Set(map[string]string{
		LabelNodeRole: role,
	}).String()
}
