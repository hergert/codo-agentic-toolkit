package statepath

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

func baseDir() (string, error) {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "codo"), nil
	}
	if dir, err := os.UserCacheDir(); err == nil && dir != "" {
		return filepath.Join(dir, "codo"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".codo"), nil
}

func repoKey(root string) string {
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}
	sum := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(sum[:])[:16]
}

func manifestPath(root string) (string, error) {
	base, err := baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "manifests", repoKey(root)+".json"), nil
}

func backupRoot(root string) (string, error) {
	base, err := baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "backups", repoKey(root)), nil
}

// ManifestPath returns the absolute path to the manifest file for the repository root.
func ManifestPath(root string) (string, error) {
	return manifestPath(root)
}

// BackupDir returns the directory path under which backups should be stored for the repo.
func BackupDir(root, timestamp string) (string, error) {
	base, err := backupRoot(root)
	if err != nil {
		return "", err
	}
	return filepath.Join(base, timestamp), nil
}

// LegacyManifestPath returns the old in-repo manifest location for migration/removal.
func LegacyManifestPath(root string) string {
	return filepath.Join(root, ".claude", ".codo-manifest.json")
}
