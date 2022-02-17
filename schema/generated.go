package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Generated struct {
	SchemaField
	starlarkHandler *starlark.Function
}

func newGenerated(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id      starlark.String
		source  starlark.String
		handler *starlark.Function
	)

	if err := starlark.UnpackArgs(
		"Generated",
		args, kwargs,
		"source", &source,
		"handler", &handler,
		"id", &id,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Generated: %s", err)
	}

	s := &Generated{}
	s.starlarkHandler = handler
	s.Source = source.GoString()
	s.Handler = handler.Name()
	s.ID = id.GoString()
	s.SchemaField.Type = "generated"

	return s, nil
}

func (s *Generated) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Generated) AttrNames() []string {
	return []string{
		"source", "handler", "id",
	}
}

func (s *Generated) Attr(name string) (starlark.Value, error) {
	switch name {

	case "source":
		return starlark.String(s.Source), nil

	case "handler":
		return s.starlarkHandler, nil
	case "id":
		return starlark.String(s.ID), nil
	default:
		return nil, nil
	}
}

func (s *Generated) String() string       { return "Generated(...)" }
func (s *Generated) Type() string         { return "Generated" }
func (s *Generated) Freeze()              {}
func (s *Generated) Truth() starlark.Bool { return true }

func (s *Generated) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
