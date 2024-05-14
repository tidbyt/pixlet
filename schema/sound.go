package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/runtime/modules/file"
)

type Sound struct {
	SchemaSound
	file *file.File
}

func newSound(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id    starlark.String
		title starlark.String
		file  *file.File
	)

	if err := starlark.UnpackArgs(
		"Sound",
		args, kwargs,
		"id", &id,
		"title", &title,
		"file", &file,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Sound: %s", err)
	}

	s := &Sound{file: file}
	s.ID = id.GoString()
	s.Title = title.GoString()
	s.Path = file.Path

	return s, nil
}

func (s *Sound) AsSchemaSound() SchemaSound {
	return s.SchemaSound
}

func (s *Sound) AttrNames() []string {
	return []string{"id", "title", "file"}
}

func (s *Sound) Attr(name string) (starlark.Value, error) {
	switch name {
	case "id":
		return starlark.String(s.ID), nil

	case "title":
		return starlark.String(s.Title), nil

	case "file":
		return s.file, nil

	default:
		return nil, nil
	}
}

func (s *Sound) String() string       { return "Sound(...)" }
func (s *Sound) Type() string         { return "Sound" }
func (s *Sound) Freeze()              {}
func (s *Sound) Truth() starlark.Bool { return true }

func (s *Sound) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
