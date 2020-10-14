package runtime

import (
	"crypto/md5"
	"fmt"
	"image"
	"time"

	"github.com/harukasan/go-libwebp/webp"
	"github.com/pkg/errors"
	starlibbase64 "github.com/qri-io/starlib/encoding/base64"
	starlibjson "github.com/qri-io/starlib/encoding/json"
	starlibhttp "github.com/qri-io/starlib/http"
	starlibmath "github.com/qri-io/starlib/math"
	starlibre "github.com/qri-io/starlib/re"
	starlibtime "github.com/qri-io/starlib/time"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/tidbyt/pixlet/render"
)

type ModuleLoader func(*starlark.Thread, string) (starlark.StringDict, error)

const (
	WebPKMin                 = 0
	WebPKMax                 = 0
	DefaultScreenDelayMillis = 50
)

func init() {
	resolve.AllowFloat = true
	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowSet = true
	resolve.AllowRecursion = true
}

type Applet struct {
	Filename    string
	Id          string
	Globals     starlark.StringDict
	src         []byte
	loader      ModuleLoader
	predeclared starlark.StringDict
	main        *starlark.Function
}

type ImageFilter func(image.Image) (image.Image, error)

type Screens struct {
	roots  []render.Root
	images []image.Image
	delay  int32
}

func ScreensFromRoots(roots []render.Root) *Screens {
	screens := Screens{
		roots: roots,
		delay: DefaultScreenDelayMillis,
	}
	if len(roots) > 0 {
		screens.delay = roots[0].Delay
	}
	return &screens
}

func ScreensFromImages(images ...image.Image) *Screens {
	screens := Screens{
		images: images,
		delay:  DefaultScreenDelayMillis,
	}
	return &screens
}

// Renders a screen to WebP. Optionally pass filters for
// postprocessing each individual frame.
func (s *Screens) Render(filters ...ImageFilter) ([]byte, error) {

	if s.images == nil {
		s.images = render.PaintRoots(true, s.roots...)
	}

	if len(s.images) == 0 {
		return []byte{}, nil
	}

	images := s.images

	if len(filters) > 0 {
		images = []image.Image{}
		for _, im := range s.images {
			for _, f := range filters {
				imFiltered, err := f(im)
				if err != nil {
					return nil, err
				}
				im = imFiltered
			}
			images = append(images, im)
		}
	}

	bounds := images[0].Bounds()
	anim, err := webp.NewAnimationEncoder(
		bounds.Dx(),
		bounds.Dy(),
		WebPKMin,
		WebPKMax,
	)
	if err != nil {
		return nil, errors.Wrap(err, "initializing encoder")
	}
	defer anim.Close()

	frameDuration := time.Duration(s.delay) * time.Millisecond
	for _, im := range images {
		if err := anim.AddFrame(im, frameDuration); err != nil {
			return nil, errors.Wrap(err, "adding frame")
		}
	}

	buf, err := anim.Assemble()
	if err != nil {
		return nil, errors.Wrap(err, "encoding animation")
	}

	return buf, nil
}

func (a *Applet) thread() *starlark.Thread {
	return &starlark.Thread{
		Name: a.Id,
		Load: a.loadModule,
		Print: func(thread *starlark.Thread, msg string) {
			fmt.Printf("[%s] %s\n", a.Filename, msg)
		},
	}
}

// Loads an applet. The script filename is used as a descriptor only,
// and the actual code should be passed in src. Optionally also pass
// loader to make additional starlark modules available to the script.
func (a *Applet) Load(filename string, src []byte, loader ModuleLoader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while executing %s: %v", a.Filename, r)
		}
	}()

	a.Filename = filename
	a.loader = loader

	a.src = src

	a.Id = fmt.Sprintf("%s/%x", filename, md5.Sum(src))

	a.predeclared = starlark.StringDict{
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
	}

	globals, err := starlark.ExecFile(a.thread(), a.Filename, a.src, a.predeclared)
	if err != nil {
		return fmt.Errorf("starlark.ExecFile: %v", err)
	}
	a.Globals = globals

	mainFun, found := globals["main"]
	if !found {
		return fmt.Errorf("%s didn't export a main() function", filename)
	}
	main, ok := mainFun.(*starlark.Function)
	if !ok {
		return fmt.Errorf("%s exported a main() that is not function", filename)
	}
	a.main = main

	return nil
}

// Runs the applet's main function, passing it configuration as a
// starlark dict.
func (a *Applet) Run(config map[string]string) (screens *Screens, err error) {
	var args starlark.Tuple
	if a.main.NumParams() > 0 {
		starlarkConfig := starlark.NewDict(len(config))
		for k, v := range config {
			starlarkConfig.SetKey(
				starlark.String(k),
				starlark.String(v),
			)
		}
		args = starlark.Tuple{starlarkConfig}
	}

	returnValue, err := a.Call(a.main, args)
	if err != nil {
		return nil, err
	}

	var roots []render.Root
	if returnRoot, ok := returnValue.(Root); ok {
		roots = []render.Root{returnRoot.AsRenderRoot()}
	} else if returnList, ok := returnValue.(*starlark.List); ok {
		roots = make([]render.Root, returnList.Len())
		iter := returnList.Iterate()
		defer iter.Done()
		i := 0
		var listVal starlark.Value
		for iter.Next(&listVal) {
			if listValRoot, ok := listVal.(Root); ok {
				roots[i] = listValRoot.AsRenderRoot()
			} else {
				return nil, fmt.Errorf(
					"expected app implementation to return Root(s) but found: %s (at index %d)",
					listVal.Type(),
					i,
				)
			}
			i++
		}
	} else {
		return nil, fmt.Errorf("expected app implementation to return Root(s) but found: %s", returnValue.Type())
	}

	return ScreensFromRoots(roots), nil
}

// Calls any callable from Applet.Globals. Pass args and receive a
// starlark Value, or an error if you're unlucky.
func (a *Applet) Call(callable *starlark.Function, args starlark.Tuple) (val starlark.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while running %s: %v", a.Filename, r)
		}
	}()

	resultVal, err := starlark.Call(a.thread(), callable, args, nil)
	if err != nil {
		evalErr, ok := err.(*starlark.EvalError)
		if ok {
			return nil, errors.New(evalErr.Backtrace())
		}
		return nil, fmt.Errorf(
			"in %s at %s: %s",
			callable.Name(),
			callable.Position().String(),
			err,
		)
	}

	return resultVal, nil
}

func (a *Applet) loadModule(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if a.loader != nil {
		mod, err := a.loader(thread, module)
		if err == nil {
			return mod, nil
		}
	}

	switch module {
	case "render.star":
		return LoadModule()

	case "cache.star":
		return LoadCacheModule()

	case "encoding/base64.star":
		return starlibbase64.LoadModule()

	case "encoding/json.star":
		return starlibjson.LoadModule()

	case "http.star":
		return starlibhttp.LoadModule()

	case "math.star":
		return starlibmath.LoadModule()

	case "re.star":
		return starlibre.LoadModule()

	case "time.star":
		return starlibtime.LoadModule()

	default:
		return nil, fmt.Errorf("invalid module: %s", module)
	}
}
