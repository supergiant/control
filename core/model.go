package core

import (
	"fmt"
	"reflect"
)

// TODO
type Resource interface {
	EtcdKey(id string) string
	InitializeModel(m Model)
}

type Model interface {
}

// NOTE this is used for ordered etcd keys, which have sequential IDs returned by etcd
type OrderedModel interface {
	SetID(id string)
}

// TODO should maybe move this to util or helper file
func GetItemsPtrAndItemType(m Model) (reflect.Value, reflect.Type) {
	// The concrete value of an interface is a pair of 32-bit words, one pointing
	// to information about the type implementing the interface, and the other
	// pointing to the underlying data in the type.
	interfaceValue := reflect.ValueOf(m)

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
		panic(fmt.Errorf("no Items field in %#v", m))
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
