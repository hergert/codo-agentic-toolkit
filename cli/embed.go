package main

import (
	"embed"
	"io/fs"
)

// embeddedPack embeds the pack directory that is temporarily copied during build
// The build process copies ../pack to cli/.embedded_pack
//
//go:embed all:.embedded_pack
var embeddedPack embed.FS

// GetEmbeddedPack returns the embedded pack filesystem for use by internal packages
func GetEmbeddedPack() (fs.FS, error) {
	return fs.Sub(embeddedPack, ".embedded_pack")
}