package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.starlark.net/starlark"
)

const (
	ReturnSchema HandlerReturnType = iota
	ReturnOptions
	ReturnString
	ReturnField
)

const (
	// SchemaFunctionName is the name of the function in Starlark that we expect
	// to be able to call to get the schema for an applet.
	SchemaFunctionName = "get_schema"
)

// Schema holds a configuration object for an applet. It holds a list of fields
// that are exported from an applet.
type Schema struct {
	Version       string        `json:"version" validate:"required"`
	Fields        []SchemaField `json:"schema" validate:"dive"`
	Notifications []SchemaField `json:"notifications,omitempty" validate:"dive"`

	Handlers map[string]SchemaHandler `json:"-"`
}

// SchemaField represents an item in the config used to confgure an applet.
type SchemaField struct {
	Type        string            `json:"type" validate:"required,oneof=color datetime dropdown generated location locationbased onoff radio text typeahead oauth2 oauth1 png notification"`
	ID          string            `json:"id" validate:"required"`
	Name        string            `json:"name,omitempty" validate:"required_for=datetime dropdown location locationbased onoff radio text typeahead png"`
	Description string            `json:"description,omitempty"`
	Icon        string            `json:"icon,omitempty" validate:"forbidden_for=generated"`
	Visibility  *SchemaVisibility `json:"visibility,omitempty" validate:"omitempty"`

	Default string         `json:"default,omitempty" validate:"required_for=dropdown onoff radio"`
	Options []SchemaOption `json:"options,omitempty" validate:"required_for=dropdown radio,dive"`
	Palette []string       `json:"palette,omitempty"`
	Sounds  []SchemaSound  `json:"sounds,omitempty" validate:"required_for=notification,dive"`

	Source  string `json:"source,omitempty" validate:"required_for=generated"`
	Handler string `json:"handler,omitempty" validate:"required_for=generated locationbased typeahead oauth2"`

	ClientID              string   `json:"client_id,omitempty" validate:"required_for=oauth2"`
	AuthorizationEndpoint string   `json:"authorization_endpoint,omitempty" validate:"required_for=oauth2"`
	Scopes                []string `json:"scopes,omitempty" validate:"required_for=oauth2"`
}

// SchemaOption represents an option in a field. For example, an item in a drop
// down menu.
type SchemaOption struct {
	Display string `json:"display"`
	Text    string `json:"text" validate:"required"` // The same as display, for legacy reasons.
	Value   string `json:"value" validate:"required"`
}

// SchemaSound represents a sound that can be played by the applet.
type SchemaSound struct {
	ID    string `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Path  string `json:"path" validate:"required"`
}

// SchemaVisibility enables conditional fields inside of the mobile app. For
// example, if a field should be invisible until a login is provided.
type SchemaVisibility struct {
	Type      string `json:"type" validate:"required,oneof=invisible disabled"`
	Condition string `json:"condition" validate:"required,oneof=equal not_equal"`
	Variable  string `json:"variable" validate:"required"`
	Value     string `json:"value"`
}

// HandlerReturnType defines an enum for the type of information we expect to
// get back from the schema function.
type HandlerReturnType int8

// SchemaHandler defines a function and and return type for getting the schema
// for an applet. This can both be the predefined schema function we expect all
// applets to have for config, but can also be used as a callback for
type SchemaHandler struct {
	Function   *starlark.Function
	ReturnType HandlerReturnType
}

func (s Schema) MarshalJSON() ([]byte, error) {
	type OriginalSchema Schema

	a := struct {
		OriginalSchema
	}{
		OriginalSchema: (OriginalSchema)(s),
	}

	// ensure that fields is serialized as "[]" and not "null",
	// even if there are no fields. otherwise the Tidbyt mobile app breaks
	if a.Fields == nil {
		a.Fields = make([]SchemaField, 0)
	}
	if a.Notifications == nil {
		a.Notifications = make([]SchemaField, 0)
	}

	js, err := json.Marshal(a)

	return js, err
}

// FromStarlark creates a new Schema from a Starlark schema object.
func FromStarlark(
	val starlark.Value,
	globals starlark.StringDict,
) (*Schema, error) {
	var schema *Schema

	starlarkSchema, ok := val.(*StarlarkSchema)
	if ok {
		schema = &starlarkSchema.Schema
		if schema.Handlers == nil {
			schema.Handlers = make(map[string]SchemaHandler)
			for name, schemaHandler := range starlarkSchema.Handlers {
				schema.Handlers[name] = schemaHandler
			}
		}
	} else {
		schemaTree, err := unmarshalStarlark(val)
		if err != nil {
			return nil, err
		}

		treeJSON, err := json.Marshal(schemaTree)
		if err != nil {
			return nil, err
		}

		schema = &Schema{
			Version:  "1",
			Handlers: make(map[string]SchemaHandler),
		}
		if err := json.Unmarshal(treeJSON, &schema.Fields); err != nil {
			return nil, err
		}
	}

	err := validateSchema(schema)
	if err != nil {
		return nil, err
	}

	for i, schemaField := range schema.Fields {
		if schemaField.Handler != "" {
			handlerValue, found := globals[schemaField.Handler]
			if !found {
				return nil, fmt.Errorf(
					"field %d references non-existent handler \"%s\"",
					i,
					schemaField.Handler)
			}

			handlerFun, ok := handlerValue.(*starlark.Function)
			if !ok {
				return nil, fmt.Errorf(
					"field %d references \"%s\" which is not a function",
					i, schemaField.Handler)
			}

			var handlerType HandlerReturnType
			switch schemaField.Type {
			case "locationbased":
				handlerType = ReturnOptions
			case "generated":
				handlerType = ReturnSchema
			case "typeahead":
				handlerType = ReturnOptions
			case "oauth2":
				handlerType = ReturnString
			case "oauth1":
				handlerType = ReturnString
			default:
				return nil, fmt.Errorf(
					"field %d of type \"%s\" can't have a handler function",
					i, schemaField.Type)
			}

			schema.Handlers[schemaField.Handler] = SchemaHandler{Function: handlerFun, ReturnType: handlerType}
		}
	}

	return schema, nil
}

// Encodes a list of schema options into validated json.
func EncodeOptions(
	starlarkOptions starlark.Value,
) (string, error) {
	optionsTree, err := unmarshalStarlark(starlarkOptions)
	if err != nil {
		return "", err
	}

	options, err := buildOptions(optionsTree)
	if err != nil {
		return "", err
	}

	validate := validator.New()
	for _, o := range options {
		err = validate.Struct(o)
		if err != nil {
			return "", err
		}

		if o.Display == "" {
			o.Display = o.Text
		}
	}

	optionsJson, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	return string(optionsJson), nil
}

// Transforms a starlark value into Go objects. The value must be a
// list, dict or string. Or a tree of these.
func unmarshalStarlark(object starlark.Value) (interface{}, error) {
	switch v := object.(type) {

	case starlark.String:
		return v.GoString(), nil

	case *starlark.List:
		goList := make([]interface{}, 0, v.Len())
		iter := v.Iterate()
		defer iter.Done()

		var listVal starlark.Value
		for i := 0; iter.Next(&listVal); i++ {
			goVal, err := unmarshalStarlark(listVal)
			if err != nil {
				return nil, err
			}
			goList = append(goList, goVal)
		}

		return goList, nil

	case *starlark.Dict:
		goMap := make(map[string]interface{})

		for _, key := range v.Keys() {
			strKey, ok := key.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("dict keys must be string")
			}
			goKey := strKey.GoString()

			value, _, _ := v.Get(key)
			goVal, err := unmarshalStarlark(value)
			if err != nil {
				return nil, err
			}

			goMap[goKey] = goVal
		}
		return goMap, nil
	case *Option:
		return v.AsSchemaOption(), nil
	case Field:
		return v.AsSchemaField(), nil
	}

	return nil, fmt.Errorf("type %s not allowed in schema", object.Type())
}

// Helper. Verifies that object is string or nil and writes it to *p
// if string.
func setString(p *string, object interface{}) string {
	if object == nil {
		return ""
	}
	objectStr, ok := object.(string)
	if !ok {
		return fmt.Sprintf("expected string, found %T", object)
	}
	*p = objectStr
	return ""
}

func buildOptions(options interface{}) ([]SchemaOption, error) {
	optionsList, ok := options.([]interface{})
	if !ok {
		return nil, fmt.Errorf("options must be a list")
	}
	schemaOptions := make([]SchemaOption, 0, len(optionsList))
	for j, o := range optionsList {
		op, ok := o.(SchemaOption)
		if ok {
			schemaOptions = append(schemaOptions, op)
			continue
		}

		oMap, ok := o.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("entry %d is not dict or Option", j)
		}

		sop := SchemaOption{}
		if err := setString(&sop.Text, oMap["text"]); err != "" {
			return nil, fmt.Errorf("option %d has bad text: %s", j, err)
		}
		if err := setString(&sop.Value, oMap["value"]); err != "" {
			return nil, fmt.Errorf("option %d has bad value: %s", j, err)
		}

		schemaOptions = append(schemaOptions, sop)
	}
	return schemaOptions, nil
}

// Validates a Schema object.
func validateSchema(schema *Schema) error {
	// This custom validator function implements
	// "required_for", which makes the tagged field required
	// whenever SchemaField.Type matches one of the parameters.
	requiredFor := func(fl validator.FieldLevel) bool {
		var isSet bool
		switch fl.Field().Kind() {
		case reflect.Map, reflect.Ptr, reflect.Interface:
			isSet = !fl.Field().IsNil()
		case reflect.String, reflect.Int, reflect.Slice:
			isSet = fl.Field().IsValid() && !fl.Field().IsZero()
		default:
			return false
		}

		schemaField := fl.Parent().Interface().(SchemaField)
		for _, fieldType := range strings.Split(fl.Param(), " ") {
			if fieldType == schemaField.Type {
				return isSet
			}
		}
		return true
	}

	// This implements "forbidden_for". Same idea as required for,
	// but the opposite.
	forbiddenFor := func(fl validator.FieldLevel) bool {
		var isSet bool
		switch fl.Field().Kind() {
		case reflect.Map, reflect.Ptr, reflect.Interface:
			isSet = !fl.Field().IsNil()
		case reflect.String, reflect.Int, reflect.Slice:
			isSet = fl.Field().IsValid() && !fl.Field().IsZero()
		default:
			return false
		}

		schemaField := fl.Parent().Interface().(SchemaField)
		for _, fieldType := range strings.Split(fl.Param(), " ") {
			if fieldType == schemaField.Type {
				return !isSet
			}
		}
		return true
	}

	// NOTE: It could be helpful to also provide an "optional_for"
	// function, to make sure we catch superfluous tags.

	validate := validator.New()
	validate.RegisterValidation("required_for", requiredFor)
	validate.RegisterValidation("forbidden_for", forbiddenFor)

	err := validate.Struct(schema)
	if err != nil {
		return err
	}

	return nil
}
