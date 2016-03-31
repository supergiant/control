package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/supergiant/supergiant/types"

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

func (db *DB) compareAndSwap(key string, prevValue string, value string) (*etcd.Response, error) {
	return db.kapi.Set(context.Background(), fullKey(key), value, &etcd.SetOptions{PrevValue: prevValue})
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

func (db *DB) createInOrder(key string, value string) (*etcd.Response, error) {
	return db.kapi.CreateInOrder(context.Background(), fullKey(key), value, nil)
}

func (db *DB) getInOrder(key string) (*etcd.Response, error) {
	return db.kapi.Get(context.Background(), fullKey(key), &etcd.GetOptions{Sort: true})
}

// TODO
func (db *DB) CreateDir(key string) (*etcd.Response, error) {
	return db.kapi.Set(context.Background(), key, "", &etcd.SetOptions{Dir: true})
}

//----------------------------------------------------------------------------//
//----------------------------------------------------------------------------//
//----------------- start of resource-specific DB operations -----------------//
//----------------------------------------------------------------------------//
//----------------------------------------------------------------------------//

func decodeList(r Collection, resp *etcd.Response, out interface{}) error {
	itemsPtr, itemType := getItemsPtrAndItemType(out)

	// NOTE this this is not working, because we're not getting the right
	// underlying element *type (with pointer).
	//
	// emptyItems := reflect.MakeSlice(reflect.SliceOf(itemType), 0, 0)
	// itemsPtr.Set(emptyItems)

	for _, node := range resp.Node.Nodes {
		// Interface() is called to convert the new item Value into an interface
		// (that we can unmarshal to. The interface{} is then cast to ResourceList type.
		obj := reflect.New(itemType).Interface().(Resource)
		unmarshalNodeInto(r, node, obj)

		// Get the Value of the unmarshalled object, and append it to the slice.
		newItem := reflect.ValueOf(obj).Elem().Addr()
		newItems := reflect.Append(itemsPtr, newItem)
		itemsPtr.Set(newItems)
	}
	return nil
}

// TODO feel like there's a DRYer or cleaner way to do this
func decodeOrderedList(r Collection, resp *etcd.Response, out interface{}) error { /// ------------------- just changed to Resource from OrderedResource
	itemsPtr, itemType := getItemsPtrAndItemType(out)
	for _, node := range resp.Node.Nodes {
		// Interface() is called to convert the new item Value into an interface
		// (that we can unmarshal to. The interface{} is then cast to Resource type.

		obj := reflect.New(itemType).Interface().(OrderedResource)

		unmarshalNodeInto(r, node, obj)

		obj.SetID(lastKeySegment(node.Key))

		// Get the Value of the unmarshalled object, and append it to the slice.
		newItem := reflect.ValueOf(obj).Elem().Addr()
		newItems := reflect.Append(itemsPtr, newItem)
		itemsPtr.Set(newItems)
	}
	return nil
}

func (db *DB) List(r Collection, out interface{}) error {
	key := r.EtcdKey(nil)
	resp, err := db.get(key)
	if err != nil {
		return err
	}
	return decodeList(r, resp, out)
}

func (db *DB) Create(r Collection, id types.ID, m Resource) error {
	key := r.EtcdKey(id)

	setCreatedTimestamp(m)

	_, err := db.create(key, marshalResource(m))
	if err != nil {
		return err
	}
	// NOTE we do this here because we call it when unmarshalling normally, and
	// we don't need to do that here.
	r.InitializeResource(m)
	return nil
}

func (db *DB) Get(r Collection, id types.ID, out Resource) error {
	key := r.EtcdKey(id)
	resp, err := db.get(key)
	if err != nil {
		return err
	}
	unmarshalNodeInto(r, resp.Node, out)
	return nil
}

func (db *DB) Update(r Collection, id types.ID, m Resource) error {
	key := r.EtcdKey(id)

	setUpdatedTimestamp(m)

	_, err := db.update(key, marshalResource(m))
	if err != nil {
		return err
	}
	r.InitializeResource(m)
	return nil
}

func (db *DB) Delete(r Collection, id types.ID) error {
	key := r.EtcdKey(id)
	_, err := db.delete(key)
	return err
}

//------------------------------------------------------------------------------
func (db *DB) ListInOrder(r Collection, out interface{}) error {
	key := r.EtcdKey(nil)
	resp, err := db.getInOrder(key)
	if err != nil {
		return err
	}
	return decodeOrderedList(r, resp, out)
}

func (db *DB) CreateInOrder(r Collection, m OrderedResource) error {
	key := r.EtcdKey(nil) // ID is generated by etcd
	resp, err := db.createInOrder(key, marshalResource(m))
	if err != nil {
		return err
	}

	// We must set ID value on model, since it is auto-generated by etcd
	m.SetID(lastKeySegment(resp.Node.Key))

	return nil
}

//------------------------------------------------------------------------------

func (db *DB) CompareAndSwap(r Collection, id types.ID, old Resource, new Resource) error {
	key := r.EtcdKey(id)
	_, err := db.compareAndSwap(key, marshalResource(old), marshalResource(new))
	return err
}

func marshalResource(m Resource) string {
	out, err := json.Marshal(m.PersistableObject())
	if err != nil {
		panic(err)
	}
	return string(out)
}

func unmarshalNodeInto(r Collection, node *etcd.Node, m Resource) {
	if err := json.Unmarshal([]byte(node.Value), m); err != nil {
		panic(err.Error() + " " + node.Key + " " + node.Value)
	}
	r.InitializeResource(m)
}

// CreateInOrder stuff...
// Was going to use "base" as a word here, like with file names. But it seems
// entirely weird to me that people inventing filesystems looked at:
//
// /home/dir/filename.txt
//
// and decided that "filename.txt" was the "base name".
func lastKeySegment(key string) types.ID {
	strs := strings.Split(key, "/")
	segment := strs[len(strs)-1]
	return &segment
}
