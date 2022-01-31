package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type PhotoSelect struct {
	SchemaField
}

func newPhotoSelect(
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
		"PhotoSelect",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for PhotoSelect: %s", err)
	}

	s := &PhotoSelect{}
	s.SchemaField.Type = "png"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()

	return s, nil
}

func (s *PhotoSelect) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *PhotoSelect) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon",
	}
}

func (s *PhotoSelect) Attr(name string) (starlark.Value, error) {
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

func (s *PhotoSelect) String() string       { return "PhotoSelect(...)" }
func (s *PhotoSelect) Type() string         { return "PhotoSelect" }
func (s *PhotoSelect) Freeze()              {}
func (s *PhotoSelect) Truth() starlark.Bool { return true }

func (s *PhotoSelect) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
