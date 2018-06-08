package command

import (
	"context"
	"strings"
)

// Command is an action that can be run and cancelled on different environments ssh, shell, docker etc.
type Command struct {
	Command string
	Args    []string
	Ctx     context.Context
}

// FullCommand gets single string for command
func (c *Command) FullCommand() string {
	return c.Command + strings.Join(c.Args, " ")
}
