package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"context"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// func UniqueStrings(in []string) (out []string) {
// 	tab := make(map[string]struct{})
// 	for _, str := range in {
// 		if _, ok := tab[str]; !ok {
// 			tab[str] = struct{}{}
// 			out = append(out, str)
// 		}
// 	}
// 	return out
// }

// WaitFor event using context
func WaitFor(desc string, ctx context.Context, period time.Duration, fn func() (bool, error)) error {
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return errors.Wrap(ctx.Err(), desc)
			}
		case <-ticker.C:
			if done, err := fn(); done {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
}

type Stringer interface {
	String() string
}

// BuildSchema
func BuildSchema(model interface{}) map[string]interface{} {
	tempSchema := make(map[string]interface{})
	RecurseSchema(tempSchema, model)
	finalSchema := make(map[string]interface{})
	finalSchema["properties"] = tempSchema
	return finalSchema
}

// RecurseSchema
func RecurseSchema(schema map[string]interface{}, obj interface{}) {
	// objs = variadic interface type to take in the root model, and sub models
	// for recursion support
	fmt.Println("recurseSchema Head:", obj, reflect.TypeOf(obj).Kind())
	var objType interface{}
	if reflect.TypeOf(obj) == nil {
		objType = reflect.String
	} else {
		objType = reflect.TypeOf(obj).Kind()
	}

	switch objType {
	case reflect.Map:
		objMap := obj.(map[string]interface{})
		if len(objMap) == 0 {
			schema["type"] = "string"
		}
		for k, v := range objMap {
			schemaItem := make(map[string]interface{})
			if (v != nil) && (reflect.TypeOf(v).Kind() == reflect.Map) {
				RecurseSchema(schemaItem, v)
				schema[k] = make(map[string]interface{})
				schemaRoot := schema[k].(map[string]interface{})
				schemaRoot["type"] = "object"
				schemaRoot["properties"] = schemaItem
			} else {
				RecurseSchema(schemaItem, v)
				schema[k] = schemaItem
			}
		}
	case reflect.Struct:
		objStruct := structs.New(obj)
		var itemName string
		var itemValue interface{}
		for _, f := range objStruct.Fields() {
			schemaItem := make(map[string]interface{})
			if f.IsExported() {
				if len(f.Tag("json")) > 0 {
					itemName = f.Tag("json")
					if itemName == "-" {
						continue
					}
					itemName = strings.Replace(itemName, ",omitempty", "", 1)
				} else {
					itemName = f.Name()
				}

				itemValue = f.Value()
				if strStruct, ok := f.Value().(string); ok {
					itemValue = strStruct
				} else if strStruct, ok := f.Value().(Stringer); ok {
					itemValue = strStruct.String()
				}

				fmt.Println("Calling recurse with:", itemName, itemValue, reflect.TypeOf(itemValue).Kind())
				typeSchema := reflect.TypeOf(itemValue).Kind()
				switch typeSchema {
				case reflect.Struct:
					schemaItem["type"] = "object"
					tempschemaItem := make(map[string]interface{})
					RecurseSchema(tempschemaItem, itemValue)
					schemaItem["properties"] = tempschemaItem
				case reflect.Ptr:
					typeSchema = reflect.TypeOf(itemValue).Elem().Kind()
					switch typeSchema {
					case reflect.Struct:
						schemaItem["type"] = "object"
						tempschemaItem := make(map[string]interface{})
						RecurseSchema(tempschemaItem, itemValue)
						schemaItem["properties"] = tempschemaItem
					default:
						RecurseSchema(schemaItem, itemValue)
					}
				default:
					RecurseSchema(schemaItem, itemValue)
				}

				schema[itemName] = schemaItem

			}

		}

	case reflect.String:
		// there is an unsafe assumption here that a string is a single value pair
		schema["type"] = "string"
		if (obj != nil) && (len(obj.(string)) > 0) {
			schema["default"] = obj
		}

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		var intString string
		intType := reflect.TypeOf(obj).Kind()
		switch intType {
		case reflect.Int64:
			intString = strconv.FormatInt(obj.(int64), 10)
		case reflect.Int32:
			intString = strconv.FormatInt(int64(obj.(int32)), 10)
		case reflect.Int16:
			intString = strconv.FormatInt(int64(obj.(int16)), 10)
		case reflect.Int8:
			intString = strconv.FormatInt(int64(obj.(int8)), 10)
		case reflect.Int:
			intString = strconv.FormatInt(int64(obj.(int)), 10)
		}
		schema["type"] = "string"
		if obj != nil {
			schema["default"] = intString
		}
	case reflect.Slice:
		if byteSlice, ok := obj.([]byte); ok {
			if len(byteSlice) > 0 {
				var dat map[string]interface{}
				err := json.Unmarshal(obj.([]byte), &dat)
				if err != nil {
					panic(err)
				}
				RecurseSchema(schema, dat)
			} else {
				RecurseSchema(schema, "")
			}
		} else {
			RecurseSchema(schema, "")
		}
	case reflect.Ptr:
		var ptrDeref interface{}
		if !reflect.ValueOf(obj).IsNil() {
			ptrType := reflect.TypeOf(obj).Elem().Kind()
			switch ptrType {
			case reflect.Int64:
				intDeref := reflect.ValueOf(obj).Elem().Interface().(int64)
				ptrDeref = intDeref
			case reflect.String:
				strDeref := reflect.ValueOf(obj).Elem().Interface().(string)
				ptrDeref = strDeref
			default:
				ptrDeref = reflect.ValueOf(obj).Elem().Interface()
				//ptrDeref = "str"
			}
		} else {
			ptrDeref = ""
		}
		RecurseSchema(schema, ptrDeref)

	default:
		fmt.Printf("%s", obj)
		fmt.Printf("%s", reflect.TypeOf(obj).Elem().Kind())
		fmt.Println("Unknown type")
	}
}
