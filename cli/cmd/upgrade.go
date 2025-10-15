package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver"
	update "github.com/inconshreveable/go-update"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const releaseRepo = "hergert/codo-agentic-toolkit"

var upgradeHTTPClient = &http.Client{Timeout: 30 * time.Second}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"self-update"},
	Short:   "Upgrade the codo binary to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		if version == "dev" {
			fmt.Println("Development build; skip self-upgrade.")
			return nil
		}

		current, err := semver.Parse(strings.TrimPrefix(version, "v"))
		if err != nil {
			return fmt.Errorf("parse current version: %w", err)
		}

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("get executable path: %w", err)
		}

		release, found, err := selfupdate.DetectLatest(releaseRepo)
		if err != nil {
			return fmt.Errorf("detect latest release: %w", err)
		}
		if !found {
			return fmt.Errorf("no releases found")
		}
		if !release.Version.GT(current) {
			fmt.Printf("Already up to date (v%s)\n", current)
			return nil
		}

		tag := "v" + release.Version.String()
		assetName, err := assetForCurrentPlatform(tag)
		if err != nil {
			return err
		}

		expected, err := fetchExpectedChecksum(tag, assetName)
		if err != nil {
			return err
		}

		archiveBytes, err := downloadReleaseAsset(tag, assetName)
		if err != nil {
			return err
		}

		actualSum := sha256.Sum256(archiveBytes)
		if hex.EncodeToString(actualSum[:]) != expected {
			return fmt.Errorf("checksum mismatch for %s", assetName)
		}

		binary, err := extractBinaryFromArchive(assetName, archiveBytes)
		if err != nil {
			return err
		}

		if err := applyUpdateBinary(exe, binary); err != nil {
			return err
		}

		fmt.Printf("Updated to %s\n", tag)
		return nil
	},
}

func assetForCurrentPlatform(tag string) (string, error) {
	osName := runtime.GOOS
	archName := runtime.GOARCH
	ext := ".tar.gz"
	if osName == "windows" {
		ext = ".zip"
	}

	supported := map[string]struct{}{
		"linux":   {},
		"darwin":  {},
		"windows": {},
	}
	if _, ok := supported[osName]; !ok {
		return "", fmt.Errorf("unsupported OS %q", osName)
	}

	if archName != "amd64" && archName != "arm64" {
		return "", fmt.Errorf("unsupported architecture %q", archName)
	}

	return fmt.Sprintf("codo_%s_%s_%s%s", tag, osName, archName, ext), nil
}

func fetchExpectedChecksum(tag, asset string) (string, error) {
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/checksums.txt", releaseRepo, tag)
	resp, err := upgradeHTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("download checksum: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download checksum: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[len(fields)-1]
		if name == asset || filepath.Base(name) == asset {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("checksum for %s not found", asset)
}

func downloadReleaseAsset(tag, asset string) ([]byte, error) {
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", releaseRepo, tag, asset)
	resp, err := upgradeHTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download asset: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download asset: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func extractBinaryFromArchive(asset string, archive []byte) ([]byte, error) {
	target := "codo"
	if runtime.GOOS == "windows" {
		target = "codo.exe"
	}

	if strings.HasSuffix(asset, ".zip") {
		readerAt := bytes.NewReader(archive)
		zr, err := zip.NewReader(readerAt, int64(len(archive)))
		if err != nil {
			return nil, fmt.Errorf("open zip: %w", err)
		}
		for _, f := range zr.File {
			if filepath.Base(f.Name) != target {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open zip entry: %w", err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
		return nil, fmt.Errorf("binary %s not found in archive", target)
	}

	gr, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("open gzip: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar: %w", err)
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		if filepath.Base(hdr.Name) != target {
			continue
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, tr); err != nil {
			return nil, fmt.Errorf("copy binary: %w", err)
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("binary %s not found in archive", target)
}

func applyUpdateBinary(exe string, binary []byte) error {
	reader := bytes.NewReader(binary)
	if err := update.Apply(reader, update.Options{TargetPath: exe}); err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			return fmt.Errorf("apply update: %w (rollback failed: %v)", err, rerr)
		}
		return fmt.Errorf("apply update: %w", err)
	}
	if err := os.Chmod(exe, 0o755); err != nil {
		return fmt.Errorf("chmod updated binary: %w", err)
	}
	return nil
}
