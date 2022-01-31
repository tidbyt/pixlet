package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Typeahead struct {
	SchemaField
	starlarkHandler *starlark.Function
}

func newTypeahead(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id      starlark.String
		name    starlark.String
		desc    starlark.String
		icon    starlark.String
		handler *starlark.Function
	)

	if err := starlark.UnpackArgs(
		"Typeahead",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"handler", &handler,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Typeahead: %s", err)
	}

	s := &Typeahead{}
	s.SchemaField.Type = "typeahead"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Handler = handler.Name()
	s.starlarkHandler = handler

	return s, nil
}

func (s *Typeahead) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Typeahead) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "handler",
	}
}

func (s *Typeahead) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	case "handler":
		return s.starlarkHandler, nil

	default:
		return nil, nil
	}
}

func (s *Typeahead) String() string       { return "Typeahead(...)" }
func (s *Typeahead) Type() string         { return "Typeahead" }
func (s *Typeahead) Freeze()              {}
func (s *Typeahead) Truth() starlark.Bool { return true }

func (s *Typeahead) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
