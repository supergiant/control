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
	DeleteMachineStep = "deleteMachine"
)

type StepDeleteMachine struct {
}

func (s StepDeleteMachine) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	step, err := deleteMachineStepFor(cfg.Provider)
	if err != nil {
		return err
	}

	return step.Run(ctx, out, cfg)
}

func (s StepDeleteMachine) Name() string {
	return DeleteMachineStep
}

func (s StepDeleteMachine) Description() string {
	return DeleteMachineStep
}

func (s StepDeleteMachine) Depends() []string {
	return nil
}

func (s StepDeleteMachine) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func deleteMachineStepFor(provider clouds.Name) (steps.Step, error) {
	switch provider {
	case clouds.AWS:
		return steps.GetStep(amazon.DeleteNodeStepName), nil
	case clouds.DigitalOcean:
		return steps.GetStep(digitalocean.DeleteMachineStepName), nil
	case clouds.GCE:
		return steps.GetStep(gce.DeleteNodeStepName), nil
	}
	return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
}
