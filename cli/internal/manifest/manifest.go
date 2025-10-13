package manifest

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
	"github.com/hergert/codo-agentic-toolkit/cli/internal/statepath"
)

type Entry struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	Unmanaged bool   `json:"unmanaged,omitempty"`
}
type Manifest struct {
	Version     string   `json:"version"`
	InstalledAt string   `json:"installed_at"`
	Files       []Entry  `json:"files"`
	Stacks      []string `json:"stacks,omitempty"`
}

func repoRoot() (string, error) {
	return os.Getwd()
}

func manifestPath() (string, error) {
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	return statepath.ManifestPath(root)
}

func legacyManifestPath() (string, error) {
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	return statepath.LegacyManifestPath(root), nil
}

func Exists() bool {
	path, err := manifestPath()
	if err == nil {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	legacy, err := legacyManifestPath()
	if err == nil {
		if _, err := os.Stat(legacy); err == nil {
			return true
		}
	}
	return false
}

func Write(files []pack.File, version string) error {
	return WriteWithStacks(files, version, nil, nil)
}

func WriteWithStacks(files []pack.File, version string, stacks []string, unmanaged map[string]bool) error {
	path, err := manifestPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var entries []Entry
	for _, f := range files {
		// use on-disk hash if placed; else hash of new content
		dst := f.RelPath
		var sum string
		if b, err := os.ReadFile(dst); err == nil {
			sum = fmt.Sprintf("%x", sha256.Sum256(b))
		} else {
			b, err := f.Read()
			if err != nil {
				return err
			}
			sum = fmt.Sprintf("%x", sha256.Sum256(b))
		}
		entries = append(entries, Entry{Path: dst, SHA256: sum, Unmanaged: unmanaged != nil && unmanaged[dst]})
	}
	m := Manifest{Version: version, InstalledAt: "", Files: entries, Stacks: stacks}
	buf, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return err
	}
	if legacy, err := legacyManifestPath(); err == nil {
		_ = os.Remove(legacy)
	}
	return nil
}

func Open() (Manifest, error) {
	var m Manifest
	path, err := manifestPath()
	if err == nil {
		if f, err := os.Open(path); err == nil {
			defer f.Close()
			b, _ := io.ReadAll(f)
			err = json.Unmarshal(b, &m)
			if err == nil {
				return m, nil
			}
		}
	}
	legacy, err := legacyManifestPath()
	if err != nil {
		return m, err
	}
	f, err := os.Open(legacy)
	if err != nil {
		return m, err
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	err = json.Unmarshal(b, &m)
	return m, err
}

func Remove() {
	if path, err := manifestPath(); err == nil {
		_ = os.Remove(path)
	}
	if legacy, err := legacyManifestPath(); err == nil {
		_ = os.Remove(legacy)
	}
}
