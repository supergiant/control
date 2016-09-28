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
	core     *Core
	service  Service
	interval time.Duration
}

func (s *RecurringService) Run() {
	for _ = range time.NewTicker(s.interval).C {
		s.tick()
	}
}

func (s *RecurringService) tick() {
	defer s.recover()
	if err := s.service.Perform(); err != nil {
		s.core.Log.Error("Error in RecurringService "+s.name()+": ", err)
	}
}

func (s *RecurringService) name() string {
	return reflect.TypeOf(s.service).Elem().String()
}

func (s *RecurringService) recover() {
	if r := recover(); r != nil {
		s.core.Log.Error("Recovered in RecurringService "+s.name()+": ", r)
		s.core.Log.Debug(string(debug.Stack()))
	}
}
