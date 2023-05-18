//go:build !js && !wasm

package main

import (
	"tidbyt.dev/pixlet/cmd"
)

func init() {
	rootCmd.AddCommand(cmd.CreateCmd)
}
