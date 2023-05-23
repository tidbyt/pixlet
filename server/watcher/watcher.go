// Package watcher provides a simple mechanism to watch a file for changes.
package watcher

import (
	"path/filepath"
)

// Watcher is a structure to watch a file for changes and notify a channel.
type Watcher struct {
	filename    string
	fileChanges chan bool
}

// NewWatcher instantiates a new watcher with the provided filename and changes
// channel.
func NewWatcher(filename string, fileChanges chan bool) *Watcher {
	return &Watcher{
		filename:    filepath.FromSlash(filename),
		fileChanges: fileChanges,
	}
}
