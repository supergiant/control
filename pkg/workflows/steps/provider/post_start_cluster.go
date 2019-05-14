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
	PostStartCluster = "post_start_cluster"
)

type StepPostStartCluster struct {}

func (s StepPostStartCluster) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	postStartClusterStep, err := postStartCluster(cfg.Provider)
	if err != nil {
		return errors.Wrap(err, PostStartCluster)
	}

	for _, s := range postStartClusterStep {
		if err = s.Run(ctx, out, cfg); err != nil {
			return errors.Wrap(err, PostStartCluster)
		}
	}

	return nil
}

func (s StepPostStartCluster) Name() string {
	return PostStartCluster
}

func (s StepPostStartCluster) Description() string {
	return PostStartCluster
}

func (s StepPostStartCluster) Depends() []string {
	return nil
}

func (s StepPostStartCluster) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func postStartCluster(provider clouds.Name) ([]steps.Step, error) {
	switch provider {
	case clouds.AWS:
		return []steps.Step{
			steps.GetStep(amazon.StepCreateTags),
		}, nil
	case clouds.DigitalOcean:
		return []steps.Step{}, nil
	case clouds.Azure:
		return []steps.Step{}, nil
	case clouds.GCE:
		// TODO(stgleb): Add non-bootstrap master instances to instance groups
		return []steps.Step{}, nil
	}
	return nil, errors.Wrapf(fmt.Errorf("unknown provider: %s", provider), PreProvisionStep)
}
