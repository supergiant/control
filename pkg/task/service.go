package task

import (
	"context"

	"github.com/pkg/errors"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Service orchestrates everything related to running tasks, currently it delegated running tasks to the machinery
type Service struct {
	srv *machinery.Server
}

// NewService creates the task service instance.
func NewService(cnf *config.Config) (*Service, error) {
	srv, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, errors.Wrap(err, "task service")
	}
	return &Service{
		srv: srv,
	}, nil
}

// RegisterTaskFunction in order to invoke task handler function it should be registered in machinery
// the function could have an arbitrary number of arguments however it is mandatory to return an error as the last one.
func (s *Service) RegisterTaskFunction(taskName string, fn interface{}) error {
	err := s.srv.RegisterTask(taskName, fn)
	return errors.Wrap(err, "task service register")
}

// Send individual task to be executed.
func (s *Service) Send(ctx context.Context, task *tasks.Signature) (*result.AsyncResult, error) {
	return s.srv.SendTaskWithContext(ctx, task)
}

// SendChain executed a chain of tasks.
func (s *Service) SendChain(ctx context.Context, chain *tasks.Chain) (*result.ChainAsyncResult, error) {
	return s.srv.SendChainWithContext(ctx, chain)
}
