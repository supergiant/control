package fake_core

import (
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

type Sessions struct {
	ClientFn func(id string) *client.Client
	ListFn   func() []*model.Session
	CreateFn func(*model.Session) error
	GetFn    func(id string, m *model.Session) error
	DeleteFn func(id string) error
}

func (f *Sessions) Client(id string) *client.Client {
	if f.ClientFn == nil {
		return nil
	}
	return f.ClientFn(id)
}

func (f *Sessions) List() []*model.Session {
	if f.ListFn == nil {
		return nil
	}
	return f.ListFn()
}

func (f *Sessions) Create(m *model.Session) error {
	if f.CreateFn == nil {
		return nil
	}
	return f.CreateFn(m)
}

func (f *Sessions) Get(id string, m *model.Session) error {
	if f.GetFn == nil {
		return nil
	}
	return f.GetFn(id, m)
}

func (f *Sessions) Delete(id string) error {
	if f.DeleteFn == nil {
		return nil
	}
	return f.DeleteFn(id)
}
