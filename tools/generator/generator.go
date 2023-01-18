package generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"sort"
	"text/template"

	"tidbyt.dev/pixlet/manifest"
)

const (
	appsDir      = "apps"
	goExt        = ".go"
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
	// Internal represents a Tidbyt internal app.
	Internal
)

//go:embed templates/source.star.tmpl
var starSource string

//go:embed templates/source.go.tmpl
var goSource string

//go:embed templates/apps.go.tmpl
var appsSource string

// Generator provides a structure for generating apps.
type Generator struct {
	starTmpl     *template.Template
	goTmpl       *template.Template
	manifestTmpl *template.Template
	appsTmpl     *template.Template
	appType      AppType
	root         string
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

	goTmpl, err := template.New("go").Parse(goSource)
	if err != nil {
		return nil, err
	}

	manifestTmpl, err := template.New("manifest").Parse(goSource)
	if err != nil {
		return nil, err
	}

	appsTmpl, err := template.New("apps").Parse(appsSource)
	if err != nil {
		return nil, err
	}

	return &Generator{
		starTmpl:     starTmpl,
		goTmpl:       goTmpl,
		manifestTmpl: manifestTmpl,
		appsTmpl:     appsTmpl,
		appType:      appType,
		root:         root,
	}, nil
}

// GenerateApp creates the base app starlark, go package, and updates the app
// list.
func (g *Generator) GenerateApp(app *manifest.Manifest) (string, error) {
	if g.appType == Community {
		err := g.createDir(app)
		if err != nil {
			return "", err
		}

		err = g.generateGo(app)
		if err != nil {
			return "", err
		}

		err = g.updateApps()
		if err != nil {
			return "", err
		}
	}

	if g.appType == Local || g.appType == Community {
		err := g.writeManifest(app)
		if err != nil {
			return "", err
		}
	}

	return g.generateStarlark(app)
}

// RemoveApp removes an app from the apps directory.
func (g *Generator) RemoveApp(app *manifest.Manifest) error {
	err := g.removeDir(app)
	if err != nil {
		return err
	}

	return g.updateApps()
}

// UpdateApps generates the app list in apps.go.
func (g *Generator) UpdateApps() error {
	return g.updateApps()
}

func (g *Generator) createDir(app *manifest.Manifest) error {
	p := path.Join(g.root, appsDir, app.PackageName)
	return os.MkdirAll(p, os.ModePerm)
}

func (g *Generator) removeDir(app *manifest.Manifest) error {
	p := path.Join(g.root, appsDir, app.PackageName)
	return os.RemoveAll(p)
}

func (g *Generator) updateApps() error {
	imports := []string{
		"tidbyt.dev/community/" + appsDir + "/manifest",
	}
	packages := []string{}

	files, err := os.ReadDir(path.Join(g.root, appsDir))
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != "manifest" {
			imp := "tidbyt.dev/community/" + appsDir + "/" + f.Name()
			imports = append(imports, imp)
			packages = append(packages, f.Name())
		}
	}
	p := path.Join(g.root, appsDir, appsDir+goExt)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Strings(imports)
	sort.Strings(packages)

	a := &appsDef{
		Imports:  imports,
		Packages: packages,
	}

	return g.appsTmpl.Execute(file, a)
}

func (g *Generator) writeManifest(app *manifest.Manifest) error {
	var p string
	switch g.appType {
	case Community:
		p = path.Join(g.root, appsDir, app.PackageName, app.FileName)
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
	case Community:
		p = path.Join(g.root, appsDir, app.PackageName, app.FileName)
	case Internal:
		p = path.Join(g.root, appsDir, app.FileName)
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

func (g *Generator) generateGo(app *manifest.Manifest) error {
	p := path.Join(g.root, appsDir, app.PackageName, app.PackageName+goExt)

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.goTmpl.Execute(file, app)
}
