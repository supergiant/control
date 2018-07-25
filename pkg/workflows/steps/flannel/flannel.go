package flannel

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "flannel"

type Step struct {
	scriptTemplate *template.Template
	runner         runner.Runner
}

func New(tpl *template.Template, cfg *ssh.Config) (*Step, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	return &Step{
		scriptTemplate: tpl,
		runner:         sshRunner,
	}, nil
}

func (t *Step) Run(ctx context.Context, out io.Writer, config steps.Config) error {
	err := steps.RunTemplate(context.Background(), t.scriptTemplate, t.runner, out, config.FlannelConfig)
	if err != nil {
		return errors.Wrap(err, "Run template has failed for Install flannel job")
	}

	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return ""
}
