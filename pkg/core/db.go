package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-validator/validator"
	"github.com/jinzhu/gorm"

	"github.com/supergiant/supergiant/pkg/model"
)

type DBInterface interface {
	Create(model.Model) error
	Save(model.Model) error
	Find(out interface{}, where ...interface{}) error
	First(out interface{}, where ...interface{}) error
	Delete(m model.Model) error
	Preload(column string, conditions ...interface{}) DBInterface
	Where(query interface{}, args ...interface{}) DBInterface
	Limit(limit interface{}) DBInterface
	Offset(offset interface{}) DBInterface
	Model(value interface{}) DBInterface
	Update(attrs ...interface{}) error
	Count(interface{}) error
}

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
	// if m.GetID() == nil {
	// 	return errors.New("ID required for Delete")
	// }
	return db.DB.Delete(m).Error
}

// The following are just for the purpose of chaining and preserving our overwritten methods

func (db *DB) Preload(column string, conditions ...interface{}) DBInterface {
	return &DB{
		db.core,
		db.DB.Preload(column, conditions...),
	}
}

func (db *DB) Where(query interface{}, args ...interface{}) DBInterface {
	return &DB{
		db.core,
		db.DB.Where(query, args...),
	}
}

func (db *DB) Limit(limit interface{}) DBInterface {
	return &DB{
		db.core,
		db.DB.Limit(limit),
	}
}

func (db *DB) Offset(offset interface{}) DBInterface {
	return &DB{
		db.core,
		db.DB.Offset(offset),
	}
}

func (db *DB) Model(value interface{}) DBInterface {
	return &DB{
		db.core,
		db.DB.Model(value),
	}
}

func (db *DB) Update(attrs ...interface{}) error {
	return db.DB.Update(attrs...).Error
}

func (db *DB) Count(value interface{}) error {
	return db.DB.Count(value).Error
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

type ErrorMissingRequiredParent struct {
	key   string
	model string
}

func (err *ErrorMissingRequiredParent) Error() string {
	return fmt.Sprintf("Parent does not exist, foreign key '%s' on %s", err.key, err.model)
}

func (db *DB) validateBelongsTos(m model.Model) error {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if belongsTo := tf.ForeignKeyOf; belongsTo != nil && tf.Field.String() != "" {
			newParent := reflect.New(belongsTo.Field.Type.Elem())
			if err := db.Where("name = ?", tf.Field.String()).First(newParent.Interface()); err != nil {
				keyName := strings.Split(newParent.Elem().Type().String(), ".")[1] + "Name"
				modelName := strings.Split(reflect.ValueOf(m).Elem().Type().String(), ".")[1]
				return &ErrorMissingRequiredParent{keyName, modelName}
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

type ErrorValidationFailed struct {
	error
}

func (err *ErrorValidationFailed) Error() string {
	return "Validation failed: " + err.error.Error()
}

// validateFields takes a Model with a pointer and runs a validation on every
// field with the validate:"..." tag.
func validateFields(m model.Model) error {
	if err := validator.Validate(m); err != nil {
		return &ErrorValidationFailed{err}
	}
	return nil
}

func marshalSerializedFields(m model.Model) {
	for _, tf := range model.TaggedModelFieldsOf(m) {
		if jsonField := tf.StoreAsJSONIn; jsonField != nil {
			objField := tf.Field

			if objField.IsNil() || (reflect.Indirect(objField).Kind() == reflect.Slice && reflect.Indirect(objField).Len() == 0) {
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
		if jsonField := tf.StoreAsJSONIn; jsonField != nil {
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
