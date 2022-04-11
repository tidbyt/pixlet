package runtime

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	starlibbase64 "github.com/qri-io/starlib/encoding/base64"
	starlibcsv "github.com/qri-io/starlib/encoding/csv"
	starlibhash "github.com/qri-io/starlib/hash"
	starlibhtml "github.com/qri-io/starlib/html"
	starlibhttp "github.com/qri-io/starlib/http"
	starlibre "github.com/qri-io/starlib/re"
	starlibjson "go.starlark.net/lib/json"
	starlibmath "go.starlark.net/lib/math"
	starlibtime "go.starlark.net/lib/time"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/starlarktest"

	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/runtime/modules/animation_runtime"
	"tidbyt.dev/pixlet/runtime/modules/humanize"
	"tidbyt.dev/pixlet/runtime/modules/random"
	"tidbyt.dev/pixlet/runtime/modules/render_runtime"
	"tidbyt.dev/pixlet/runtime/modules/sunrise"
	"tidbyt.dev/pixlet/schema"
	"tidbyt.dev/pixlet/starlarkutil"
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
	Filename            string
	Id                  string
	Globals             starlark.StringDict
	SecretDecryptionKey *SecretDecryptionKey

	src         []byte
	loader      ModuleLoader
	predeclared starlark.StringDict
	main        *starlark.Function

	schema     *schema.Schema
	schemaJSON []byte
	decrypter  decrypter
}

func (a *Applet) thread(initializers ...ThreadInitializer) *starlark.Thread {
	t := &starlark.Thread{
		Name: a.Id,
		Load: a.loadModule,
		Print: func(thread *starlark.Thread, msg string) {
			fmt.Printf("[%s] %s\n", a.Filename, msg)
		},
	}

	if a.decrypter != nil {
		a.decrypter.attachToThread(t)
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

	if a.SecretDecryptionKey != nil {
		a.decrypter, err = a.SecretDecryptionKey.decrypterForApp(a)
		if err != nil {
			return errors.Wrapf(err, "preparing secret key for %s", a.Filename)
		}
	}

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

	schemaFun, _ := a.Globals[schema.SchemaFunctionName].(*starlark.Function)
	if schemaFun != nil {
		schemaVal, err := a.Call(schemaFun, nil)
		if err != nil {
			return errors.Wrapf(err, "calling schema function for %s", a.Filename)
		}

		a.schema, err = schema.FromStarlark(schemaVal, a.Globals)
		if err != nil {
			return errors.Wrapf(err, "parsing schema for %s", a.Filename)
		}

		a.schemaJSON, err = json.Marshal(a.schema)
		if err != nil {
			return errors.Wrapf(err, "serializing schema to JSON for %s", a.Filename)
		}
	}

	return nil
}

// Runs the applet's main function, passing it configuration as a
// starlark dict.
func (a *Applet) Run(config map[string]string, initializers ...ThreadInitializer) (roots []render.Root, err error) {
	var args starlark.Tuple
	if a.main.NumParams() > 0 {
		starlarkConfig := AppletConfig(config)
		args = starlark.Tuple{starlarkConfig}
	}

	returnValue, err := a.Call(a.main, args, initializers...)
	if err != nil {
		return nil, err
	}

	if returnRoot, ok := returnValue.(render_runtime.Rootable); ok {
		roots = []render.Root{returnRoot.AsRenderRoot()}
	} else if returnList, ok := returnValue.(*starlark.List); ok {
		roots = make([]render.Root, returnList.Len())
		iter := returnList.Iterate()
		defer iter.Done()
		i := 0
		var listVal starlark.Value
		for iter.Next(&listVal) {
			if listValRoot, ok := listVal.(render_runtime.Rootable); ok {
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

// CallSchemaHandler calls a schema handler, passing it a single
// string parameter and returning a single string value.
func (app *Applet) CallSchemaHandler(ctx context.Context, handlerName, parameter string) (result string, err error) {
	handler, found := app.schema.Handlers[handlerName]
	if !found {
		return "", fmt.Errorf("no exported handler named '%s'", handlerName)
	}

	resultVal, err := app.Call(
		handler.Function,
		starlark.Tuple{starlark.String(parameter)},
		attachContext(ctx),
	)
	if err != nil {
		return "", fmt.Errorf("calling schema handler %s: %v", handlerName, err)
	}

	switch handler.ReturnType {
	case schema.ReturnOptions:
		options, err := schema.EncodeOptions(resultVal)
		if err != nil {
			return "", err
		}
		return options, nil

	case schema.ReturnSchema:
		sch, err := schema.FromStarlark(resultVal, app.Globals)
		if err != nil {
			return "", err
		}

		s, err := json.Marshal(sch)
		if err != nil {
			return "", errors.Wrap(err, "serializing schema to JSON")
		}

		return string(s), nil

	case schema.ReturnString:
		str, ok := starlark.AsString(resultVal)
		if !ok {
			return "", fmt.Errorf(
				"expected %s to return a string or string-like value",
				handler.Function.Name(),
			)
		}
		return str, nil
	}

	return "", fmt.Errorf("a very unexpected error happened for handler \"%s\"", handlerName)
}

// GetSchema returns the config for the applet.
func (app *Applet) GetSchema() string {
	return string(app.schemaJSON)
}

func attachContext(ctx context.Context) ThreadInitializer {
	return func(thread *starlark.Thread) *starlark.Thread {
		starlarkutil.AttachThreadContext(ctx, thread)
		return thread
	}
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
		return render_runtime.LoadRenderModule()

	case "animation.star":
		return animation_runtime.LoadAnimationModule()

	case "schema.star":
		return schema.LoadModule()

	case "cache.star":
		return LoadCacheModule()

	case "secret.star":
		return LoadSecretModule()

	case "xpath.star":
		return LoadXPathModule()

	case "encoding/base64.star":
		return starlibbase64.LoadModule()

	case "encoding/csv.star":
		return starlibcsv.LoadModule()

	case "encoding/json.star":
		return starlark.StringDict{
			starlibjson.Module.Name: starlibjson.Module,
		}, nil

	case "hash.star":
		return starlibhash.LoadModule()

	case "http.star":
		return starlibhttp.LoadModule()

	case "html.star":
		return starlibhtml.LoadModule()

	case "humanize.star":
		return humanize.LoadModule()

	case "math.star":
		return starlark.StringDict{
			starlibmath.Module.Name: starlibmath.Module,
		}, nil

	case "re.star":
		return starlibre.LoadModule()

	case "sunrise.star":
		return sunrise.LoadModule()

	case "time.star":
		return starlark.StringDict{
			starlibtime.Module.Name: starlibtime.Module,
		}, nil

	case "random.star":
		return random.LoadModule()

	case "assert.star":
		return starlarktest.LoadAssertModule()

	default:
		return nil, fmt.Errorf("invalid module: %s", module)
	}
}
