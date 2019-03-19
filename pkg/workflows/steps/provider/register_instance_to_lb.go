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
	RegisterInstanceStepName = "register_instance"
)

type RegisterInstanceToLoadBalancer struct {
}

func (s *RegisterInstanceToLoadBalancer) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	step, err := registerMachineToLoadBalancer(cfg.Provider)
	if err != nil {
		return err
	}

	return step.Run(ctx, out, cfg)
}

func (s *RegisterInstanceToLoadBalancer) Name() string {
	return CreateMachineStep
}

func (s *RegisterInstanceToLoadBalancer) Description() string {
	return CreateMachineStep
}

func (s *RegisterInstanceToLoadBalancer) Depends() []string {
	return nil
}

func (s *RegisterInstanceToLoadBalancer) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func registerMachineToLoadBalancer(provider clouds.Name) (steps.Step, error) {
	switch provider {
	case clouds.AWS:
		return steps.GetStep(amazon.RegisterInstanceStepName), nil
	// TODO(stgleb): rest of providers TBD
	}
	return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
}
