package pack

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Resolve downloads and verifies a pack from GitHub releases
func Resolve(tag string) (string, error) {
	baseURL := "https://github.com/hergert/codo-agentic-toolkit/releases"
	if tag == "latest" {
		tag = "latest/download"
	} else {
		tag = "download/" + tag
	}

	packURL := fmt.Sprintf("%s/%s/dotclaude-pack.zip", baseURL, tag)
	checksumURL := fmt.Sprintf("%s/%s/dotclaude-pack.sha256", baseURL, tag)

	// Download to ~/.codo/packs/<tag>/
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(home, ".codo", "packs", strings.ReplaceAll(tag, "/", "_"))
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	zipPath := filepath.Join(cacheDir, "dotclaude-pack.zip")
	
	// Download pack
	if err := downloadFile(zipPath, packURL); err != nil {
		return "", fmt.Errorf("failed to download pack: %w", err)
	}

	// Download and verify checksum
	checksumPath := filepath.Join(cacheDir, "dotclaude-pack.sha256")
	if err := downloadFile(checksumPath, checksumURL); err != nil {
		return "", fmt.Errorf("failed to download checksum: %w", err)
	}

	// Read expected checksum
	expectedBytes, err := os.ReadFile(checksumPath)
	if err != nil {
		return "", err
	}
	expected := strings.TrimSpace(string(expectedBytes))

	// Calculate actual checksum
	actual, err := fileChecksum(zipPath)
	if err != nil {
		return "", err
	}

	if actual != expected {
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}

	// Extract pack
	extractDir := filepath.Join(cacheDir, "dotclaude")
	if err := extractZip(zipPath, extractDir); err != nil {
		return "", fmt.Errorf("failed to extract pack: %w", err)
	}

	return extractDir, nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func fileChecksum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			rc.Close()
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			rc.Close()
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}