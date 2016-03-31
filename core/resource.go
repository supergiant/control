package core

import (
	"fmt"
	"reflect"

	"github.com/supergiant/supergiant/types"
)

// Collection is an interface for defining behavior of a collection of
// Resources.
type Collection interface {
	EtcdKey(id types.ID) string

	// InitializeResource is called when unmarshalling objects from etcd.
	// Primarily, it sets a pointer to the Collection on the Resource.
	InitializeResource(r Resource)
}

// Resource is an interface used mainly for generalized marshalling purposes for
// resource types.
type Resource interface {
	// PersistableObject returns the portion (or whole) of the Resource meant to
	// be marshalled and stored in the DB. This allows for dynamic "virtual"
	// values to be loaded and displayed in the API without needing to be stored.
	PersistableObject() interface{}
}

// OrderedResource is similar to Resource, but provides a SetID() method to
// set an auto-generated ID from etcd on the Resource.
type OrderedResource interface {
	PersistableObject() interface{}
	SetID(id types.ID)
}

// TODO should maybe move this to util or helper file
// getItemsPtrAndItemType takes a Resource, which must be of the List type, and
// returns a pointer to the Items slice of the List and the underlying item type
// of the slice.
func getItemsPtrAndItemType(r interface{}) (reflect.Value, reflect.Type) {
	// The concrete value of an interface is a pair of 32-bit words, one pointing
	// to information about the type implementing the interface, and the other
	// pointing to the underlying data in the type.
	interfaceValue := reflect.ValueOf(r)

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
		panic(fmt.Errorf("no Items field in %#v", r))
	}

	// Items field is a slice here... (not a pointer)

	// Must first get the pointer of the slice with Addr(), so we can then call
	// Elem() to make it settable.
	itemsPtr := itemsField //.Addr() //.Interface()
	// Type() returns the underlying element type of the slice, and Elem()
	// allows us to utilize the type with reflect.New().
	itemType := itemsPtr.Type().Elem().Elem()

	// fmt.Println(fmt.Sprintf("m: %#v", m))
	// fmt.Println(fmt.Sprintf("interfaceValue: %#v", interfaceValue))
	// fmt.Println(fmt.Sprintf("modelValue: %#v", modelValue))
	// fmt.Println(fmt.Sprintf("itemsField: %#v", itemsField))
	// fmt.Println(fmt.Sprintf("itemsPtr: %#v", itemsPtr))
	// fmt.Println(fmt.Sprintf("itemType: %#v", itemType))

	return itemsPtr, itemType
}

// TODO all the following stuff should be better documented.

// // used to initialize Meta object primarily
// func initResource(r Resource) {
// 	meta := reflect.ValueOf(types.NewMeta()).Elem()
// 	getFieldValue(r.PersistableObject(), "Meta").Set(meta)
// }

func getFieldValue(r Resource, f string) reflect.Value {
	field := reflect.ValueOf(r).Elem().FieldByName(f)
	if !field.IsValid() {
		panic(fmt.Errorf("No %s field in %#v", f, r))
	}
	return field
}

func newTimestampValue() reflect.Value {
	return reflect.ValueOf(types.NewTimestamp()) //.Elem()
}

func setCreatedTimestamp(r Resource) {
	getFieldValue(r, "Created").Set(newTimestampValue())
}

func setUpdatedTimestamp(r Resource) {
	getFieldValue(r, "Updated").Set(newTimestampValue())
}
