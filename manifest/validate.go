package manifest

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// Our longest app name to date. This can be updated, but it will need to
	// be tested in the mobile app.
	MaxNameLength = 17

	// Our longest app summary to date. This can be updated, but it will need to
	// be tested in the mobile app.
	MaxSummaryLength = 27

	dash       = '-'
	underscore = '_'
)

var punctuation []string = []string{
	".",
	"!",
	"?",
}

var titleCaser cases.Caser

func init() {
	titleCaser = cases.Title(language.English, cases.NoLower)
}

// ValidateName ensures the app name provided adheres to the standards for app
// names. We're picky here because these will display in the Tidbyt mobile app
// and need to display properly.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if name != titleCase(name) {
		return fmt.Errorf("'%s' should be title case, 'Fuzzy Clock' for example", name)
	}

	if len(name) > MaxNameLength {
		return fmt.Errorf("app names need to be less then %d characters", MaxNameLength)
	}

	return nil
}

// ValidateSummary ensures the app summary provided adheres to the standards
// for app summaries. We're picky here because these will display in the Tidbyt
// mobile app and need to display properly.
func ValidateSummary(summary string) error {
	if summary == "" {
		return fmt.Errorf("summary cannot be empty")
	}

	if len(summary) > MaxSummaryLength {
		return fmt.Errorf("app summaries need to be less then %d characters", MaxSummaryLength)
	}

	for _, punct := range punctuation {
		if strings.HasSuffix(summary, punct) {
			return fmt.Errorf("app summaries should not end in punctuation")
		}
	}

	words := strings.Split(summary, " ")
	if len(words) > 0 && words[0] != titleCaser.String(words[0]) {
		return fmt.Errorf("app summaries should start with an uppercased character")
	}

	return nil
}

// ValidateDesc ensures the app description provided adheres to the standards
// for app descriptions. We're picky here because these will display in the
// Tidbyt mobile app and need to display properly.
func ValidateDesc(desc string) error {
	if desc == "" {
		return fmt.Errorf("desc cannot be empty")
	}

	found := false
	for _, punct := range punctuation {
		if strings.HasSuffix(desc, punct) {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("app descriptions should end in punctuation")
	}

	words := strings.Split(desc, " ")
	if len(words) > 0 && words[0] != titleCaser.String(words[0]) {
		return fmt.Errorf("app descriptions should start with an uppercased character")
	}

	return nil
}

// ValidateAuthor ensures the app author provided adheres to the standards
// for app author. We're picky here because these will display in the
// Tidbyt mobile app and need to display properly.
func ValidateAuthor(author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	// I don't know what validation where need here just yet. We're going to
	// have to eyeball it in pull requests until we get a sense of what doesn't
	// work.
	return nil
}

func ValidatePackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("package names cannot be empty")
	}

	if packageName != strings.ToLower(packageName) {
		return fmt.Errorf("package names should be lower case")
	}

	for _, r := range packageName {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return fmt.Errorf("package names can only contain letters, numbers, or an underscore character")
		}
	}
	return nil
}

// ValidateFileName ensures the file name appears appropriately for starlark
// source code.
func ValidateFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("fileName cannot be empty")
	}

	if !strings.HasSuffix(fileName, ".star") {
		return fmt.Errorf("file names should end in .star: '%s'", fileName)
	}

	testName := strings.TrimSuffix(fileName, ".star")

	if testName != strings.ToLower(testName) {
		return fmt.Errorf("file names should be lower case")
	}

	for _, r := range testName {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == underscore) {
			return fmt.Errorf("file names can only contain letters, numbers, or an underscore character")
		}
	}

	return nil
}

// ValidateID ensures the id will parse when we go to add it to our database
// internally.
func ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if id != strings.ToLower(id) {
		return fmt.Errorf("ids should be lower case, %s != %s", id, strings.ToLower(id))
	}

	for _, r := range id {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == dash) {
			return fmt.Errorf("ids can only contain letters, numbers, or a dash character")
		}
	}

	return nil
}

func titleCase(input string) string {
	words := strings.Split(input, " ")
	smallwords := " a an on the to of "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") && word != string(word[0]) {
			words[index] = word
		} else {
			words[index] = titleCaser.String(word)
		}
	}

	return strings.Join(words, " ")
}
