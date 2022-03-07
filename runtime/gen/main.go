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

const StarlarkHeaderTemplate = "./gen/header.tmpl"
const StarlarkWidgetTemplate = "./gen/widget.tmpl"
const DocumentationTemplate = "./gen/docs.tmpl"
const RenderDirectory = "../render"
const CodeOut = "./generated.go"
const DocumentationOut = "../docs/widgets.md"
const DocumentationDirectory = "../docs/"

var RenderWidgets = []render.Widget{
	&render.Text{},
	&render.Image{},
	render.Row{},
	render.Column{},
	render.Stack{},
	render.Padding{},
	render.Box{},
	render.Circle{},
	render.Marquee{},
	render.Animation{},
	render.WrappedText{},
	animation.AnimatedPositioned{},
}

// Defines the starlark version of a render.Widget
type Attribute struct {
	Render        string
	Starlark      string
	Required      bool
	ReadOnly      bool
	Type          string
	Documentation string
}

type StarlarkWidget struct {
	Name          string
	FullName      string
	AttrAll       []*Attribute
	AttrString    []*Attribute
	AttrInt       []*Attribute
	AttrColor     []*Attribute
	AttrCurve     []*Attribute
	AttrBool      []*Attribute
	AttrChildren  []*Attribute
	AttrChild     []*Attribute
	AttrInsets    []*Attribute
	HasSize       bool
	HasPtrRcvr    bool
	RequiresInit  bool
	Documentation string
	Examples      []string
}

type StarlarkHeader struct {
	Widget []StarlarkWidget
}

func nilOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func starlarkWidgetFromRenderWidget(w render.Widget) *StarlarkWidget {
	sw := StarlarkWidget{}

	val := reflect.ValueOf(w)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
		sw.HasPtrRcvr = true
	}
	if val.Kind() != reflect.Struct {
		panic("widget is neither struct nor pointer to struct, wtf?")
	}

	typ := val.Type()

	sw.Name = typ.Name()
	sw.FullName = typ.String()

	if _, hasSize := w.(render.WidgetStaticSize); hasSize {
		sw.HasSize = true
	}

	if _, requiresInit := w.(render.WidgetWithInit); requiresInit {
		sw.RequiresInit = true
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.PkgPath != "" || field.Anonymous {
			// Field is not an exposed attribute
			continue
		}

		// Widget fields can be tagged `starlark:"<name>[<param>...]"` to
		// control attribute name in starlark.
		//
		// Additional supported flags:
		// "required" - field is required on instantiation
		// "readonly" - field is read-only, and not passed to constructor
		attr := &Attribute{
			Render: field.Name,
		}
		fieldTag, ok := field.Tag.Lookup("starlark")
		if ok {
			tag := strings.Split(fieldTag, ",")
			attr.Starlark = strings.TrimSpace(tag[0])
			for _, t := range tag[1:] {
				t = strings.TrimSpace(t)
				if t == "required" {
					attr.Required = true
				} else if t == "readonly" {
					attr.ReadOnly = true
				} else {
					panic(fmt.Sprintf(
						"%s.%s has unsupported tag: '%s'",
						typ.Name(), field.Name, tag[1],
					))
				}
			}
		}
		if attr.Starlark == "" {
			attr.Starlark = strings.ToLower(field.Name)
		}

		sw.AttrAll = append(sw.AttrAll, attr)

		switch field.Type.Kind() {
		case reflect.Int:
			sw.AttrInt = append(sw.AttrInt, attr)
			attr.Type = "int"
		case reflect.String:
			sw.AttrString = append(sw.AttrString, attr)
			attr.Type = "str"
		case reflect.Bool:
			sw.AttrBool = append(sw.AttrBool, attr)
			attr.Type = "bool"
		case reflect.Slice:
			sw.AttrChildren = append(sw.AttrChildren, attr)
			attr.Type = "list"
		case reflect.Struct:
			insetsType := reflect.TypeOf((*render.Insets)(nil)).Elem()
			if field.Type == insetsType {
				sw.AttrInsets = append(sw.AttrInsets, attr)
				attr.Type = "insets"
			} else {
				panic(fmt.Sprintf(
					"%s.%s has unsupported type",
					typ.Name(), field.Name,
				))
			}
		case reflect.Interface:
			colorType := reflect.TypeOf((*color.Color)(nil)).Elem()
			curveType := reflect.TypeOf((*animation.Curve)(nil)).Elem()
			widgetType := reflect.TypeOf((*render.Widget)(nil)).Elem()

			if field.Type.Implements(colorType) {
				sw.AttrColor = append(sw.AttrColor, attr)
				attr.Type = "color"
			} else if field.Type.Implements(curveType) {
				sw.AttrCurve = append(sw.AttrCurve, attr)
				attr.Type = "curve"
			} else if field.Type.Implements(widgetType) {
				sw.AttrChild = append(sw.AttrChild, attr)
				attr.Type = "Widget"
			} else {
				panic(fmt.Sprintf(
					"%s.%s has unsupported type",
					typ.Name(), field.Name,
				))
			}
		}
	}

	// Reorder AttrAll so that required fields appear first
	sort.SliceStable(sw.AttrAll, func(i, j int) bool {
		return sw.AttrAll[i].Required && !sw.AttrAll[j].Required
	})

	return &sw
}

func attachWidgetDocs(widgets []*StarlarkWidget) {

	// Parse all .go files in pixlet/render and extract all type doc comments
	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, RenderDirectory, nil, parser.ParseComments)
	nilOrPanic(err)
	pkg := doc.New(astPkgs["render"], "tidbyt.dev/pixlet/render", 0)
	nilOrPanic(err)
	docs := map[string]string{}
	for _, type_ := range pkg.Types {
		docs[type_.Name] = type_.Doc
	}

	// These match our attribute docs and example blocks
	docRe, err := regexp.Compile(`(?m)^DOC\(([^)]+)\): +(.+)$`)
	nilOrPanic(err)
	exampleRe, err := regexp.Compile(`(?s)EXAMPLE BEGIN(.*?)EXAMPLE END`)
	nilOrPanic(err)

	for _, widget := range widgets {
		// Widget doc is full comment sans attribute docs and examples
		widget.Documentation = strings.TrimSpace(string(
			docRe.ReplaceAllString(
				exampleRe.ReplaceAllString(docs[widget.Name], ""),
				"",
			),
		))

		// Attribute docs
		attrDocs := map[string]string{}
		for _, group := range docRe.FindAllStringSubmatch(docs[widget.Name], -1) {
			attrDocs[group[1]] = group[2]
		}
		for _, attr := range widget.AttrAll {
			attr.Documentation = attrDocs[attr.Render]
		}

		// Examples
		examples := []string{}
		for _, group := range exampleRe.FindAllStringSubmatch(docs[widget.Name], -1) {
			examples = append(examples, strings.TrimSpace(group[1]))
		}
		widget.Examples = examples

	}
}

func generateCode(widgets []*StarlarkWidget) {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	headerTemplateContent, err := ioutil.ReadFile(StarlarkHeaderTemplate)
	nilOrPanic(err)

	headerTemplate, err := template.New("header").Funcs(funcMap).Parse(string(headerTemplateContent))
	nilOrPanic(err)

	widgetTemplateContent, err := ioutil.ReadFile(StarlarkWidgetTemplate)
	nilOrPanic(err)

	widgetTemplate, err := template.New("widget").Funcs(funcMap).Parse(string(widgetTemplateContent))
	nilOrPanic(err)

	outf, err := os.Create(CodeOut)
	nilOrPanic(err)
	defer outf.Close()

	var buf bytes.Buffer
	err = headerTemplate.Execute(&buf, widgets)
	nilOrPanic(err)

	for _, data := range widgets {
		err = widgetTemplate.Execute(&buf, data)
		nilOrPanic(err)
	}

	formatted, err := format.Source(buf.Bytes())
	nilOrPanic(err)
	outf.Write(formatted)
}

func generateDocs(widgets []*StarlarkWidget) {
	funcMap := template.FuncMap{
		"attributeRow": func(attr *Attribute) string {
			var name, type_, descr string

			if attr.Required {
				name = fmt.Sprintf("**%s**", attr.Starlark)
			} else {
				name = attr.Starlark
			}

			switch attr.Type {
			case "str", "int", "bool", "list", "Widget":
				type_ = attr.Type
			case "color":
				type_ = "str"
			case "insets":
				type_ = "int / (int, int, int, int)"
			default:
				panic(fmt.Sprintf("bad type: %s", attr.Type))
			}

			descr = attr.Documentation

			return fmt.Sprintf("| %s | %s | %s |", name, type_, descr)
		},
	}

	docsTemplateContent, err := ioutil.ReadFile(DocumentationTemplate)
	nilOrPanic(err)

	docsTemplate, err := template.New("docs").Funcs(funcMap).Parse(string(docsTemplateContent))
	nilOrPanic(err)

	outf, err := os.Create(DocumentationOut)
	nilOrPanic(err)
	defer outf.Close()

	err = docsTemplate.Execute(outf, widgets)
	nilOrPanic(err)

	for _, widget := range widgets {
		for i, example := range widget.Examples {
			err = ioutil.WriteFile(
				fmt.Sprintf("%s/%s_%d.star", DocumentationDirectory, widget.Name, i),
				[]byte(example),
				0644)
			nilOrPanic(err)
		}
	}
}

func main() {
	widgets := []*StarlarkWidget{}
	for _, w := range RenderWidgets {
		widgets = append(widgets, starlarkWidgetFromRenderWidget(w))
	}

	sort.SliceStable(widgets, func(i, j int) bool {
		return widgets[i].Name < widgets[j].Name
	})

	attachWidgetDocs(widgets)
	generateCode(widgets)
	generateDocs(widgets)
}
