package command

import (
	"context"
	"strings"
)

type Command struct {
	Command string
	Args    []string
	Ctx     context.Context
}

func (c *Command) FullCommand() string {
	return c.Command + strings.Join(c.Args, " ")
}
