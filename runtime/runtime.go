package runtime

import (
	"crypto/md5"
	"fmt"

	"github.com/pkg/errors"
	starlibbase64 "github.com/qri-io/starlib/encoding/base64"
	starlibhttp "github.com/qri-io/starlib/http"
	starlibre "github.com/qri-io/starlib/re"
	starlibjson "go.starlark.net/lib/json"
	starlibmath "go.starlark.net/lib/math"
	starlibtime "go.starlark.net/lib/time"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"tidbyt.dev/pixlet/render"
)

type ModuleLoader func(*starlark.Thread, string) (starlark.StringDict, error)

// ThreadInitializer is called when building a Starlark thread to run an applet
// on. It can customize the thread by overriding behavior or attaching thread
// local data.
type ThreadInitializer func(thread *starlark.Thread) *starlark.Thread

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

func (a *Applet) thread(initializers ...ThreadInitializer) *starlark.Thread {
	t := &starlark.Thread{
		Name: a.Id,
		Load: a.loadModule,
		Print: func(thread *starlark.Thread, msg string) {
			fmt.Printf("[%s] %s\n", a.Filename, msg)
		},
	}

	for _, init := range initializers {
		t = init(t)
	}

	return t
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
func (a *Applet) Run(config map[string]string, initializers ...ThreadInitializer) (roots []render.Root, err error) {
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

	returnValue, err := a.Call(a.main, args, initializers...)
	if err != nil {
		return nil, err
	}

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

	return roots, nil
}

// Calls any callable from Applet.Globals. Pass args and receive a
// starlark Value, or an error if you're unlucky.
func (a *Applet) Call(callable *starlark.Function, args starlark.Tuple, initializers ...ThreadInitializer) (val starlark.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while running %s: %v", a.Filename, r)
		}
	}()

	resultVal, err := starlark.Call(a.thread(initializers...), callable, args, nil)
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

	case "xpath.star":
		return LoadXPathModule()

	case "encoding/base64.star":
		return starlibbase64.LoadModule()

	case "encoding/json.star":
		return starlark.StringDict{
			starlibjson.Module.Name: starlibjson.Module,
		}, nil

	case "http.star":
		return starlibhttp.LoadModule()

	case "math.star":
		return starlark.StringDict{
			starlibmath.Module.Name: starlibmath.Module,
		}, nil

	case "re.star":
		return starlibre.LoadModule()

	case "time.star":
		return starlark.StringDict{
			starlibtime.Module.Name: starlibtime.Module,
		}, nil

	default:
		return nil, fmt.Errorf("invalid module: %s", module)
	}
}
