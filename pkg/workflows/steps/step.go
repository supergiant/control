package steps

import (
	"context"
	"io"
	"sync"
)

type Step interface {
	Run(context.Context, io.Writer, *Config) error
	Name() string
	Description() string
	Depends() []string
	Rollback(context.Context, io.Writer, *Config) error
}

var (
	m       sync.RWMutex
	stepMap map[string]Step
)

func init() {
	stepMap = make(map[string]Step)
}

func RegisterStep(stepName string, step Step) {
	m.Lock()
	defer m.Unlock()
	stepMap[stepName] = step
}

func GetStep(stepName string) Step {
	m.RLock()
	defer m.RUnlock()
	return stepMap[stepName]
}
