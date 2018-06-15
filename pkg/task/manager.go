package task

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/pkg/errors"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v1/backends"
)

// Manager is a simple wrapper around task engine, currently it is machinery
type Manager struct {
	srv *machinery.Server
}

func NewManager(cnf *config.Config) (*Manager, error) {
	srv, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, errors.Wrap(err, "task manager")
	}

	return &Manager{
		srv: srv,
	}, nil
}

func (m *Manager) RegisterTask(taskName string, fn interface{}) error {
	err := m.srv.RegisterTask(taskName, fn)
	return errors.Wrap(err, "task manager register")
}

func (m *Manager) SendTask(task *tasks.Signature) (*backends.AsyncResult, error) {
	return m.srv.SendTask(task)
}
