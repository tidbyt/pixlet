// Package manifest provides structures and primitives to define apps.
package manifest

import "strings"

// Manifest is a structure to define a starlark applet for Tidbyt in Go.
type Manifest struct {
	// ID is the unique identifier of this app. It has to be globally unique,
	// which means it cannot conflict with any of our private apps.
	ID string `json:"id"`
	// Name is the name of the applet. Ex. "Fuzzy Clock"
	Name string `json:"name"`
	// Summary is the short form of what this applet does. Ex. "Human readable
	// time".
	Summary string `json:"summary"`
	// Desc is the long form of what this applet does. Ex. "Display the time in
	// a groovy, human-readable way."
	Desc string `json:"desc"`
	// Author is the person or organization who contributed this applet. Ex,
	// "Max Timkovich"
	Author string `json:"author"`
	// FileName is the name of the starlark source file.
	FileName string `json:"file_name"`
	// PackageName is the name of the go package where this app lives.
	PackageName string `json:"package_name"`
	// Source is the starlark source code for this applet using the go `embed`
	// module.
	Source []byte `json:"-"`
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
