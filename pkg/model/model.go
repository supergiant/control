package model

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// Model is an interface that defines the required behaviors of all stored
// (whether in memory or on disk) Supergiant data structure objects.
type Model interface {
	GetID() interface{}
	GetUUID() string
	SetUUID()
	SetActionStatus(*ActionStatus)
	SetPassiveStatus()
}

// BaseModel implements the Model interface, and is composed into all persisted
// Supergiant resources.
type BaseModel struct {
	ID        *int64        `gorm:"primary_key" json:"id,omitempty" sg:"readonly"`
	UUID      string        `json:"uuid" gorm:"not null" sg:"readonly"`
	CreatedAt time.Time     `json:"created_at" sg:"readonly"`
	UpdatedAt time.Time     `json:"updated_at" sg:"readonly"`
	Status    *ActionStatus `gorm:"-" json:"status,omitempty"`

	PassiveStatus     string `gorm:"-" json:"passive_status,omitempty"`
	PassiveStatusOkay bool   `gorm:"-" json:"passive_status_okay,omitempty"`
}

// ActionStatus holds all the information pertaining to any running or failed
// Async Actions, and is rendered on the model on display (not persisted).
type ActionStatus struct {
	Description    string `json:"description"`
	MaxRetries     int    `json:"max_retries"`
	Retries        int    `json:"retries"`
	Error          string `json:"error,omitempty"`
	Cancelled      bool   `json:"cancelled,omitempty"`
	TotalSteps     int    `json:"total_steps,omitempty"`
	StepsCompleted int    `json:"steps_completed,omitempty"`
}

// GetID returns the model ID.
func (m *BaseModel) GetID() interface{} {
	return m.ID
}

// GetUUID returns the model UUID.
func (m *BaseModel) GetUUID() string {
	return m.UUID
}

// SetUUID sets the model UUID.
func (m *BaseModel) SetUUID() {
	if m.UUID == "" {
		u, _ := uuid.NewV4()
		m.UUID = u.String()
	}
}

// SetActionStatus takes an *ActionStatus and sets it on the model.
func (m *BaseModel) SetActionStatus(status *ActionStatus) {
	m.Status = status
}

// SetPassiveStatus implements the Model interface, but does nothing. It can be
// used by models to set render the PassiveStatus and PassiveStatusOkay fields.
func (m *BaseModel) SetPassiveStatus() {
}

// Helpers

// RootFieldJSONNames takes a model and returns the json name tag for every
// top-level field.
func RootFieldJSONNames(m Model) (fields []string) {
	mt := reflect.TypeOf(m).Elem()
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, strings.Split(mt.Field(i).Tag.Get("json"), ",")[0])
	}
	return
}

//------------------------------------------------------------------------------

// BelongsToField holds the reflect StructField (type) and Value of a model
// field representing the parent object in a belongs_to relation.
type BelongsToField struct {
	Field reflect.StructField
	Value reflect.Value
}

// TaggedModelField holds all the information extracted from "sg" tags defined
// on a model's field.
type TaggedModelField struct {
	FieldName     string
	Field         reflect.Value
	Readonly      bool
	Private       bool
	Immutable     bool
	Default       interface{}
	StoreAsJSONIn *reflect.Value
	ForeignKeyOf  *BelongsToField
}

func taggedModelFieldOf(obj reflect.Value, field reflect.StructField, fieldValue reflect.Value) *TaggedModelField {
	tag := field.Tag.Get("sg")
	parts := strings.Split(tag, ",")

	out := new(TaggedModelField)
	out.FieldName = field.Name
	out.Field = fieldValue

	for _, part := range parts {
		subparts := strings.Split(part, "=")
		switch len(subparts) {
		case 1:

			switch subparts[0] {

			case "readonly":
				out.Readonly = true

			case "private":
				out.Private = true

			case "immutable":
				out.Immutable = true

			default:
				panic("Could not parse Model tag " + tag)
			}

		case 2: // e.g. default=10

			switch subparts[0] {

			case "default":
				switch kind := fieldValue.Kind(); kind {
				case reflect.String:
					out.Default = subparts[1] // already a string
				case reflect.Int:
					integer, err := strconv.Atoi(subparts[1])
					if err != nil {
						panic(err)
					}
					out.Default = integer
				default:
					panic("Cannot parse tag default with value " + subparts[1])
				}

			case "store_as_json_in":
				jsonField := obj.FieldByName(subparts[1])
				out.StoreAsJSONIn = &jsonField

			default:
				panic("Could not parse Model tag " + tag)
			}

		default:
			panic("Could not parse Model tag " + tag)
		}
	}

	return out
}

func gatherTaggedModelFieldsInto(obj reflect.Value, taggedFields *[]*TaggedModelField) {
	objType := obj.Type()

	for i := 0; i < obj.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := obj.Field(i)

		// 1. if we see an SG tag, pass it to the tag parsing func, and continue
		// 2. if no SG tag, AND it's a struct (or ptr to), then we have to call recursively
		// 3. if no SG tag, and it's NOT a struct, we don't care

		// Foreign key
		isID := field.Type.Kind() == reflect.String && strings.Contains(field.Tag.Get("gorm"), "index")
		rxp := regexp.MustCompile(`^(\w+)Name$`)
		if isID && rxp.MatchString(field.Name) {
			outFieldName := rxp.FindStringSubmatch(field.Name)[1]
			belongsToField, found := objType.FieldByName(outFieldName)
			if found {
				taggedField := &TaggedModelField{
					Field: fieldValue,
					ForeignKeyOf: &BelongsToField{
						Field: belongsToField,
						Value: obj.FieldByName(outFieldName),
					},
				}
				*taggedFields = append(*taggedFields, taggedField)
			}
		}

		// SG tag
		if tag := field.Tag.Get("sg"); tag != "" {
			taggedField := taggedModelFieldOf(obj, field, fieldValue)
			*taggedFields = append(*taggedFields, taggedField)
			// continue
		}

		indirectFieldValue := reflect.Indirect(fieldValue)

		switch {
		case fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.Ptr && fieldValue.Type().Elem().Elem().Kind() == reflect.Struct:
			for j := 0; j < fieldValue.Len(); j++ {
				gatherTaggedModelFieldsInto(fieldValue.Index(j).Elem(), taggedFields)
			}
		case indirectFieldValue.Kind() == reflect.Struct:
			gatherTaggedModelFieldsInto(indirectFieldValue, taggedFields)
		}
	}
}

// TaggedModelFieldsOf takes a Model and returns every *TaggedModelField defined
// (fields are gathered recursively).
func TaggedModelFieldsOf(r Model) (taggedFields []*TaggedModelField) {
	resourceValue := reflect.ValueOf(r).Elem()
	gatherTaggedModelFieldsInto(resourceValue, &taggedFields)
	return
}

// ZeroReadonlyFields takes a Model with pointer, and zeroes any fields with
// the tag sg:"readonly".
func ZeroReadonlyFields(r Model) {
	for _, tf := range TaggedModelFieldsOf(r) {
		if tf.Readonly {
			tf.Field.Set(reflect.Zero(tf.Field.Type()))
		}
	}
}

// ZeroPrivateFields takes a Model with pointer, and zeroes any fields with
// the tag sg:"private".
func ZeroPrivateFields(r Model) {
	for _, tf := range TaggedModelFieldsOf(r) {
		if tf.Private {
			tf.Field.Set(reflect.Zero(tf.Field.Type()))
		}
	}
}

// ErrorChangedImmutableField is an error for when a model has a value defined
// in a field with the sg tag "immutable" (used for updates).
type ErrorChangedImmutableField struct {
	fieldName string
}

func (err *ErrorChangedImmutableField) Error() string {
	return err.fieldName + " cannot be changed"
}

// CheckImmutableFields returns an error if any fields with the sg tag
// "immutable" have values (are not zero).
func CheckImmutableFields(m Model) error {
	for _, tf := range TaggedModelFieldsOf(m) {
		if tf.Immutable {
			isZero := reflect.DeepEqual(tf.Field.Interface(), reflect.Zero(tf.Field.Type()).Interface())
			if !isZero {
				return &ErrorChangedImmutableField{tf.FieldName}
			}
		}
	}
	return nil
}
