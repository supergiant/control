package client

import (
	"fmt"
	"reflect"

	"github.com/supergiant/supergiant/pkg/model"
)

// We don't use this directly, but instead compose other fake collections with
type CollectionInterface interface {
	List(model.List) error
	Create(model.Model) error
	Get(interface{}, model.Model) error
	GetWithIncludes(interface{}, model.Model, []string) error
	Update(interface{}, model.Model) error
	Delete(interface{}, model.Model) error
}

type Collection struct {
	client   *Client
	basePath string
}

func (c *Collection) List(list model.List) error {
	return c.client.request("GET", c.basePath, nil, list, list.QueryValues())
}

func (c *Collection) Get(id interface{}, item model.Model) error {
	return c.client.request("GET", c.memberPath(id), nil, item, nil)
}

func (c *Collection) GetWithIncludes(id interface{}, item model.Model, includes []string) error {
	queryValues := map[string][]string{"include": includes}
	return c.client.request("GET", c.memberPath(id), nil, item, queryValues)
}

func (c *Collection) Create(item model.Model) error {
	return c.client.request("POST", c.basePath, item, item, nil)
}

func (c *Collection) Update(id interface{}, item model.Model) error {
	return c.client.request("PATCH", c.memberPath(id), item, item, nil)
}

func (c *Collection) Delete(id interface{}, item model.Model) error {
	return c.client.request("DELETE", c.memberPath(id), nil, item, nil)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (c *Collection) memberPath(id interface{}) string {
	// id can be *int64 or string, so we must de-ref *int64
	indirectID := reflect.Indirect(reflect.ValueOf(id)).Interface()
	return fmt.Sprintf("%s/%v", c.basePath, indirectID)
}
