// Package manifest provides structures and primitives to define apps.
package manifest

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

const ManifestFileName = "manifest.yaml"

// Manifest is a structure to define a starlark applet for Tidbyt in Go.
type Manifest struct {
	// ID is the unique identifier of this app. It has to be globally unique,
	// which means it cannot conflict with any of our private apps.
	ID string `json:"id" yaml:"id"`
	// Name is the name of the applet. Ex. "Fuzzy Clock"
	Name string `json:"name" yaml:"name"`
	// Summary is the short form of what this applet does. Ex. "Human readable
	// time".
	Summary string `json:"summary" yaml:"summary"`
	// Desc is the long form of what this applet does. Ex. "Display the time in
	// a groovy, human-readable way."
	Desc string `json:"desc" yaml:"desc"`
	// Author is the person or organization who contributed this applet. Ex,
	// "Max Timkovich"
	Author string `json:"author" yaml:"author"`
	// FileName is the name of the starlark source file.
	FileName string `json:"fileName" yaml:"fileName"`
	// PackageName is the name of the go package where this app lives.
	PackageName string `json:"packageName" yaml:"packageName"`
	// Source is the starlark source code for this applet using the go `embed`
	// module.
	Source []byte `json:"-" yaml:"-"`
}

// LoadManifest reads a manifest from an io.Reader, with the most common reader
// being a file from os.Open. It returns a manifest or an error if it could not
// be parsed.
func LoadManifest(r io.Reader) (*Manifest, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read manifest: %w", err)
	}

	manifest := &Manifest{}
	err = yaml.Unmarshal(b, manifest)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal manifest: %w", err)
	}

	return manifest, nil
}

// WriteManifest writes a manifest to the supplied writer, with the most common
// writer being a file. Any issue will return an error.
func (m *Manifest) WriteManifest(w io.Writer) error {
	b, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal manifest: %w", err)
	}

	fmt.Fprintf(w, "---\n")

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("could not write manifest: %w", err)
	}

	return nil
}

// Validate ensures all fields of the manifest are valid and returns an error
// if they are not.
func (m Manifest) Validate() error {
	err := ValidateID(m.ID)
	if err != nil {
		return err
	}

	err = ValidateName(m.Name)
	if err != nil {
		return err
	}

	err = ValidateSummary(m.Summary)
	if err != nil {
		return err
	}

	err = ValidateDesc(m.Desc)
	if err != nil {
		return err
	}

	err = ValidateAuthor(m.Author)
	if err != nil {
		return err
	}

	err = ValidateFileName(m.FileName)
	if err != nil {
		return err
	}

	err = ValidatePackageName(m.PackageName)
	if err != nil {
		return err
	}

	return nil
}

// GeneratePackageName creates a suitable go package name from an app name.
func GeneratePackageName(name string) string {
	packageName := strings.ReplaceAll(name, "-", "")
	packageName = strings.ReplaceAll(packageName, "_", "")
	return strings.ToLower(strings.Join(strings.Fields(packageName), ""))
}

// GenerateID creates a suitable ID from an app name.
func GenerateID(name string) string {
	id := strings.ReplaceAll(name, "_", "-")
	return strings.ToLower(strings.Join(strings.Fields(id), "-"))
}

// GenerateFileName creates a suitable file name for the starlark source.
func GenerateFileName(name string) string {
	fileName := strings.ReplaceAll(name, "-", "_")
	return strings.ToLower(strings.Join(strings.Fields(fileName), "_")) + ".star"
}
