// Package bundle provides primitives for bundling apps for portability.
package bundle

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/nlepage/go-tarfs"

	"tidbyt.dev/pixlet/manifest"
)

const (
	// AppSourceName is the required name to be used for the source file. This
	// is needed so we know what we're searching for in the archive. Note, we
	// rename an app filename from a manifest when we create the bundle to
	// ensure it can be unpacked. We could get around this if we loaded all
	// files in the bundle, though we risk abuse with really large bundles.
	AppSourceName = "app.star"
	// AppBundleName is the standard name for a created bundle.
	AppBundleName = "bundle.tar.gz"
)

// AppBundle represents the unpacked bundle in our system.
type AppBundle struct {
	Manifest *manifest.Manifest
	Source   fs.FS
}

func FromFS(fs fs.FS) (*AppBundle, error) {
	m, err := fs.Open(manifest.ManifestFileName)
	if err != nil {
		return nil, fmt.Errorf("could not open manifest: %w", err)
	}
	defer m.Close()

	man, err := manifest.LoadManifest(m)
	if err != nil {
		return nil, fmt.Errorf("could not load manifest: %w", err)
	}

	// Create app bundle struct
	return &AppBundle{
		Manifest: man,
		Source:   fs,
	}, nil
}

// FromDir translates a directory containing an app manifest and source
// into an AppBundle.
func FromDir(dir string) (*AppBundle, error) {
	return FromFS(os.DirFS(dir))
}

// LoadBundle loads a compressed archive into an AppBundle.
func LoadBundle(in io.Reader) (*AppBundle, error) {
	gzr, err := gzip.NewReader(in)
	if err != nil {
		return nil, fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gzr.Close()

	// read the entire tarball into memory so that we can seek
	// around it, and so that the underlying reader can be closed.
	var b bytes.Buffer
	io.Copy(&b, gzr)

	r := bytes.NewReader(b.Bytes())
	fs, err := tarfs.New(r)
	if err != nil {
		return nil, fmt.Errorf("creating tarfs: %w", err)
	}

	return FromFS(fs)
}
