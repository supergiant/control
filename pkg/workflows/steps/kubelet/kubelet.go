package kubelet

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/supergiant/control/pkg/model"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
)

const (
	StepName = "kubelet"

	// nodeLabelRole specifies the role of a node
	nodeLabelRole = "kubernetes.io/role"
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
	err := steps.RunTemplate(ctx, t.script, config.Runner, out, nil)

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
	return []string{docker.StepName, manifest.StepName}
}

func getNodeLables(role string) string {
	return labels.Set(map[string]string{
		nodeLabelRole: role,
	}).String()
}

// TODO: role should be a port of config, it's used by a few tasks
func toRole(isMaster bool) string {
	if isMaster {
		return string(model.RoleMaster)
	}
	return string(model.RoleNode)
}
