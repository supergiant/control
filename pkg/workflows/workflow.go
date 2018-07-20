package workflows

import "github.com/supergiant/supergiant/pkg/workflows/steps"

type WorkFlow struct {
	s []*steps.Step
}


func New(steps []*steps.Step) *WorkFlow {
	return &WorkFlow{
		s: steps,
	}
}


func (w *WorkFlow) Run() error {
	for _, step := range w.s {
		if err := step.Run(); err != nil {
			return err
		}
	}

	return nil
}