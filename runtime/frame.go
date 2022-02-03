package runtime

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render"
)

type Root struct {
	render.Root
	starlarkChild  starlark.Value
	starlarkDelay  starlark.Int
	starlarkMaxAge starlark.Int
}

func newRoot(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var child starlark.Value
	var delay starlark.Int
	var maxAge starlark.Int

	if err := starlark.UnpackArgs(
		"Root",
		args, kwargs,
		"child", &child,
		"delay?", &delay,
		"max_age?", &maxAge,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Root: %s", err)
	}

	r := Root{
		starlarkChild:  child,
		starlarkDelay:  delay,
		starlarkMaxAge: maxAge,
	}
	r.Delay = int32(delay.BigInt().Int64())
	r.MaxAge = int32(maxAge.BigInt().Int64())

	w, ok := child.(Widget)
	if !ok {
		return nil, fmt.Errorf("invalid type for child: %s (expected a Widget)", child.Type())
	}
	r.Child = w.AsRenderWidget()

	return r, nil
}

func (r Root) AsRenderRoot() render.Root {
	return r.Root
}

func (r Root) AttrNames() []string {
	return []string{
		"child",
		"delay",
		"max_age",
	}
}

func (r Root) Attr(name string) (starlark.Value, error) {
	switch name {
	case "child":
		return r.starlarkChild, nil

	case "delay":
		return r.starlarkDelay, nil

	case "max_age":
		return r.starlarkMaxAge, nil

	default:
		return nil, nil
	}
}

func (r Root) String() string       { return "Root(...)" }
func (r Root) Type() string         { return "Root" }
func (r Root) Freeze()              {}
func (r Root) Truth() starlark.Bool { return true }

func (r Root) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(r, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
