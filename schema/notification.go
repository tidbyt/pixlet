package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Notification struct {
	SchemaField
	Builder        *starlark.Function `json:"-"`
	starlarkSounds *starlark.List
}

func newNotification(
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
		sounds  *starlark.List
		builder *starlark.Function
	)

	if err := starlark.UnpackArgs(
		"Notification",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"sounds", &sounds,
		"builder", &builder,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Notification: %s", err)
	}

	s := &Notification{}
	s.SchemaField.Type = "notification"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Builder = builder

	var soundVal starlark.Value
	soundIter := sounds.Iterate()
	defer soundIter.Done()
	for i := 0; soundIter.Next(&soundVal); {
		if _, isNone := soundVal.(starlark.NoneType); isNone {
			continue
		}

		o, ok := soundVal.(*Sound)
		if !ok {
			return nil, fmt.Errorf(
				"expected options to be a list of Sound but found: %s (at index %d)",
				soundVal.Type(),
				i,
			)
		}

		s.Sounds = append(s.Sounds, o.AsSchemaSound())
	}
	s.starlarkSounds = sounds

	return s, nil
}

func (s *Notification) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Notification) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "sounds", "builder",
	}
}

func (s *Notification) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	case "sounds":
		return s.starlarkSounds, nil

	case "builder":
		return s.Builder, nil

	default:
		return nil, nil
	}
}

func (s *Notification) String() string       { return "Notification(...)" }
func (s *Notification) Type() string         { return "Notification" }
func (s *Notification) Freeze()              {}
func (s *Notification) Truth() starlark.Bool { return true }

func (s *Notification) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
