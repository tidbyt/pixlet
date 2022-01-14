package starlarkutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

type contextKey string

func TestThreadContext(t *testing.T) {
	var key contextKey
	key = "foo"
	val := "bar"

	ctx := context.WithValue(
		context.Background(),
		key, val,
	)

	thread := &starlark.Thread{}
	AttachThreadContext(ctx, thread)

	ctxFromThread := ThreadContext(thread)
	assert.Same(t, ctx, ctxFromThread)
	assert.Equal(t, val, ctx.Value(key))
}

func TestThreadWithoutContext(t *testing.T) {
	thread := &starlark.Thread{}
	assert.NotNil(t, ThreadContext(thread))
}
