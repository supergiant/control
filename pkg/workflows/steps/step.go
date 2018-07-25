package steps

import (
	"context"
)

type Status string

const (
	StatusTodo    Status = "todo"
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Step interface {
	Run(ctx context.Context, config Config) error
	Name() string
	Description() string
}
