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
	"slices"

	"tidbyt.dev/pixlet/manifest"
	"tidbyt.dev/pixlet/runtime"
)

type WriteOption interface{}

type withoutRuntimeOption struct{}

// WithoutRuntime is a WriteOption that can be used to write the bundle without
// using the runtime to determine the files to include in the bundle. Instead,
// all files in the source FS will be included in the bundle.
//
// This is useful when writing a bundle that is known not to contain any
// unnecessary files, when loading and rewriting a bundle that was already
// tree-shaken, or when loading the entire runtime is not possible for
// performance or security reasons.
func WithoutRuntime() WriteOption {
	return withoutRuntimeOption{}
}

// WriteBundleToPath is a helper to be able to write the bundle to a provided
// directory.
func (b *AppBundle) WriteBundleToPath(dir string, opts ...WriteOption) error {
	path := filepath.Join(dir, AppBundleName)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file for bundle: %w", err)
	}
	defer f.Close()

	return b.WriteBundle(f, opts...)
}

// WriteBundle writes a compressed archive to the provided writer.
func (ab *AppBundle) WriteBundle(out io.Writer, opts ...WriteOption) error {
	var bundleFiles []string

	if slices.Contains(opts, WithoutRuntime()) {
		// we can't use the runtime to determine the files to include in the
		// bundle, so we'll just include everything in the source FS.
		err := fs.WalkDir(ab.Source, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walking directory: %w", err)
			}
			if !d.IsDir() {
				bundleFiles = append(bundleFiles, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("walking source FS: %w", err)
		}
	} else {
		// we don't want to naively write the entire source FS to the tarball,
		// since it could contain a lot of extraneous files. instead, run the
		// applet and interrogate it for the files it needs to include in the
		// bundle.
		app, err := runtime.NewAppletFromFS(ab.Manifest.ID, ab.Source, runtime.WithPrintDisabled())
		if err != nil {
			return fmt.Errorf("loading applet for bundling: %w", err)
		}
		bundleFiles = app.PathsForBundle()
	}

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
