package starlarkutil

import (
	"context"

	"go.starlark.net/starlark"
)

const (
	// ThreadContextKey is the name of the Starlark thread-local that we use to
	// pass context around.
	ThreadContextKey = "tidbyt.dev/pixlet/starlarkutil/$ctx"
)

// AttachThreadContext attaches context to a Starlark thread so that it can be
// retrieved latter with `ThreadContext`.
func AttachThreadContext(ctx context.Context, thread *starlark.Thread) {
	thread.SetLocal(ThreadContextKey, ctx)
}

// ThreadContext returns the context that was attached to a Starlark thread
// by `AttachThreadContext`. If no context is attached to the thread, it
// returns a new, empty context.
func ThreadContext(thread *starlark.Thread) context.Context {
	ctx, ok := thread.Local(ThreadContextKey).(context.Context)
	if !ok {
		ctx = context.Background()
	}
	return ctx
}
