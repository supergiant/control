package model

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type Model interface {
	GetID() interface{}
	GetUUID() string
	SetUUID()
	SetActionStatus(*ActionStatus)
	SetPassiveStatus()
}

type BaseModel struct {
	ID        *int64        `gorm:"primary_key" json:"id,omitempty" sg:"readonly"`
	UUID      string        `json:"uuid,omitempty" sg:"readonly"`
	CreatedAt time.Time     `json:"created_at,omitempty" sg:"readonly"` // TODO won't be omitted cuz not *time.Time
	UpdatedAt time.Time     `json:"updated_at,omitempty" sg:"readonly"`
	Status    *ActionStatus `gorm:"-" json:"status,omitempty"`

	PassiveStatus     string `gorm:"-" json:"passive_status,omitempty"`
	PassiveStatusOkay bool   `gorm:"-" json:"passive_status_okay,omitempty"`
}

type ActionStatus struct {
	Description string `json:"description"`
	MaxRetries  int    `json:"max_retries"`
	Retries     int    `json:"retries"`
	Error       string `json:"error,omitempty"`
	Cancelled   bool   `json:"cancelled,omitempty"`
}

func (m *BaseModel) GetID() interface{} {
	return m.ID
}

func (m *BaseModel) GetUUID() string {
	return m.UUID
}

func (m *BaseModel) SetUUID() {
	if m.UUID == "" {
		m.UUID = uuid.NewV4().String()
	}
}

func (m *BaseModel) SetActionStatus(status *ActionStatus) {
	m.Status = status
}

func (m *BaseModel) SetPassiveStatus() {
}

//------------------------------------------------------------------------------------- helpers below

func RootFieldJSONNames(m Model) (fields []string) {
	mt := reflect.TypeOf(m).Elem()
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, strings.Split(mt.Field(i).Tag.Get("json"), ",")[0])
	}
	return
}

//------------------------------------------------------------------------------

type BelongsToField struct {
	Field reflect.StructField
	Value reflect.Value
}

type TaggedModelField struct {
	Field         reflect.Value
	Readonly      bool
	Private       bool
	Default       interface{}
	StoreAsJsonIn *reflect.Value
	ForeignKeyOf  *BelongsToField
}

func taggedModelFieldOf(obj reflect.Value, field reflect.StructField, fieldValue reflect.Value) *TaggedModelField {
	tag := field.Tag.Get("sg")
	parts := strings.Split(tag, ",")

	out := new(TaggedModelField)
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
				out.StoreAsJsonIn = &jsonField

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
