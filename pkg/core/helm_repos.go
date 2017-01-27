package core

import "github.com/supergiant/supergiant/pkg/model"

type HelmRepos struct {
	Collection
}

func (c *HelmRepos) Delete(id *int64, m *model.HelmRepo) ActionInterface {
	return &Action{
		Status: &model.ActionStatus{
			Description: "deleting",
			MaxRetries:  5,
		},
		Core:  c.Core,
		Scope: c.Core.DB.Preload("Charts"),
		Model: m,
		ID:    id,
		Fn: func(a *Action) error {
			// Delete chart records directly
			for _, chart := range m.Charts {
				if err := c.Core.DB.Delete(chart); err != nil {
					return err
				}
			}
			return c.Collection.Delete(id, m)
		},
	}
}
