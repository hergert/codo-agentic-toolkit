package main

import (
	"archive/zip"
	"bytes"
	"io/fs"

	_ "embed"
)

// embeddedPackZip stores a zipped copy of the base pack for offline installs.
//
//go:embed internal/pack/embedded_base.zip
var embeddedPackZip []byte

// GetEmbeddedPack returns the embedded pack as an fs.FS backed by the zip archive.
func GetEmbeddedPack() (fs.FS, error) {
	if len(embeddedPackZip) == 0 {
		return nil, fs.ErrNotExist
	}
	reader, err := zip.NewReader(bytes.NewReader(embeddedPackZip), int64(len(embeddedPackZip)))
	if err != nil {
		return nil, err
	}
	return reader, nil
}
