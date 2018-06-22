package runner

import (
	"context"
	"io"
)

// Command is an action that can be run and cancelled on different environments ssh, shell, docker etc.
type Command struct {
	Ctx context.Context

	Script string

	Out io.Writer
	Err io.Writer
}

func NewCommand(ctx context.Context, script string, out, err io.Writer) *Command {
	return &Command{
		ctx,
		script,
		out,
		err,
	}
}
