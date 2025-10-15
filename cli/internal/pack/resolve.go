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
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// headOK checks if a URL is accessible via HEAD request
func headOK(u string) bool {
	req, err := http.NewRequest(http.MethodHead, u, nil)
	if err != nil {
		return false
	}
	r, err := httpClient.Do(req)
	if err != nil {
		gr, gerr := http.NewRequest(http.MethodGet, u, nil)
		if gerr != nil {
			return false
		}
		gr.Header.Set("Range", "bytes=0-0")
		resp, gerr := httpClient.Do(gr)
		if gerr != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent
	}
	r.Body.Close()
	return r.StatusCode == http.StatusOK
}

// Resolve downloads and verifies a pack from GitHub releases
func Resolve(tag string) (string, error) {
	baseURL := "https://github.com/hergert/codo-agentic-toolkit/releases"

	// Normalize tag for URL construction
	urlTag := tag
	if tag == "" || tag == "latest" {
		urlTag = "latest/download"
		tag = "latest"
	} else {
		urlTag = "download/" + tag
	}

	// Try both naming patterns: flat names (recommended) and tag-suffixed names (legacy)
	var packURL, checksumURL string
	patterns := []struct {
		zip string
		sha string
	}{
		{"dotclaude-pack.zip", "dotclaude-pack.sha256"},
		{fmt.Sprintf("dotclaude-pack-%s.zip", tag), fmt.Sprintf("dotclaude-pack-%s.sha256", tag)},
	}

	for _, p := range patterns {
		testPackURL := fmt.Sprintf("%s/%s/%s", baseURL, urlTag, p.zip)
		testChecksumURL := fmt.Sprintf("%s/%s/%s", baseURL, urlTag, p.sha)
		if headOK(testPackURL) && headOK(testChecksumURL) {
			packURL = testPackURL
			checksumURL = testChecksumURL
			break
		}
	}

	if packURL == "" {
		return "", fmt.Errorf("no pack asset found for tag=%s", tag)
	}

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
	expectedLine := strings.TrimSpace(string(expectedBytes))
	fields := strings.Fields(expectedLine)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty checksum file: %s", checksumPath)
	}
	expected := fields[0]

	// Calculate actual checksum
	actual, err := fileChecksum(zipPath)
	if err != nil {
		return "", err
	}

	if actual != expected {
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}

	// Extract pack
	extractDir := filepath.Join(cacheDir, "pack")
	if err := extractZip(zipPath, extractDir); err != nil {
		return "", fmt.Errorf("failed to extract pack: %w", err)
	}

	normalized, err := canonicalPackRoot(extractDir)
	if err != nil {
		return "", err
	}

	return normalized, nil
}

func canonicalPackRoot(root string) (string, error) {
	candidates := []string{root}
	files, err := os.ReadDir(root)
	if err == nil {
		for _, entry := range files {
			if entry.IsDir() && strings.EqualFold(entry.Name(), "pack") {
				candidates = append([]string{filepath.Join(root, entry.Name())}, candidates...)
				break
			}
		}
	}

	for _, dir := range candidates {
		if _, err := os.Stat(filepath.Join(dir, "dotclaude")); err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("dotclaude directory not found in pack (checked %v)", candidates)
}

func downloadFile(filepath string, url string) error {
	resp, err := httpClient.Get(url)
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
