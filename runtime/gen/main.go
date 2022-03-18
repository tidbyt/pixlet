package main

// Generates starlark bindings for the pixlet/render package.
//
// Also produces widget documentation and extracts example snippets
// that can be run with docs/gen.go to produce images for the widget
// docs.

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"image/color"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/render/animation"
)

const DocumentationDirectory = "../docs/"

// Given a `reflect.Type` representing a pointer or slice, get the pointed-to or element type.
func decay(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		return t.Elem()
	}

	return t
}

// Given an `interface{}` return a `reflect.Type`, with pointer or slice unwrapped.
func toDecayedType(v interface{}) reflect.Type {
	return decay(reflect.TypeOf(v))
}

type Package struct {
	Name           string
	Directory      string
	ImportPath     string
	HeaderTemplate string
	TypeTemplate   string
	CodePath       string
	DocTemplate    string
	DocPath        string
	GoRootName     string
	GoWidgetName   string
	Types          []reflect.Value
}

// A list of packages and their types to generate code and documentation for.
var Packages = []Package{
	{
		Name:           "render",
		Directory:      "../render",
		ImportPath:     "tidbyt.dev/pixlet/render",
		HeaderTemplate: "./gen/header/render.tmpl",
		TypeTemplate:   "./gen/type.tmpl",
		CodePath:       "./modules/render_runtime/generated.go",
		DocTemplate:    "./gen/docs/render.tmpl",
		DocPath:        "../docs/widgets.md",
		GoRootName:     "Root",
		GoWidgetName:   "Widget",
		Types: []reflect.Value{
			reflect.ValueOf(new(render.Root)),
			reflect.ValueOf(new(render.Text)),
			reflect.ValueOf(new(render.Image)),
			reflect.ValueOf(new(render.Row)),
			reflect.ValueOf(new(render.Column)),
			reflect.ValueOf(new(render.Stack)),
			reflect.ValueOf(new(render.Padding)),
			reflect.ValueOf(new(render.Box)),
			reflect.ValueOf(new(render.Circle)),
			reflect.ValueOf(new(render.Marquee)),
			reflect.ValueOf(new(render.Animation)),
			reflect.ValueOf(new(render.WrappedText)),
		},
	},
	{
		Name:           "animation",
		Directory:      "../render/animation",
		ImportPath:     "tidbyt.dev/pixlet/render/animation",
		HeaderTemplate: "./gen/header/animation.tmpl",
		TypeTemplate:   "./gen/type.tmpl",
		CodePath:       "./modules/animation_runtime/generated.go",
		DocTemplate:    "./gen/docs/animation.tmpl",
		DocPath:        "../docs/animation.md",
		GoRootName:     "render_runtime.Root",
		GoWidgetName:   "render_runtime.Widget",
		Types: []reflect.Value{
			reflect.ValueOf(new(animation.AnimatedPositioned)),
		},
	},
}

// Defines how to generate code and documentation for type.
type Type struct {
	GoType        string
	DocType       string
	TemplatePath  string
	GenerateField bool
}

// A map of Go types to an `Attribute` definition.
var TypeMap = map[reflect.Type]Type{
	// Primitive types
	toDecayedType(new(string)): {
		GoType:       "starlark.String",
		DocType:      "str",
		TemplatePath: "./gen/attr/string.tmpl",
	},
	toDecayedType(new(int)): {
		GoType:       "starlark.Int",
		DocType:      "int",
		TemplatePath: "./gen/attr/int.tmpl",
	},
	toDecayedType(new(int32)): {
		GoType:       "starlark.Int",
		DocType:      "int",
		TemplatePath: "./gen/attr/int32.tmpl",
	},
	toDecayedType(new(float64)): {
		GoType:       "starlark.Value",
		DocType:      "float / int",
		TemplatePath: "./gen/attr/float.tmpl",
	},
	toDecayedType(new(bool)): {
		GoType:       "starlark.Bool",
		DocType:      "bool",
		TemplatePath: "./gen/attr/bool.tmpl",
	},

	// Render types
	toDecayedType(new(render.Insets)): {
		GoType:       "starlark.Value",
		DocType:      "int / (int, int, int, int)",
		TemplatePath: "./gen/attr/insets.tmpl",
	},
	toDecayedType(new(render.Widget)): {
		GoType:       "starlark.Value",
		DocType:      "Widget",
		TemplatePath: "./gen/attr/child.tmpl",
	},
	toDecayedType(new([]render.Widget)): {
		GoType:       "*starlark.List",
		DocType:      "[Widget]",
		TemplatePath: "./gen/attr/children.tmpl",
	},
	toDecayedType(new(color.Color)): {
		GoType:        "starlark.String",
		DocType:       `color`,
		TemplatePath:  "./gen/attr/color.tmpl",
		GenerateField: true,
	},

	// Animation types
	toDecayedType(new(animation.Curve)): {
		GoType:       "starlark.Value",
		DocType:      `str / function`,
		TemplatePath: "./gen/attr/curve.tmpl",
	},
}

// Defines a generated "Go to Starlark" attribute.
// This definition is passed to the templating engine.
type GeneratedAttr struct {
	GoName        string
	GoType        string
	GoWidgetName  string
	StarlarkName  string
	GenerateField bool
	IsRequired    bool
	IsReadOnly    bool

	// Template and generated code for handling this attribute.
	Template *template.Template
	Code     string

	// Documentation for this attribute.
	Documentation string
	DocType       string
}

// Defines a generated "Go to Starlark" binding type.
// This definition is passed to the templating engine.
type GeneratedType struct {
	GoName            string
	GoNameWithPackage string
	GoRootName        string
	GoWidgetName      string
	Attributes        []*GeneratedAttr
	HasSize           bool
	HasInit           bool
	Documentation     string
	Examples          []string
}

func nilOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Given a `reflect.Value`, return all its fields, including fields of anonymous composed types.
func allFields(val reflect.Value) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		t := typ.Field(i)
		v := val.Field(i)

		if t.Anonymous && t.Type.Kind() == reflect.Struct {
			fields = append(fields, allFields(v)...)
		} else {
			fields = append(fields, t)
		}
	}

	return fields
}

// Given a `reflect.StructField`, return a `GeneratedAttr` parse its `starlark:` field tag.
func toGeneratedAttribute(typ reflect.Type, field reflect.StructField) (*GeneratedAttr, error) {
	result := &GeneratedAttr{
		GoName:       field.Name,
		StarlarkName: strings.ToLower(field.Name),
	}

	// Fields can be tagged `starlark:"<name>[<param>...]"` to control the attribute name in Starlark.
	//
	// Additional supported flags:
	//   * "required" - field is required on instantiation
	//   * "readonly" - field is read-only, and not passed to constructor
	//
	if tag, ok := field.Tag.Lookup("starlark"); ok {
		attrs := strings.Split(tag, ",")
		if len(attrs) == 0 {
			return nil, fmt.Errorf("%s.%s has invalid tag: '%s'", typ.Name(), field.Name, tag)
		}

		result.StarlarkName = strings.TrimSpace(attrs[0])

		for _, attr := range attrs[1:] {
			attr = strings.TrimSpace(attr)
			if attr == "required" {
				result.IsRequired = true
			} else if attr == "readonly" {
				result.IsReadOnly = true
			} else {
				return nil, fmt.Errorf("%s.%s has unsupported tag attribute: '%s'", typ.Name(), field.Name, attr)
			}
		}
	}

	if result.StarlarkName == "" {
		result.StarlarkName = strings.ToLower(field.Name)
	}

	return result, nil
}

func toGeneratedType(pkg Package, val reflect.Value) (*GeneratedType, error) {
	result := &GeneratedType{}

	typ := val.Type()

	if decay(typ) == toDecayedType(new(render.Root)) {
		result.GoRootName = pkg.GoRootName
	}

	if typ.ConvertibleTo(toDecayedType(new(render.Widget))) {
		result.GoWidgetName = pkg.GoWidgetName
	}

	if typ.ConvertibleTo(toDecayedType(new(render.WidgetStaticSize))) {
		result.HasSize = true
	}

	if typ.ConvertibleTo(toDecayedType(new(render.WidgetWithInit))) {
		result.HasInit = true
	}

	// Unwrap any pointer types.
	val = reflect.Indirect(val)
	typ = val.Type()

	if val.Kind() != reflect.Struct {
		panic("type is neither struct nor pointer to struct, wtf?")
	}

	result.GoName = typ.Name()
	result.GoNameWithPackage = typ.String()

	for _, field := range allFields(val) {
		if field.PkgPath != "" || field.Anonymous {
			// Field is not an exposed attribute
			continue
		}

		if attr, err := toGeneratedAttribute(typ, field); err == nil {
			result.Attributes = append(result.Attributes, attr)

			if t, ok := TypeMap[field.Type]; ok {
				attr.GoType = t.GoType
				attr.GoWidgetName = pkg.GoWidgetName
				attr.DocType = t.DocType
				attr.Template = loadTemplate("attr", t.TemplatePath)
				attr.GenerateField = t.GenerateField
			} else {
				return nil, fmt.Errorf("%s.%s has unsupported type", typ.Name(), field.Name)
			}
		} else {
			return nil, err
		}
	}

	// Reorder attributes so that required fields appear first.
	sort.SliceStable(result.Attributes, func(i, j int) bool {
		return result.Attributes[i].IsRequired && !result.Attributes[j].IsRequired
	})

	return result, nil
}

func loadTemplate(name, path string) *template.Template {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	content, err := ioutil.ReadFile(path)
	nilOrPanic(err)

	template, err := template.New(name).Funcs(funcMap).Parse(string(content))
	nilOrPanic(err)

	return template
}

func renderTemplateToFile(tmpl *template.Template, data interface{}, path string) {
	outf, err := os.Create(path)
	nilOrPanic(err)
	defer outf.Close()
	err = tmpl.Execute(outf, data)
	nilOrPanic(err)
}

func renderTemplateToBuffer(tmpl *template.Template, data interface{}, buf *bytes.Buffer) {
	err := tmpl.Execute(buf, data)
	nilOrPanic(err)
}

func renderTemplateToString(tmpl *template.Template, data interface{}) string {
	var buf bytes.Buffer
	renderTemplateToBuffer(tmpl, data, &buf)
	return string(buf.Bytes())
}

func attachDocs(pkg Package, types []*GeneratedType) {
	// Parse all .go files in pixlet/render packages and extract all type doc comments
	fset := token.NewFileSet()
	docs := map[string]string{}

	astPkgs, err := parser.ParseDir(fset, pkg.Directory, nil, parser.ParseComments)
	nilOrPanic(err)
	pkgDoc := doc.New(astPkgs[pkg.Name], pkg.ImportPath, 0)
	nilOrPanic(err)
	for _, type_ := range pkgDoc.Types {
		docs[type_.Name] = type_.Doc
	}

	// These match our attribute docs and example blocks
	docRe, err := regexp.Compile(`(?m)^DOC\(([^)]+)\): +(.+)$`)
	nilOrPanic(err)
	exampleRe, err := regexp.Compile(`(?s)EXAMPLE BEGIN(.*?)EXAMPLE END`)
	nilOrPanic(err)

	for _, type_ := range types {
		// Widget doc is full comment sans attribute docs and examples
		type_.Documentation = strings.TrimSpace(string(
			docRe.ReplaceAllString(
				exampleRe.ReplaceAllString(docs[type_.GoName], ""),
				"",
			),
		))

		// Attribute docs
		attrDocs := map[string]string{}
		for _, group := range docRe.FindAllStringSubmatch(docs[type_.GoName], -1) {
			attrDocs[group[1]] = group[2]
		}
		for _, attr := range type_.Attributes {
			attr.Documentation = attrDocs[attr.GoName]
		}

		// Examples
		examples := []string{}
		for _, group := range exampleRe.FindAllStringSubmatch(docs[type_.GoName], -1) {
			examples = append(examples, strings.TrimSpace(group[1]))
		}
		type_.Examples = examples

	}
}

func generateCode(pkg Package, types []*GeneratedType) {
	// First render templates for each attribute.
	for _, type_ := range types {
		for _, attr := range type_.Attributes {
			attr.Code = renderTemplateToString(attr.Template, attr)
		}
	}

	// Then render templates for the header and for each type.
	headerTmpl := loadTemplate("header", pkg.HeaderTemplate)
	typeTmpl := loadTemplate("type", pkg.TypeTemplate)

	outf, err := os.Create(pkg.CodePath)
	nilOrPanic(err)
	defer outf.Close()

	var buf bytes.Buffer
	renderTemplateToBuffer(headerTmpl, types, &buf)

	for _, typ := range types {
		renderTemplateToBuffer(typeTmpl, typ, &buf)
	}

	// Format and write the source to disk.
	source, err := format.Source(buf.Bytes())
	nilOrPanic(err)
	outf.Write(source)
}

func generateDocs(pkg Package, types []*GeneratedType) {
	template := loadTemplate("docs", pkg.DocTemplate)

	renderTemplateToFile(template, types, pkg.DocPath)

	for _, typ := range types {
		for i, example := range typ.Examples {
			err := ioutil.WriteFile(
				fmt.Sprintf("%s/%s_%d.star", DocumentationDirectory, typ.GoName, i),
				[]byte(example),
				0644)
			nilOrPanic(err)
		}
	}
}

func main() {
	// Generate code and documentation for each package.
	for _, pkg := range Packages {
		types := []*GeneratedType{}

		for _, typ := range pkg.Types {
			if result, err := toGeneratedType(pkg, typ); err == nil {
				types = append(types, result)
			} else {
				panic(err)
			}
		}

		sort.SliceStable(types, func(i, j int) bool {
			return types[i].GoName < types[j].GoName
		})

		attachDocs(pkg, types)
		generateCode(pkg, types)
		generateDocs(pkg, types)
	}
}
