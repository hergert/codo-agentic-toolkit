package pack

import (
	"archive/zip"
	"bytes"
	"io/fs"
)

var embeddedBaseZip []byte

// LoadEmbeddedPack returns the embedded pack filesystem if the generated
// byte slice is present. When building from a clean tree, go generate must
// refresh embeddedBaseZip before compilation.
func LoadEmbeddedPack() (fs.FS, error) {
	if len(embeddedBaseZip) == 0 {
		return nil, fs.ErrNotExist
	}
	reader, err := zip.NewReader(bytes.NewReader(embeddedBaseZip), int64(len(embeddedBaseZip)))
	if err != nil {
		return nil, err
	}
	return reader, nil
}
