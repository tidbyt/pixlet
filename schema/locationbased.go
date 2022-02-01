package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type LocationBased struct {
	SchemaField
	starlarkHandler *starlark.Function
}

func newLocationBased(
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
		"LocationBased",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"handler", &handler,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for LocationBased: %s", err)
	}

	s := &LocationBased{}
	s.SchemaField.Type = "locationbased"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Handler = handler.Name()
	s.starlarkHandler = handler

	return s, nil
}

func (s *LocationBased) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *LocationBased) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "handler",
	}
}

func (s *LocationBased) Attr(name string) (starlark.Value, error) {
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

func (s *LocationBased) String() string       { return "LocationBased(...)" }
func (s *LocationBased) Type() string         { return "LocationBased" }
func (s *LocationBased) Freeze()              {}
func (s *LocationBased) Truth() starlark.Bool { return true }

func (s *LocationBased) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
