package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type DateTime struct {
	SchemaField
}

func newDateTime(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id   starlark.String
		name starlark.String
		desc starlark.String
		icon starlark.String
	)

	if err := starlark.UnpackArgs(
		"DateTime",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for DateTime: %s", err)
	}

	s := &DateTime{}
	s.SchemaField.Type = "datetime"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()

	return s, nil
}

func (s *DateTime) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *DateTime) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon",
	}
}

func (s *DateTime) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	default:
		return nil, nil
	}
}

func (s *DateTime) String() string       { return "DateTime(...)" }
func (s *DateTime) Type() string         { return "DateTime" }
func (s *DateTime) Freeze()              {}
func (s *DateTime) Truth() starlark.Bool { return true }

func (s *DateTime) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
