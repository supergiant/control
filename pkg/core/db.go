package core

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-validator/validator"
	"github.com/jinzhu/gorm"
	"github.com/supergiant/supergiant/pkg/model"
)

type DB struct {
	core *Core
	*gorm.DB
}

func (db *DB) Create(m model.Model) error {
	m.SetUUID()
	setDefaultFields(m)
	marshalSerializedFields(m)
	if err := db.validateBelongsTos(m); err != nil {
		return err
	}
	if err := validateFields(m); err != nil {
		return err
	}
	return db.Set("gorm:save_associations", true).Create(m).Error
}

func (db *DB) Save(m model.Model) error {
	marshalSerializedFields(m)
	if err := validateFields(m); err != nil {
		return err
	}
	return db.Set("gorm:save_associations", false).Save(m).Error
}

func (db *DB) Find(out interface{}, where ...interface{}) error {
	if err := db.DB.Find(out, where...).Error; err != nil {
		return err
	}
	items := reflect.ValueOf(out).Elem()
	for i := 0; i < items.Len(); i++ {
		m := items.Index(i).Interface().(model.Model)
		unmarshalSerializedFields(m)
	}
	return nil
}

func (db *DB) First(out interface{}, where ...interface{}) error {
	if err := db.DB.First(out, where...).Error; err != nil {
		return err
	}
	m := out.(model.Model)
	unmarshalSerializedFields(m)
	return nil
}

func (db *DB) Delete(m model.Model) error {
	if m.GetID() == nil {
		return errors.New("ID required for Delete")
	}
	return db.DB.Delete(m).Error
}

// The following are just for the purpose of chaining and preserving our overwritten methods

func (db *DB) Preload(column string, conditions ...interface{}) *DB {
	return &DB{
		db.core,
		db.DB.Preload(column, conditions...),
	}
}

func (db *DB) Where(query interface{}, args ...interface{}) *DB {
	return &DB{
		db.core,
		db.DB.Where(query, args...),
	}
}

func (db *DB) Limit(limit interface{}) *DB {
	return &DB{
		db.core,
		db.DB.Limit(limit),
	}
}

func (db *DB) Offset(offset interface{}) *DB {
	return &DB{
		db.core,
		db.DB.Offset(offset),
	}
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func (db *DB) validateBelongsTos(m model.Model) error {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if belongsTo := tf.ForeignKeyOf; belongsTo != nil && !tf.Field.IsNil() {

			newParent := reflect.New(belongsTo.Field.Type.Elem())

			if err := db.First(newParent.Interface(), tf.Field.Elem().Int()); err != nil {
				return err
			}
			belongsTo.Value.Set(newParent)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Helpers                                                                    //
////////////////////////////////////////////////////////////////////////////////

// setDefaultFields takes a Model with a pointer and sets the default value
// on all fields with the tag sg:"default=something".
func setDefaultFields(m model.Model) {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if tf.Default == nil {
			continue
		}
		existingValueIsZero := tf.Field.Interface() == reflect.Zero(tf.Field.Type()).Interface()
		if existingValueIsZero {
			tf.Field.Set(reflect.ValueOf(tf.Default))
		}
	}
}

// validateFields takes a Model with a pointer and runs a validation on every
// field with the validate:"..." tag.
func validateFields(m model.Model) error {
	return validator.Validate(m)
}

func marshalSerializedFields(m model.Model) {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if jsonField := tf.StoreAsJsonIn; jsonField != nil {
			objField := tf.Field

			if objField.IsNil() {
				continue
			}

			out, err := json.Marshal(objField.Interface())
			if err != nil {
				panic(err)
			}

			jsonField.SetBytes(out)
		}
	}
}

func unmarshalSerializedFields(m model.Model) {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if jsonField := tf.StoreAsJsonIn; jsonField != nil {
			objField := tf.Field

			if jsonField.Len() == 0 {
				continue
			}

			var unmarshalTo reflect.Value

			if objField.Kind() == reflect.Map {
				objField.Set(reflect.MakeMap(objField.Type()))
				unmarshalTo = objField.Addr()

			} else if objField.Kind() == reflect.Slice {
				objField.Set(reflect.MakeSlice(objField.Type(), 0, 0))
				unmarshalTo = objField.Addr()

			} else { // *struct
				objField.Set(reflect.New(objField.Type().Elem()))
				unmarshalTo = objField
			}

			err := json.Unmarshal(jsonField.Bytes(), unmarshalTo.Interface())
			if err != nil {
				panic(err)
			}
		}
	}
}
