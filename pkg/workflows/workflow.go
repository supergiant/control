package workflows

import (
	"context"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type WorkFlow struct {
	s      []steps.Step
	config Config
}

func New(steps []steps.Step, config Config) *WorkFlow {
	return &WorkFlow{
		s:      steps,
		config: config,
	}
}

func (w *WorkFlow) Run(ctx context.Context) error {
	for _, step := range w.s {
		if err := step.Run(ctx, w.config); err != nil {
			return err
		}
	}

	return nil
}
