package generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"

	"tidbyt.dev/pixlet/manifest"
)

const (
	appsDir      = "apps"
	manifestName = "manifest.yaml"
)

// AppType defines the type of app to generate using this package. There are
// several types of apps that are defined slightly differently. The ideal state
// is one type of app no matter where the app exists, but that's not the current
// reality.
type AppType int64

const (
	// Community represents an app that will be published in the community repo.
	Community AppType = iota
	// Local represents an app that is local and not meant to be published.
	Local
	//LocalCustom is exactly like Local but puts it in a differnt folder/name
	LocalCustom
	// Internal represents a Tidbyt internal app.
	Internal
)

//go:embed templates/source.star.tmpl
var starSource string

// Generator provides a structure for generating apps.
type Generator struct {
	starTmpl *template.Template
	appType  AppType
	root     string
}

type appsDef struct {
	Imports  []string
	Packages []string
}

// NewGenerator creates an instantiated generator with the templates parsed.
func NewGenerator(appType AppType, root string) (*Generator, error) {
	starTmpl, err := template.New("star").Parse(starSource)
	if err != nil {
		return nil, err
	}

	return &Generator{
		starTmpl: starTmpl,
		appType:  appType,
		root:     root,
	}, nil
}

// GenerateApp creates the base app starlark, go package, and updates the app
// list.
func (g *Generator) GenerateApp(app *manifest.Manifest) (string, error) {
	if g.appType == Community || g.appType == Internal || g.appType==LocalCustom {
		err := g.createDir(app)
		if err != nil {
			return "", err
		}
	}

	err := g.writeManifest(app)
	if err != nil {
		return "", err
	}

	return g.generateStarlark(app)
}

// RemoveApp removes an app from the apps directory.
func (g *Generator) RemoveApp(app *manifest.Manifest) error {
	return g.removeDir(app)
}

func (g *Generator) createDir(app *manifest.Manifest) error {
	// p := path.Join(g.root, appsDir, app.PackageName)
	var p string
	if g.appType==LocalCustom {
		p = g.root //appDir and/or PackageName is already accounted for here
	} else {
		p = path.Join(g.root, appsDir, app.PackageName)
	}
	return os.MkdirAll(p, os.ModePerm)
}

func (g *Generator) removeDir(app *manifest.Manifest) error {
	// p := path.Join(g.root, appsDir, app.PackageName)
	var p string
	if g.appType==LocalCustom {
		p = g.root //appDir and/or PackageName is already accounted for here
	} else {
		p = path.Join(g.root, appsDir, app.PackageName)
	}
	return os.RemoveAll(p)
}

func (g *Generator) writeManifest(app *manifest.Manifest) error {
	var p string
	switch g.appType {
	case Community, Internal:
		p = path.Join(g.root, appsDir, app.PackageName, manifestName)
	default:
		p = path.Join(g.root, manifestName)
	}

	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("couldn't create manifest file: %w", err)
	}
	defer f.Close()

	return app.WriteManifest(f)
}

func (g *Generator) generateStarlark(app *manifest.Manifest) (string, error) {
	var p string
	switch g.appType {
	case Community, Internal:
		p = path.Join(g.root, appsDir, app.PackageName, app.FileName)
	default:
		p = path.Join(g.root, app.FileName)
	}

	file, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = g.starTmpl.Execute(file, app)
	if err != nil {
		return "", err
	}

	return p, nil
}
