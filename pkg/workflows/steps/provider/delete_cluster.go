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
	DeleteClusterStepName = "deleteCluster"
)

type DeleteCluster struct {
}

func (s DeleteCluster) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	steps, err := cleanUpStepsFor(cfg.Provider)
	if err != nil {
		return errors.Wrap(err, DeleteClusterStepName)
	}
	for _, s := range steps {
		if err = s.Run(ctx, out, cfg); err != nil {
			return errors.Wrap(err, DeleteClusterStepName)
		}
	}

	return nil
}

func (s DeleteCluster) Name() string {
	return DeleteClusterStepName
}

func (s DeleteCluster) Description() string {
	return DeleteClusterStepName
}

func (s DeleteCluster) Depends() []string {
	return nil
}

func (s DeleteCluster) Rollback(context.Context, io.Writer, *steps.Config) error {
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
			steps.GetStep(digitalocean.DeleteClusterMachines),
			steps.GetStep(digitalocean.DeleteDeleteKeysStepName),
		}, nil
	case clouds.GCE:
		return []steps.Step{
			steps.GetStep(gce.DeleteNodeStepName),
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
}
