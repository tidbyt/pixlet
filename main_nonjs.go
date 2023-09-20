//go:build !js && !wasm

package main

import (
	"tidbyt.dev/pixlet/cmd/private"
)

func init() {
	rootCmd.AddCommand(private.CreateCmd)
}
