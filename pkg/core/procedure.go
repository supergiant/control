package core

import (
	"github.com/supergiant/supergiant/pkg/model"
)

type Procedure struct {
	Core   *Core
	Name   string
	Model  model.Model
	Action *Action
	steps  []*Step
}

type Step struct {
	desc string
	fn   func() error
}

func (p *Procedure) AddStep(desc string, fn func() error) {
	p.steps = append(p.steps, &Step{desc, fn})
}

func (p *Procedure) Run() error {
	for n, step := range p.steps {

		if p.Action.Status.StepsCompleted > n {
			continue
		}

		p.Core.Log.Infof("Running step of %s procedure: %s", p.Name, step.desc)
		if err := step.fn(); err != nil {
			return err
		}

		// If there is no error, it means we've moved past whatever error there may
		// have been from a previous try of this step.
		p.Action.Status.Error = ""
		p.Action.Status.StepsCompleted = n + 1

		// We save here so that attributes changed on model during fn() are saved
		if err := p.Core.DB.Save(p.Model); err != nil {
			return err
		}
	}
	return nil
}
