package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Option struct {
	SchemaOption
}

func newOption(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		display starlark.String
		value   starlark.String
	)

	if err := starlark.UnpackArgs(
		"Option",
		args, kwargs,
		"display", &display,
		"value", &value,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Option: %s", err)
	}

	s := &Option{}
	s.SchemaOption.Text = display.GoString()
	s.SchemaOption.Value = value.GoString()

	return s, nil
}

func (s *Option) AsSchemaOption() SchemaOption {
	return s.SchemaOption
}

func (s *Option) AttrNames() []string {
	return []string{
		"display", "value",
	}
}

func (s *Option) Attr(name string) (starlark.Value, error) {
	switch name {

	case "display":
		return starlark.String(s.Text), nil

	case "value":
		return starlark.String(s.Value), nil

	default:
		return nil, nil
	}
}

func (s *Option) String() string       { return "Option(...)" }
func (s *Option) Type() string         { return "Option" }
func (s *Option) Freeze()              {}
func (s *Option) Truth() starlark.Bool { return true }

func (s *Option) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
