package schema

import (
	"fmt"
	"sync"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "schema"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		handlerType := starlarkstruct.FromStringDict(
			starlark.String("HandlerType"),
			map[string]starlark.Value{
				"Schema":  starlark.MakeInt(int(ReturnSchema)),
				"Options": starlark.MakeInt(int(ReturnOptions)),
				"String":  starlark.MakeInt(int(ReturnString)),
				"Field":   starlark.MakeInt(int(ReturnField)),
			},
		)

		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"Schema":        starlark.NewBuiltin("Schema", newSchema),
					"Toggle":        starlark.NewBuiltin("Toggle", newToggle),
					"Option":        starlark.NewBuiltin("Option", newOption),
					"Dropdown":      starlark.NewBuiltin("Dropdown", newDropdown),
					"Location":      starlark.NewBuiltin("Location", newLocation),
					"Text":          starlark.NewBuiltin("Text", newText),
					"LocationBased": starlark.NewBuiltin("LocationBased", newLocationBased),
					"DateTime":      starlark.NewBuiltin("DateTime", newDateTime),
					"OAuth2":        starlark.NewBuiltin("OAuth2", newOAuth2),
					"PhotoSelect":   starlark.NewBuiltin("PhotoSelect", newPhotoSelect),
					"Typeahead":     starlark.NewBuiltin("Typeahead", newTypeahead),
					"Handler":       starlark.NewBuiltin("Handler", newHandler),
					"HandlerType":   handlerType,
					"Generated":     starlark.NewBuiltin("Generated", newGenerated),
					"Color":         starlark.NewBuiltin("Color", newColor),
				},
			},
		}
	})

	return module, nil
}

type Field interface {
	AsSchemaField() SchemaField
}

type StarlarkSchema struct {
	Schema
	Handlers         map[string]SchemaHandler
	starlarkFields   *starlark.List
	starlarkHandlers *starlark.List
}

func newSchema(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		version  starlark.String
		fields   *starlark.List
		handlers *starlark.List
	)

	if err := starlark.UnpackArgs(
		"Schema",
		args, kwargs,
		"version", &version,
		"fields", &fields,
		"handlers?", &handlers,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Schema: %s", err)
	}

	if version.GoString() != "1" {
		return nil, fmt.Errorf("only schema version 1 is supported, not: %s", version.GoString())
	}

	s := &StarlarkSchema{
		Schema: Schema{
			Version: version.GoString(),
		},
		Handlers:         map[string]SchemaHandler{},
		starlarkFields:   fields,
		starlarkHandlers: handlers,
	}

	if s.starlarkFields != nil {
		fieldIter := s.starlarkFields.Iterate()
		defer fieldIter.Done()

		var fieldVal starlark.Value
		for i := 0; fieldIter.Next(&fieldVal); {
			if _, isNone := fieldVal.(starlark.NoneType); isNone {
				continue
			}

			f, ok := fieldVal.(Field)
			if !ok {
				return nil, fmt.Errorf(
					"expected fields to be a list of Field but found: %s (at index %d)",
					fieldVal.Type(),
					i,
				)
			}

			s.Schema.Fields = append(s.Schema.Fields, f.AsSchemaField())
		}
	}

	if s.starlarkHandlers != nil {
		handlerIter := s.starlarkHandlers.Iterate()
		defer handlerIter.Done()

		var handlerVal starlark.Value
		for i := 0; handlerIter.Next(&handlerVal); {
			handler, ok := handlerVal.(*Handler)
			if !ok {
				return nil, fmt.Errorf(
					"expected handlers to hold Handler but found: %s (at index %d)",
					handlerVal.Type(),
					i,
				)
			}
			s.Handlers[handler.Function.Name()] = handler.SchemaHandler
		}
	}

	return s, nil
}

func (s StarlarkSchema) AttrNames() []string {
	return []string{
		"version",
		"fields",
		"handlers",
	}
}

func (s StarlarkSchema) Attr(name string) (starlark.Value, error) {
	switch name {
	case "version":
		return starlark.String(s.Version), nil

	case "fields":
		return s.starlarkFields, nil

	case "handlers":
		return s.starlarkHandlers, nil

	default:
		return nil, nil
	}
}

func (s StarlarkSchema) String() string       { return "Schema(...)" }
func (s StarlarkSchema) Type() string         { return "StarlarkSchema" }
func (s StarlarkSchema) Freeze()              {}
func (s StarlarkSchema) Truth() starlark.Bool { return true }

func (s StarlarkSchema) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
