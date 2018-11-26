package testutils

import (
	"io"
	"strings"

	"github.com/supergiant/control/pkg/runner"
)

type MockRunner struct {
	Err error
}

func (m *MockRunner) Run(command *runner.Command) error {
	if m.Err != nil {
		return m.Err
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}
