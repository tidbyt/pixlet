package starlarkutil

import "go.starlark.net/starlark"

const (
	// ThreadOnExitKey is the key used to store functions that should be called
	// when a thread exits.
	ThreadOnExitKey = "tidbyt.dev/pixlet/runtime/on_exit"
)

type threadOnExitFunc func()

func AddOnExit(thread *starlark.Thread, fn threadOnExitFunc) {
	if onExit, ok := thread.Local(ThreadOnExitKey).(*[]threadOnExitFunc); ok {
		*onExit = append(*onExit, fn)
	} else {
		thread.SetLocal(ThreadOnExitKey, &[]threadOnExitFunc{fn})
	}
}

func RunOnExitFuncs(thread *starlark.Thread) {
	if onExit, ok := thread.Local(ThreadOnExitKey).(*[]threadOnExitFunc); ok {
		for _, fn := range *onExit {
			fn()
		}
	}
}
