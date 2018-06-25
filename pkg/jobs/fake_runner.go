package jobs

import (
	"errors"
	"io"
	"strings"


	"github.com/supergiant/supergiant/pkg/runner"
)

type FakeRunner struct {
	ErrMsg string
}

func (f *FakeRunner) Run(command *runner.Command) error {
	if len(f.ErrMsg) > 0 {
		return errors.New(f.ErrMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}
