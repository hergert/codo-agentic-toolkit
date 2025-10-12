package pack

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEmbeddedBaseZipDeterministic ensures the checked-in embedded_base.zip
// matches the output of scripts/embed-pack.sh. This guards against the
// release workflow publishing a different pack than the one tested locally.
func TestEmbeddedBaseZipDeterministic(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working dir: %v", err)
	}

	embedPath := filepath.Join(wd, "embedded_base.zip")
	original, err := os.ReadFile(embedPath)
	if err != nil {
		t.Fatalf("read embedded zip: %v", err)
	}

	restore := true
	defer func() {
		if restore {
			if writeErr := os.WriteFile(embedPath, original, 0o644); writeErr != nil {
				t.Fatalf("restore embedded zip: %v", writeErr)
			}
		}
	}()

	root := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	scriptPath := filepath.Join(root, "scripts", "embed-pack.sh")

	cmd := exec.Command(scriptPath)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("running embed-pack.sh failed: %v\n%s", err, string(out))
	}

	updated, err := os.ReadFile(embedPath)
	if err != nil {
		t.Fatalf("read regenerated zip: %v", err)
	}

	if !bytes.Equal(original, updated) {
		restore = false
		t.Fatalf("embedded_base.zip differs from scripts/embed-pack.sh output")
	}

	// No changes detected, keep the file as-is.
	restore = false
}
