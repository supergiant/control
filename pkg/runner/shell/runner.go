package shell

import (
	"os/exec"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner/command"
)

// Runner is implementation of runner interface for shell
type Runner struct{}

// Run command on shell
func (r *Runner) Run(c command.Command) error {
	cmd := exec.CommandContext(c.Ctx, c.Command, c.Args...)

	// Start a process
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to start process: %s")
	}

	return err
}
