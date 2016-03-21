package model

import (
	"encoding/json"
	"fmt"
	"reflect"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

const (
	baseDir = "/supergiant"
)

type DB struct {
	kapi etcd.KeysAPI
}

func NewDB(endpoints []string) *DB {
	etcdClient, err := etcd.New(etcd.Config{Endpoints: endpoints})
	if err != nil {
		panic(err)
	}
	db := DB{etcd.NewKeysAPI(etcdClient)}

	// TODO
	db.CreateDir(baseDir)

	return &db
}

func fullKey(key string) string {
	return fmt.Sprintf("%s%s", baseDir, key)
}

func (db *DB) CompareAndSwap(key string, prevValue string, value string) (*etcd.Response, error) {
	return db.kapi.Set(context.Background(), fullKey(key), value, &etcd.SetOptions{PrevValue: prevValue})
}

func (db *DB) CreateInOrder(dir string, value string) (*etcd.Response, error) {
	return db.kapi.CreateInOrder(context.Background(), fullKey(dir), value, nil)
}

func (db *DB) GetInOrder(dir string) (*etcd.Response, error) {
	return db.kapi.Get(context.Background(), fullKey(dir), &etcd.GetOptions{Sort: true})
}

func (db *DB) CreateDir(key string) (*etcd.Response, error) {
	return db.kapi.Set(context.Background(), key, "", &etcd.SetOptions{Dir: true})
}

func (db *DB) create(key string, value string) (*etcd.Response, error) {
	return db.kapi.Create(context.Background(), fullKey(key), value)
}

func (db *DB) get(key string) (*etcd.Response, error) {
	return db.kapi.Get(context.Background(), fullKey(key), nil)
}

func (db *DB) update(key string, value string) (*etcd.Response, error) {
	return db.kapi.Update(context.Background(), fullKey(key), value)
}

func (db *DB) delete(key string) (*etcd.Response, error) {
	return db.kapi.Delete(context.Background(), fullKey(key), nil)
}

// TODO we should maybe rename these methods, or make a sub-class like ResourceDB

func (db *DB) List(r Resource, out Model) error {
	key := r.EtcdKey("")
	resp, err := db.get(key)
	if err != nil {
		return err
	}

	// The concrete value of an interface is a pair of 32-bit words, one pointing
	// to information about the type implementing the interface, and the other
	// pointing to the underlying data in the type.
	interfaceValue := reflect.ValueOf(out)

	// In this case, we expect out to have been passed as a pointer, so that
	// interfaceValue's real value is actually:
	//
	// [ pointer ] --> [ AppList type ]
	// [ pointer ] --> [ pointer to instance of AppList ]
	//
	// So, calling this will dereference the pointer, providing the underlying
	// value of AppList. It makes AppList addressable AND settable.
	// NOTE it will also panic if out was not passed as a pointer.
	modelValue := interfaceValue.Elem()

	// Items field on any ModelList should be a slice of the relevant Model.
	itemsField := modelValue.FieldByName("Items")
	if !itemsField.IsValid() {
		panic(fmt.Errorf("no Items field in %#v", out))
	}

	// Must first get the pointer of the slice with Addr(), so we can then call
	// Elem() to make it settable.
	itemsPtr := itemsField.Addr().Elem()

	for _, node := range resp.Node.Nodes {
		// Type() returns the underlying element type of the slice, and Elem()
		// allows us to utilize the type with reflect.New().
		itemType := itemsPtr.Type().Elem()

		// Interface() is called to convert the new item Value into an interface
		// (that we can unmarshal to. The interface{} is then cast to Model type.
		obj := reflect.New(itemType).Interface().(Model)
		unmarshalNodeInto(node, obj)

		// Get the Value of the unmarshalled object, and append it to the slice.
		newItem := reflect.ValueOf(obj).Elem()
		newItems := reflect.Append(itemsPtr, newItem)
		itemsPtr.Set(newItems)
	}
	return nil
}

func (db *DB) Create(r Resource, id string, m Model) error {
	key := r.EtcdKey(id)
	_, err := db.create(key, marshalModel(m))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(r Resource, id string, out Model) error {
	key := r.EtcdKey(id)
	resp, err := db.get(key)
	if err != nil {
		return err
	}
	unmarshalNodeInto(resp.Node, out)
	return nil
}

func (db *DB) Update(r Resource, id string, m Model) error {
	key := r.EtcdKey(id)
	_, err := db.update(key, marshalModel(m))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Delete(r Resource, id string) error {
	key := r.EtcdKey(id)
	_, err := db.delete(key)
	return err
}

func marshalModel(m Model) string {
	out, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func unmarshalNodeInto(node *etcd.Node, m Model) {
	if err := json.Unmarshal([]byte(node.Value), m); err != nil {
		panic(err)
	}
}
