package storageclass

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/supergiant/control/pkg/runner"
)

//TODO cleanup all fakerunners and move them to separate package
type fakeRunner struct {
	errMsg string
}

func (f *fakeRunner) Run(command *runner.Command) error {
	if len(f.errMsg) > 0 {
		return errors.New(f.errMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}

func TestStep_Run(t *testing.T) {

}
