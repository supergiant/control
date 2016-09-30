package fake_core

import (
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type DB struct {
	CreateFn  func(model.Model) error
	SaveFn    func(model.Model) error
	FindFn    func(out interface{}, where ...interface{}) error
	FirstFn   func(out interface{}, where ...interface{}) error
	DeleteFn  func(m model.Model) error
	PreloadFn func(column string, conditions ...interface{}) core.DBInterface
	WhereFn   func(query interface{}, args ...interface{}) core.DBInterface
	LimitFn   func(limit interface{}) core.DBInterface
	OffsetFn  func(offset interface{}) core.DBInterface
	ModelFn   func(value interface{}) core.DBInterface
	UpdateFn  func(attrs ...interface{}) error
	CountFn   func(interface{}) error
}

func (db *DB) Create(m model.Model) error {
	if db.CreateFn == nil {
		return nil
	}
	return db.CreateFn(m)
}

func (db *DB) Save(m model.Model) error {
	if db.SaveFn == nil {
		return nil
	}
	return db.SaveFn(m)
}

func (db *DB) Find(out interface{}, where ...interface{}) error {
	if db.FindFn == nil {
		return nil
	}
	return db.FindFn(out, where...)
}

func (db *DB) First(out interface{}, where ...interface{}) error {
	if db.FirstFn == nil {
		return nil
	}
	return db.FirstFn(out, where...)
}

func (db *DB) Delete(m model.Model) error {
	if db.DeleteFn == nil {
		return nil
	}
	return db.DeleteFn(m)
}

func (db *DB) Preload(column string, conditions ...interface{}) core.DBInterface {
	if db.PreloadFn == nil {
		return nil
	}
	return db.PreloadFn(column, conditions...)
}

func (db *DB) Where(query interface{}, args ...interface{}) core.DBInterface {
	if db.WhereFn == nil {
		return nil
	}
	return db.WhereFn(query, args...)
}

func (db *DB) Limit(limit interface{}) core.DBInterface {
	if db.LimitFn == nil {
		return nil
	}
	return db.LimitFn(limit)
}

func (db *DB) Offset(offset interface{}) core.DBInterface {
	if db.OffsetFn == nil {
		return nil
	}
	return db.OffsetFn(offset)
}

func (db *DB) Model(value interface{}) core.DBInterface {
	if db.ModelFn == nil {
		return nil
	}
	return db.ModelFn(value)
}

func (db *DB) Update(attrs ...interface{}) error {
	if db.UpdateFn == nil {
		return nil
	}
	return db.UpdateFn(attrs...)
}

func (db *DB) Count(value interface{}) error {
	if db.CountFn == nil {
		return nil
	}
	return db.CountFn(value)
}
