// Package bundle provides primitives for bundling apps for portability.
package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	Source   []byte
	Manifest *manifest.Manifest
}

// InitFromPath translates a directory containing an app manifest and source
// into an AppBundle.
func InitFromPath(dir string) (*AppBundle, error) {
	// Load manifest
	path := filepath.Join(dir, manifest.ManifestFileName)
	m, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open manifest: %w", err)
	}
	defer m.Close()

	man, err := manifest.LoadManifest(m)
	if err != nil {
		return nil, fmt.Errorf("could not load manifest: %w", err)
	}

	// Load source
	path = filepath.Join(dir, man.FileName)
	s, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open app source: %w", err)
	}
	defer s.Close()

	src, err := io.ReadAll(s)
	if err != nil {
		return nil, fmt.Errorf("could not read app source: %w", err)
	}

	// Create app bundle struct
	return &AppBundle{
		Manifest: man,
		Source:   src,
	}, nil
}

// LoadBundle loads a compressed archive into an AppBundle.
func LoadBundle(in io.Reader) (*AppBundle, error) {
	gzr, err := gzip.NewReader(in)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	ab := &AppBundle{}

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			// If there are no more files in the bundle, validate and return it.
			if ab.Manifest == nil {
				return nil, fmt.Errorf("could not find manifest in archive")
			}
			if ab.Source == nil {
				return nil, fmt.Errorf("could not find source in archive")
			}
			return ab, nil
		case err != nil:
			// If there is an error, return immediately.
			return nil, fmt.Errorf("could not read archive: %w", err)
		case header == nil:
			// If for some reason we end up with a blank header, continue to the
			// next one.
			continue
		case header.Name == AppSourceName:
			// Load the app source.
			buff := make([]byte, header.Size)
			_, err := io.ReadFull(tr, buff)
			if err != nil {
				return nil, fmt.Errorf("could not read source from archive: %w", err)
			}
			ab.Source = buff
		case header.Name == manifest.ManifestFileName:
			// Load the app manifest.
			buff := make([]byte, header.Size)
			_, err := io.ReadFull(tr, buff)
			if err != nil {
				return nil, fmt.Errorf("could not read manifest from archive: %w", err)
			}

			man, err := manifest.LoadManifest(bytes.NewReader(buff))
			if err != nil {
				return nil, fmt.Errorf("could not load manifest: %w", err)
			}
			ab.Manifest = man
		}
	}
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
	// Setup writers.
	gzw := gzip.NewWriter(out)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Write manifest.
	buff := &bytes.Buffer{}
	err := ab.Manifest.WriteManifest(buff)
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

	// Write source.
	hdr = &tar.Header{
		Name: AppSourceName,
		Mode: 0600,
		Size: int64(len(ab.Source)),
	}
	err = tw.WriteHeader(hdr)
	if err != nil {
		return fmt.Errorf("could not write source header: %w", err)
	}
	_, err = tw.Write(ab.Source)
	if err != nil {
		return fmt.Errorf("could not write source to archive: %w", err)
	}

	return nil
}
