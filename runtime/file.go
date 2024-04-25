package runtime

import (
	"fmt"
	"io"
	"io/fs"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type File struct {
	fsys fs.FS
	path string
}

func (f File) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlark.String("File"), starlark.StringDict{
		"path":    starlark.String(f.path),
		"readall": starlark.NewBuiltin("readall", f.readall),
	})
}

func (f File) readall(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mode starlark.String
	if err := starlark.UnpackArgs("readall", args, kwargs, "mode?", &mode); err != nil {
		return nil, err
	}

	r, err := f.reader(string(mode))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return r.read(thread, nil, nil, nil)
}

func (f File) reader(mode string) (*Reader, error) {
	var binaryMode bool
	switch mode {
	case "", "r", "rt":
		binaryMode = false

	case "rb":
		binaryMode = true

	default:
		return nil, fmt.Errorf("unsupported mode: %s", mode)
	}

	fl, err := f.fsys.Open(f.path)
	if err != nil {
		return nil, err
	}
	return &Reader{fl, binaryMode}, nil
}

type Reader struct {
	io.ReadCloser
	binaryMode bool
}

func (r Reader) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlark.String("Reader"), starlark.StringDict{
		"read": starlark.NewBuiltin("read", r.read),
		"close": starlark.NewBuiltin("close", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return nil, r.Close()
		}),
	})
}

// read reads the contents of the file. The Starlark signature is:
//
//	read(size=-1) -> bytes
func (r Reader) read(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	starlarkSize := starlark.MakeInt(-1)
	if err := starlark.UnpackArgs("read", args, kwargs, "size?", &starlarkSize); err != nil {
		return nil, err
	}

	var size int
	if err := starlark.AsInt(starlarkSize, &size); err != nil {
		return nil, fmt.Errorf("size is not an int")
	}

	returnType := func(buf []byte) starlark.Value {
		if r.binaryMode {
			return starlark.Bytes(buf)
		} else {
			return starlark.String(buf)
		}
	}

	if size < 0 {
		// read and return all bytes
		buf, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return returnType(buf), nil
	} else {
		// read and return size bytes
		buf := make([]byte, size)
		_, err := r.Read(buf)
		if err != nil {
			return nil, err
		}
		return returnType(buf), nil
	}
}
