package schema_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var notificationSource = `
load("assert.star", "assert")
load("schema.star", "schema")
load("sound.mp3", "file")

sounds = [
	schema.Sound(
		title = "Ding!",
		file = file,
	),

]

s = schema.Notification(
	id = "notification1",
	name = "New message",
	desc = "A new message has arrived",
	icon = "message",
	sounds = sounds,
)

assert.eq(s.id, "notification1")
assert.eq(s.name, "New message")
assert.eq(s.desc, "A new message has arrived")
assert.eq(s.icon, "message")

assert.eq(s.sounds[0].title, "Ding!")
assert.eq(s.sounds[0].file, file)

def main():
	return []
`

func TestNotification(t *testing.T) {
	vfs := fstest.MapFS{
		"sound.mp3":         &fstest.MapFile{Data: []byte("sound data")},
		"notification.star": &fstest.MapFile{Data: []byte(notificationSource)},
	}
	app, err := runtime.NewAppletFromFS("sound", vfs)
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
