package main

import (
	"github.com/hergert/codo-agentic-toolkit/cli/cmd"
	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
)

func init() {
	if fs, err := pack.LoadEmbeddedPack(); err == nil {
		pack.SetEmbeddedFS(fs)
	}
}

func main() {
	cmd.Execute()
}
