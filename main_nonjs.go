//go:build !js && !wasm

package main

import (
	"tidbyt.dev/pixlet/cmd"
	"tidbyt.dev/pixlet/cmd/private"
)

func init() {
	rootCmd.AddCommand(private.PrivateCmd)
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.ServeCmd)
}
