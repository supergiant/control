package client

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/supergiant/supergiant/pkg/model"
)

type Collection struct {
	client   *Client
	basePath string
}

func (c *Collection) List(items interface{}) error {
	return c.client.request("GET", c.basePath, nil, items, nil)
}

func (c *Collection) Get(id interface{}, item model.Model) error {
	return c.client.request("GET", c.memberPath(id), nil, item, nil)
}

func (c *Collection) GetWithIncludes(id *int64, item model.Model, includes []string) error {
	queryValues := map[string]string{"includes": strings.Join(includes, " ")}
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
