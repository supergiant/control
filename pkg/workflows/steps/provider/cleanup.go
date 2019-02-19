package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
)

const (
	CleanUpStep = "cleanUp"
)

type StepCleanUp struct {
}

func (s StepCleanUp) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	steps, err := cleanUpStepsFor(cfg.Provider)
	if err != nil {
		return errors.Wrap(err, CleanUpStep)
	}
	for _, s := range steps {
		if err = s.Run(ctx, out, cfg); err != nil {
			return errors.Wrap(err, CleanUpStep)
		}
	}

	return nil
}

func (s StepCleanUp) Name() string {
	return CleanUpStep
}

func (s StepCleanUp) Description() string {
	return CleanUpStep
}

func (s StepCleanUp) Depends() []string {
	return nil
}

func (s StepCleanUp) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func cleanUpStepsFor(provider clouds.Name) ([]steps.Step, error) {
	// TODO: use provider interface
	switch provider {
	case clouds.AWS:
		return []steps.Step{
			steps.GetStep(amazon.DeleteClusterMachinesStepName),
			steps.GetStep(amazon.DeleteSecurityGroupsStepName),
			steps.GetStep(amazon.DisassociateRouteTableStepName),
			steps.GetStep(amazon.DeleteSubnetsStepName),
			steps.GetStep(amazon.DeleteRouteTableStepName),
			steps.GetStep(amazon.DeleteInternetGatewayStepName),
			steps.GetStep(amazon.DeleteKeyPairStepName),
			steps.GetStep(amazon.DeleteVPCStepName),
		}, nil
	case clouds.DigitalOcean:
		return []steps.Step{
			steps.GetStep(digitalocean.DeleteMachineStepName),
			steps.GetStep(digitalocean.DeleteDeleteKeysStepName),
		}, nil
	case clouds.GCE:
		return []steps.Step{
			steps.GetStep(gce.DeleteNodeStepName),
		}, nil
	case clouds.Azure:
		return []steps.Step{
			//TODO DELETION
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
}
