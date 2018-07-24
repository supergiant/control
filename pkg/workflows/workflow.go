package workflows

import (
	"context"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type WorkFlow struct {
	Type   string       `json:"type"` // master or node
	Steps  []steps.Step `json:"steps"`
	Config steps.Config `json:"config"`
}

func New(steps []steps.Step, config steps.Config) *WorkFlow {
	return &WorkFlow{
		Steps:  steps,
		Config: config,
	}
}

func (w *WorkFlow) Run(ctx context.Context) error {
	for _, step := range w.Steps {
		if err := step.Run(ctx, w.Config); err != nil {
			return err
		}
	}

	return nil
}
