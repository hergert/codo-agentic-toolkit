package pack

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack/zipbuild"
)

// TestEmbeddedBaseZipDeterministic ensures the generated embedded bytes match
// a freshly built archive of pack/. This guards against drift between the
// embedded fallback and the distributable pack.
func TestEmbeddedBaseZipDeterministic(t *testing.T) {
	t.Helper()

	root := filepath.Clean(filepath.Join("..", "..", ".."))
	packDir := filepath.Join(root, "pack")
	packDir, err := filepath.Abs(packDir)
	if err != nil {
		t.Fatalf("abs pack dir: %v", err)
	}
	generated, err := zipbuild.Build(packDir)
	if err != nil {
		t.Fatalf("build zip: %v", err)
	}

	if len(embeddedBaseZip) == 0 {
		t.Fatalf("embeddedBaseZip not generated; run go generate ./cli/internal/pack")
	}

	if !bytes.Equal(embeddedBaseZip, generated) {
		t.Fatalf("embedded pack bytes differ from pack directory snapshot")
	}
}
