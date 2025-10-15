package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const releaseRepo = "hergert/codo-agentic-toolkit"

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"self-update"},
	Short:   "Upgrade the codo binary to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := semver.Parse(strings.TrimPrefix(version, "v"))
		if err != nil {
			return fmt.Errorf("parse current version: %w", err)
		}

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("get executable path: %w", err)
		}

		latest, err := selfupdate.UpdateCommand(exe, v, releaseRepo)
		if err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		if latest == nil || !latest.Version.GT(v) {
			fmt.Printf("Already up to date (v%s)\n", v)
			return nil
		}

		fmt.Printf("Updated to v%s\n", latest.Version)
		return nil
	},
}
