package etcd

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	tm "github.com/supergiant/supergiant/pkg/templatemanager"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
)

const StepName = "etcd"

type Step struct {
	scriptTemplate *template.Template
}

func Init() {
	steps.RegisterStep(StepName, New(tm.GetTemplate(StepName)))
}

func New(tpl *template.Template) *Step {
	return &Step{
		scriptTemplate: tpl,
	}
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.EtcdConfig.Name = config.Node.Id
	config.EtcdConfig.AdvertiseHost = config.Node.PrivateIp
	ctx2, _ := context.WithTimeout(ctx, config.EtcdConfig.Timeout)

	err := steps.RunTemplate(ctx2, s.scriptTemplate,
		config.Runner, out, config.EtcdConfig)
	if err != nil {
		return errors.Wrap(err, "install etcd step")
	}

	// Notify other task that one master is ready
	if config.IsMaster {
		config.Done()
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
