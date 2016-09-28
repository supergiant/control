package core

import (
	"reflect"
	"sync"

	"github.com/imdario/mergo"
	"github.com/supergiant/supergiant/pkg/model"
)

type Collection struct {
	Core *Core
}

func (c *Collection) Create(m model.Model) error {
	return c.Core.DB.Create(m)
}

func (c *Collection) Get(id *int64, m model.Model) error {
	return c.Core.DB.First(m, *id)
}

func (c *Collection) GetWithIncludes(id *int64, m model.Model, includes []string) error {
	scope := c.Core.DB
	for _, include := range includes {
		scope = scope.Preload(include)
	}
	return scope.First(m, *id)
}

func (c *Collection) Update(id *int64, oldM model.Model, m model.Model) error {
	// model.ZeroReadonlyFields(m)

	// oldM := reflect.New(reflect.TypeOf(m)).Elem().Elem().Interface()

	if err := c.Core.DB.First(oldM, *id); err != nil {
		return err
	}

	// Merge old item attributes into the empty fields of the newItem
	if err := mergo.Merge(m, oldM); err != nil {
		return err
	}

	return c.Core.DB.Save(m)
}

func (c *Collection) Delete(id *int64, m model.Model) error { // Loaded so we can render out
	if err := c.Core.DB.First(m, *id); err != nil {
		return err
	}
	return c.Core.DB.Delete(m)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (c *Collection) inParallel(model interface{}, fn func(interface{}) error) (err error) {
	mv := reflect.ValueOf(model)
	count := mv.Len()

	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(idx int) {
			defer wg.Done()
			err = fn(mv.Index(idx).Interface())
		}(i)
	}

	wg.Wait()
	return
}
