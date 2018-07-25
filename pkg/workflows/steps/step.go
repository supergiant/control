package steps

import (
	"context"
	"io"
)

type Status string

const (
	StatusTodo    Status = "todo"
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Step interface {
	Run(context.Context,io.Writer,Config) error
	Name() string
	Description() string
}
