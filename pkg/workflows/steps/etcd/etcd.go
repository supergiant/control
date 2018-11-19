package etcd

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
	"fmt"
)

const StepName = "etcd"

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

func New(tpl *template.Template) *Step {
	return &Step{
		script: tpl,
	}
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.EtcdConfig.Name = config.Node.ID
	config.EtcdConfig.AdvertiseHost = config.Node.PrivateIp
	ctx2, _ := context.WithTimeout(ctx, config.EtcdConfig.Timeout)

	err := steps.RunTemplate(ctx2, s.script,
		config.Runner, out, config.EtcdConfig)
	if err != nil {
		return errors.Wrap(err, "install etcd step")
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
	return []string{docker.StepName}
}
