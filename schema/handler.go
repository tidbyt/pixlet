package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Handler struct {
	SchemaHandler
}

func newHandler(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		handler *starlark.Function
		type_   starlark.Int
	)

	if err := starlark.UnpackArgs(
		"Handler",
		args, kwargs,
		"handler", &handler,
		"type", &type_,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Handler: %s", err)
	}

	handlerType := HandlerReturnType(type_.BigInt().Int64())
	if handlerType != ReturnSchema &&
		handlerType != ReturnOptions &&
		handlerType != ReturnString &&
		handlerType != ReturnField {
		return nil, fmt.Errorf("invalid handler type %d", int(handlerType))
	}

	s := &Handler{}
	s.Function = handler
	s.ReturnType = handlerType

	return s, nil
}

func (s *Handler) AsSchemaHandler() SchemaHandler {
	return s.SchemaHandler
}

func (s *Handler) AttrNames() []string {
	return []string{
		"handler", "type",
	}
}

func (s *Handler) Attr(name string) (starlark.Value, error) {
	switch name {

	case "handler":
		return s.Function, nil

	case "type":
		return starlark.MakeInt(int(s.ReturnType)), nil

	default:
		return nil, nil
	}
}

func (s *Handler) String() string       { return "Handler(...)" }
func (s *Handler) Type() string         { return "Handler" }
func (s *Handler) Freeze()              {}
func (s *Handler) Truth() starlark.Bool { return true }

func (s *Handler) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
