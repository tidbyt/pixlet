package runtime

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
	"tidbyt.dev/pixlet/runtime/modules/testing"
)

type TestError string

func (te TestError) String() string {
	return string(te)
}

type TestResult struct {
	FunctionName string
	Errors       []TestError
}

func (tr *TestResult) Success() bool {
	return len(tr.Errors) == 0
}

type TestResultReporter struct {
	Errors []TestError
}

func (trr *TestResultReporter) Error(args ...interface{}) {
	trr.Errors = append(trr.Errors, TestError(fmt.Sprint(args...)))
}

func (a *Applet) LoadTests(filename string, src []byte, config map[string]string, initializers ...ThreadInitializer) (err error) {
	testGlobals, err := starlark.ExecFile(a.thread(), filename, src, a.Globals)
	if err != nil {
		return fmt.Errorf("starlark.ExecFile: %v", err)
	}

	for k, v := range testGlobals {
		a.Globals[k] = v
	}

	return nil
}

func (a *Applet) RunTests(config map[string]string, initializers ...ThreadInitializer) ([]TestResult, error) {
	testResults := []TestResult{}
	for key, fun := range a.Globals {
		if strings.HasPrefix(key, "test_") {
			testFun, ok := fun.(*starlark.Function)
			if ok {
				testErrors := a.runTest(testFun, config, initializers...)
				result := TestResult{
					FunctionName: key,
					Errors:       testErrors,
				}
				testResults = append(testResults, result)
			}
		}
	}

	return testResults, nil
}

func (a *Applet) runTest(testFunction *starlark.Function, config map[string]string, initializers ...ThreadInitializer) []TestError {
	var args starlark.Tuple

	if testFunction.NumParams() > 0 {
		starlarkConfig := AppletConfig(config)
		args = starlark.Tuple{starlarkConfig}
	}

	testReporter := &TestResultReporter{}
	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		starlarktest.SetReporter(thread, testReporter)
		return thread
	})

	returnValue, err := a.Call(testFunction, args, initializers...)
	if err != nil {
		testReporter.Error(err)
		return testReporter.Errors
	}

	if returnValue != starlark.None {
		testReporter.Error(fmt.Errorf("Error: expected test implementation to return None but found: %s", returnValue.Type()))
		return testReporter.Errors
	}

	return testReporter.Errors
}

func TestingLoader(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "testing.star" {
		return testing.LoadModule()
	}
	return nil, fmt.Errorf("module not found in custom loader %s", module)
}
