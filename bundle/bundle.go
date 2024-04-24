// Package bundle provides primitives for bundling apps for portability.
package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/quay/claircore/pkg/tarfs"

	"tidbyt.dev/pixlet/manifest"
	"tidbyt.dev/pixlet/runtime"
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

func fromFS(fs fs.FS) (*AppBundle, error) {
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

// InitFromPath translates a directory containing an app manifest and source
// into an AppBundle.
func InitFromPath(dir string) (*AppBundle, error) {
	return fromFS(os.DirFS(dir))
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

	return fromFS(fs)
}

// WriteBundleToPath is a helper to be able to write the bundle to a provided
// directory.
func (b *AppBundle) WriteBundleToPath(dir string) error {
	path := filepath.Join(dir, AppBundleName)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file for bundle: %w", err)
	}
	defer f.Close()

	return b.WriteBundle(f)
}

// WriteBundle writes a compressed archive to the provided writer.
func (ab *AppBundle) WriteBundle(out io.Writer) error {
	// we don't want to naively write the entire source FS to the tarball,
	// since it could contain a lot of extraneous files. instead, run the
	// applet and interrogate it for the files it needs to include in the
	// bundle.
	app, err := runtime.NewAppletFromFS(ab.Manifest.ID, ab.Source, runtime.WithPrintDisabled())
	if err != nil {
		return fmt.Errorf("loading applet for bundling: %w", err)
	}
	bundleFiles := app.PathsForBundle()

	// Setup writers.
	gzw := gzip.NewWriter(out)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Write manifest.
	buff := &bytes.Buffer{}
	err = ab.Manifest.WriteManifest(buff)
	if err != nil {
		return fmt.Errorf("could not write manifest to buffer: %w", err)
	}
	b := buff.Bytes()

	hdr := &tar.Header{
		Name: manifest.ManifestFileName,
		Mode: 0600,
		Size: int64(len(b)),
	}
	err = tw.WriteHeader(hdr)
	if err != nil {
		return fmt.Errorf("could not write manifest header: %w", err)
	}
	_, err = tw.Write(b)
	if err != nil {
		return fmt.Errorf("could not write manifest to archive: %w", err)
	}

	// write sources.
	for _, path := range bundleFiles {
		stat, err := fs.Stat(ab.Source, path)
		if err != nil {
			return fmt.Errorf("could not stat %s: %w", path, err)
		}

		hdr, err := tar.FileInfoHeader(stat, "")
		if err != nil {
			return fmt.Errorf("creating header for %s: %w", path, err)
		}
		hdr.Name = filepath.ToSlash(path)

		err = tw.WriteHeader(hdr)
		if err != nil {
			return fmt.Errorf("writing header for %s: %w", path, err)
		}

		if !stat.IsDir() {
			file, err := ab.Source.Open(path)
			if err != nil {
				return fmt.Errorf("opening file %s: %w", path, err)
			}

			written, err := io.Copy(tw, file)
			if err != nil {
				return fmt.Errorf("writing file %s: %w", path, err)
			} else if written != stat.Size() {
				return fmt.Errorf("did not write entire file %s: %w", path, err)
			}
		}
	}

	return nil
}
