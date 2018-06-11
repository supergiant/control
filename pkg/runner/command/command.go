package command

import (
	"context"
	"io"
	"strings"
)

// Command is an action that can be run and cancelled on different environments ssh, shell, docker etc.
type Command struct {
	Ctx context.Context

	Command string
	Args    []string

	Out io.Writer
	Err io.Writer
}

func NewCommand(ctx context.Context, cmd string, args []string, out, err io.Writer) *Command {
	return &Command{
		ctx,
		cmd,
		args,
		out,
		err,
	}
}

// FullCommand gets single string for command
func (c *Command) FullCommand() string {
	return c.Command + strings.Join(c.Args, " ")
}
