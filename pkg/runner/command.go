package runner

import (
	"context"
	"github.com/pkg/errors"
	"io"
)

var (
	ErrNilContext = errors.New("nil context")
	ErrNilWriter  = errors.New("writer is nil")
)

// Command is an action that can be run and cancelled on different environments ssh, shell, docker etc.
type Command struct {
	Ctx context.Context

	Script string

	Out io.Writer
	Err io.Writer
}

//  TODO(stgleb): Use single io.Writer for gathering command output
func NewCommand(ctx context.Context, script string, out, err io.Writer) (*Command, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if out == nil || err == nil {
		return nil, ErrNilWriter
	}

	return &Command{
		ctx,
		script,
		out,
		err,
	}, nil
}
