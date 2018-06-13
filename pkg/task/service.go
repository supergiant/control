package task

import (
	"context"

	"github.com/pkg/errors"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Service orchestrates everything related to running tasks
type Service struct {
	srv *machinery.Server
}

func NewService(cnf *config.Config) (*Service, error) {
	srv, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, errors.Wrap(err, "task service")
	}
	return &Service{
		srv: srv,
	}, nil
}

func (s *Service) RegisterTask(taskName string, fn interface{}) error {
	err := s.srv.RegisterTask(taskName, fn)
	return errors.Wrap(err, "task service register")
}

func (s *Service) SendTask(ctx context.Context, task *tasks.Signature) (*backends.AsyncResult, error) {
	return s.srv.SendTaskWithContext(ctx, task)
}

func (s *Service) SendChain(ctx context.Context, chain *tasks.Chain) (*backends.ChainAsyncResult, error) {
	return s.srv.SendChainWithContext(ctx, chain)
}
