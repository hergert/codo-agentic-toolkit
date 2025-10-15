package main

import (
	"github.com/hergert/codo-agentic-toolkit/cli/cmd"
	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
)

func init() {
	// Initialize the embedded pack for internal packages
	if fs, err := GetEmbeddedPack(); err == nil {
		pack.SetEmbeddedFS(fs)
	}
}

func main() {
	cmd.Execute()
}
