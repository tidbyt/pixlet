package schema

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Toggle struct {
	SchemaField
}

func newToggle(
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
		def  starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Toggle",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default?", &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Toggle: %s", err)
	}

	s := &Toggle{}
	s.SchemaField.Type = "onoff"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Default = strconv.FormatBool(bool(def))

	return s, nil
}

func (s *Toggle) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Toggle) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default",
	}
}

func (s *Toggle) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	case "default":
		b, _ := strconv.ParseBool(s.Default)
		return starlark.Bool(b), nil

	default:
		return nil, nil
	}
}

func (s *Toggle) String() string       { return "Toggle(...)" }
func (s *Toggle) Type() string         { return "Toggle" }
func (s *Toggle) Freeze()              {}
func (s *Toggle) Truth() starlark.Bool { return true }

func (s *Toggle) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
