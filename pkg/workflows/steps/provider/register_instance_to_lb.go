package provider

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/workflows/steps/openstack"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
)

const (
	RegisterInstanceStepName = "register_instance"
)

type RegisterInstanceToLoadBalancer struct {
}

func (s *RegisterInstanceToLoadBalancer) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	var step steps.Step

	switch cfg.Provider {
	case clouds.AWS:
		step = steps.GetStep(amazon.RegisterInstanceStepName)
	// TODO(stgleb): rest of providers TBD
	case clouds.DigitalOcean:
		// Load balancing in DO is made by tags
		return nil
	case clouds.GCE:
		return nil
	case clouds.Azure:
		return nil
	case clouds.OpenStack:
		step = steps.GetStep(openstack.RegisterInstancetoLBStepName)
	default:
		return errors.Wrapf(fmt.Errorf("unknown provider: %s", cfg.Provider), RegisterInstanceStepName)
	}

	return step.Run(ctx, out, cfg)
}

func (s *RegisterInstanceToLoadBalancer) Name() string {
	return RegisterInstanceStepName
}

func (s *RegisterInstanceToLoadBalancer) Description() string {
	return RegisterInstanceStepName
}

func (s *RegisterInstanceToLoadBalancer) Depends() []string {
	return nil
}

func (s *RegisterInstanceToLoadBalancer) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
