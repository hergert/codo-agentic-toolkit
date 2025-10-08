package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/fsops"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update the codo binary to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, found, err := selfupdate.DetectLatest("hergert/codo-agentic-toolkit")
		if err != nil {
			return err
		}
		v, err := semver.Parse(version)
		if err != nil {
			return err
		}
		if !found || !latest.Version.GT(v) {
			fmt.Println("codo is up to date:", version)
			return nil
		}
		sumURL := ""
		for _, asset := range latest.ReleaseAssets {
			if strings.HasSuffix(asset.GetName(), ".sha256") {
				sumURL = asset.GetBrowserDownloadURL()
				break
			}
		}
		if sumURL == "" {
			return fmt.Errorf("checksum asset not found for %s", latest.Version)
		}

		checksum, err := downloadString(sumURL)
		if err != nil {
			return fmt.Errorf("download checksum: %w", err)
		}

		tmpPath, err := downloadAsset(latest.AssetURL)
		if err != nil {
			return err
		}
		defer os.Remove(tmpPath)

		digest, err := fileSha256(tmpPath)
		if err != nil {
			return err
		}
		if !strings.Contains(checksum, digest) {
			return fmt.Errorf("checksum mismatch: expected entry in %s, got %s", sumURL, digest)
		}

		exe, err := os.Executable()
		if err != nil {
			exe = os.Args[0]
		}
		if err := selfupdate.ReplaceExecutable(tmpPath, exe); err != nil {
			return err
		}
		fmt.Println("Updated to", latest.Version)
		return nil
	},
}

func fileSha256(path string) (string, error) {
	return fsops.Sha256File(path)
}

func downloadString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func downloadAsset(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	tmp, err := os.CreateTemp("", "codo-update-*")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	defer tmp.Close()
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(path)
		return "", err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(path)
		return "", err
	}
	return path, nil
}
