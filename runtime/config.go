package runtime

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type AppletConfig map[string]string

func (a AppletConfig) AttrNames() []string {
	return []string{
		"get",
		"str",
		"bool",
	}
}

func (a AppletConfig) Attr(name string) (starlark.Value, error) {
	switch name {

	case "get", "str":
		return starlark.NewBuiltin("str", a.getString), nil

	case "bool":
		return starlark.NewBuiltin("bool", a.getBoolean), nil

	default:
		return nil, nil
	}
}

func (a AppletConfig) Get(key starlark.Value) (starlark.Value, bool, error) {
	switch v := key.(type) {
	case starlark.String:
		val, found := a[v.GoString()]
		return starlark.String(val), found, nil
	default:
		return nil, false, nil
	}
}

func (a AppletConfig) String() string       { return "AppletConfig(...)" }
func (a AppletConfig) Type() string         { return "AppletConfig" }
func (a AppletConfig) Freeze()              {}
func (a AppletConfig) Truth() starlark.Bool { return true }

func (a AppletConfig) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(a, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func (a AppletConfig) getString(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	var def starlark.Value
	def = starlark.None

	if err := starlark.UnpackPositionalArgs(
		"str", args, kwargs, 1,
		&key, &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for config.str: %v", err)
	}

	val, ok := a[key.GoString()]
	if !ok {
		return def, nil
	} else {
		return starlark.String(val), nil
	}
}

func (a AppletConfig) getBoolean(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	var def starlark.Value
	def = starlark.None

	if err := starlark.UnpackPositionalArgs(
		"bool", args, kwargs, 1,
		&key, &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for config.bool: %v", err)
	}

	val, ok := a[key.GoString()]
	if !ok {
		return def, nil
	} else {
		b, _ := strconv.ParseBool(val)
		return starlark.Bool(b), nil
	}
}
