package core

import (
	"reflect"
	"runtime/debug"
	"time"
)

type Service interface {
	Perform() error
}

type RecurringService struct {
	core *Core

	// can provide function directly (which takes priority), or Service interface
	fn       func() error
	service  Service
	tag      string
	interval time.Duration
}

func (s *RecurringService) Run() {
	s.tick() // we want to run Service once immediately before waiting interval
	for range time.NewTicker(s.interval).C {
		s.tick()
	}
}

func (s *RecurringService) tick() {
	defer s.recover()
	if err := s.perform(); err != nil {

		s.core.Log.Error("Error in RecurringService ["+s.tag+"] : ", err)
	}
}

func (s *RecurringService) perform() error {
	if s.fn != nil {
		return s.fn()
	}
	return s.service.Perform()
}

func (s *RecurringService) name() string {
	// TODO no name if no service provided...
	if s.service == nil {
		return ""
	}
	return reflect.TypeOf(s.service).Elem().String()
}

func (s *RecurringService) recover() {
	if r := recover(); r != nil {
		s.core.Log.Error("Recovered in RecurringService "+s.name()+": ", r)
		s.core.Log.Debug(string(debug.Stack()))
	}
}
