package tools

import (
	"io/fs"
	"os"
	"path/filepath"
)

type SingleFileFS struct {
	Path   string
	baseFS fs.FS
}

func NewSingleFileFS(filePath string) *SingleFileFS {
	return &SingleFileFS{
		Path:   filePath,
		baseFS: os.DirFS(filepath.Dir(filePath)),
	}
}

func (sfs *SingleFileFS) Open(name string) (fs.File, error) {
	if name != "." && name != filepath.Base(sfs.Path) {
		return nil, fs.ErrNotExist
	}

	return sfs.baseFS.Open(name)
}
