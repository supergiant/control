package addons

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/dashboard"
)

const StepName = "addons"

var (
	Default = []string{
		dashboard.StepName,
	}
)

type Step struct {
}

func (s Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	for _, name := range config.Kube.Addons {
		s := steps.GetStep(name)
		if s == nil {
			continue
		}
		if err := s.Run(ctx, out, config); err != nil {
			return errors.Wrapf(err, "install kubernetes addons: %s", name)
		}
	}
	return nil
}

func (s Step) Name() string {
	return StepName
}

func (s Step) Description() string {
	return "Install kubernetes addons"
}

func (s Step) Depends() []string {
	return nil
}
