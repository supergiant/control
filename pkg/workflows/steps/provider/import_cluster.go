package provider

import (
	"context"
	"fmt"
	"io"


	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/clouds"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
)

const (
	ImportClusterStepName = "import_cluster"
)

type ImportClusterStep struct {
}

func (s ImportClusterStep) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	var importSteps []steps.Step

	switch cfg.Provider {
	case clouds.AWS:
		importSteps = []steps.Step{
			steps.GetStep(amazon.ImportClusterMachinesStepName),
			steps.GetStep(amazon.ImportSubnetsStepName),
		}
	default:
		return errors.New(fmt.Sprintf("unsupported provider: %s", cfg.Provider))
	}

	for _, s := range importSteps {
		if err := s.Run(ctx, out, cfg); err != nil {
			return errors.Wrap(err, ImportClusterStepName)
		}
	}

	return nil
}

func (s ImportClusterStep) Name() string {
	return ImportClusterStepName
}

func (s ImportClusterStep) Description() string {
	return ImportClusterStepName
}

func (s ImportClusterStep) Depends() []string {
	return nil
}

func (s ImportClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
