package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/azure"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
)

const (
	PreProvisionStep = "preProvision"
)

type StepPreProvision struct {
}

func (s StepPreProvision) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	preProvisionSteps, err := prepProvisionStepFor(cfg.Provider)
	if err != nil {
		return errors.Wrap(err, PreProvisionStep)
	}
	for _, s := range preProvisionSteps {
		if err = s.Run(ctx, out, cfg); err != nil {
			return errors.Wrap(err, PreProvisionStep)
		}
	}

	return nil
}

func (s StepPreProvision) Name() string {
	return PreProvisionStep
}

func (s StepPreProvision) Description() string {
	return PreProvisionStep
}

func (s StepPreProvision) Depends() []string {
	return nil
}

func (s StepPreProvision) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func prepProvisionStepFor(provider clouds.Name) ([]steps.Step, error) {
	// TODO: use provider interface
	switch provider {
	case clouds.AWS:
		return []steps.Step{
			steps.GetStep(amazon.StepFindAMI),
			steps.GetStep(amazon.StepCreateVPC),
			steps.GetStep(amazon.StepCreateSecurityGroups),
			steps.GetStep(amazon.StepNameCreateInstanceProfiles),
			steps.GetStep(amazon.ImportKeyPairStepName),
			steps.GetStep(amazon.StepCreateInternetGateway),
			steps.GetStep(amazon.StepCreateSubnets),
			steps.GetStep(amazon.StepCreateRouteTable),
			steps.GetStep(amazon.StepAssociateRouteTable),
			steps.GetStep(amazon.StepCreateLoadBalancer),
		}, nil
	case clouds.DigitalOcean:
		return []steps.Step{
			// TODO(stgleb): Apply security stuff here
			steps.GetStep(digitalocean.CreateLoadBalancerStepName),
		}, nil
	case clouds.Azure:
		return []steps.Step{
			steps.GetStep(azure.GetAuthorizerStepName),
			steps.GetStep(azure.CreateGroupStepName),
			steps.GetStep(azure.CreateLBStepName),
			steps.GetStep(azure.CreateVNetAndSubnetsStepName),
			steps.GetStep(azure.CreateSecurityGroupStepName),
		}, nil
	case clouds.GCE:
		return []steps.Step{
			steps.GetStep(gce.CreateIPAddressStepName),
			steps.GetStep(gce.CreateTargetPullStepName),
			steps.GetStep(gce.CreateForwardingRulesStepName),
		}, nil
	}
	return nil, errors.Wrapf(fmt.Errorf("unknown provider: %s", provider), PreProvisionStep)
}
