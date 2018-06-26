package testutils

import (
	"io"
	"strings"

	"github.com/supergiant/supergiant/pkg/runner"
)

type FakeRunner struct {
	Err error
}

func (f *FakeRunner) Run(command *runner.Command) error {
	if f.Err != nil {
		return f.Err
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}
