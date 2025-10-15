package pack

import (
	"io/fs"
)

// embeddedFS holds the embedded filesystem set from main package
var embeddedFS fs.FS

// SetEmbeddedFS sets the embedded filesystem (called from main)
func SetEmbeddedFS(fs fs.FS) {
	embeddedFS = fs
}

// GetEmbeddedBaseFS returns the embedded dotclaude pack filesystem
func GetEmbeddedBaseFS() (fs.FS, error) {
	if embeddedFS == nil {
		return nil, fs.ErrNotExist
	}
	return embeddedFS, nil
}
