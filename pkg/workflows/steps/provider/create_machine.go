package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/azure"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
)

const (
	CreateMachineStep = "createMachine"
)

type StepCreateMachine struct {
}

func (s StepCreateMachine) Run(ctx context.Context, out io.Writer, cfg *steps.Config) error {
	if cfg == nil {
		return errors.New("invalid config")
	}

	step, err := createMachineStepFor(cfg.Provider)
	if err != nil {
		return err
	}
	// TODO: check it on steps.GetStep()
	if step == nil {
		return errors.Wrap(sgerrors.ErrRawError, "createMachine step not found")
	}

	return step.Run(ctx, out, cfg)
}

func (s StepCreateMachine) Name() string {
	return CreateMachineStep
}

func (s StepCreateMachine) Description() string {
	return CreateMachineStep
}

func (s StepCreateMachine) Depends() []string {
	return nil
}

func (s StepCreateMachine) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func createMachineStepFor(provider clouds.Name) (steps.Step, error) {
	switch provider {
	case clouds.AWS:
		return steps.GetStep(amazon.StepNameCreateEC2Instance), nil
	case clouds.DigitalOcean:
		return steps.GetStep(digitalocean.CreateMachineStepName), nil
	case clouds.GCE:
		return steps.GetStep(gce.CreateInstanceStepName), nil
	case clouds.Azure:
		return steps.GetStep(azure.CreateVMStepName), nil
	}
	return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
}
