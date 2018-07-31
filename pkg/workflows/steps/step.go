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
	Run(context.Context, io.Writer, Config) error
	Name() string
	Description() string
}

var (
	stepMap map[string]Step
)

func init() {
	stepMap = make(map[string]Step)
}

func RegisterStep(stepName string, step Step) {
	stepMap[stepName] = step
}

func GetStep(stepName string) Step {
	return stepMap[stepName]
}
