package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment requirements",
	RunE: func(cmd *cobra.Command, args []string) error {
		checks := []struct{ name, bin string }{
			{"git", "git"},
		}
		for _, c := range checks {
			if _, err := exec.LookPath(c.bin); err != nil {
				fmt.Println("✗", c.name)
			} else {
				fmt.Println("✓", c.name)
			}
		}

		// Check Python version
		if err := checkPython(); err != nil {
			fmt.Println("✗ Python 3:", err)
		}

		for _, h := range []string{
			".claude/hooks/pre_tool_use.py",
			".claude/hooks/post_tool_use.py",
			".claude/hooks/user_prompt_submit.py",
		} {
			if fi, err := os.Stat(h); err == nil && fi.Mode()&0o111 != 0 {
				fmt.Println("✓ exec", h)
			} else {
				fmt.Println("✗ exec", h)
			}
		}

		if err := emitUpgradeStatus(); err != nil {
			fmt.Println("✗ upgrade check:", err)
		}
		return nil
	},
}

func checkPython() error {
	candidates := [][]string{
		{"python3", "--version"},
		{"py", "-3", "--version"},
		{"python", "--version"},
	}
	for _, c := range candidates {
		if _, err := exec.LookPath(c[0]); err != nil {
			continue
		}
		out, err := exec.Command(c[0], c[1:]...).Output()
		if err != nil {
			continue
		}
		text := strings.ToLower(string(out))
		if strings.Contains(text, "python 3") {
			fmt.Printf("✓ %s %s\n", c[0], strings.TrimSpace(string(out)))
			return nil
		}
	}
	return fmt.Errorf("Python 3 not found (required for hooks)")
}

func emitUpgradeStatus() error {
	if version == "dev" {
		fmt.Println("✓ upgrade check: development build (skipping release detection)")
		return nil
	}

	v, err := semver.Parse(strings.TrimPrefix(version, "v"))
	if err != nil {
		return fmt.Errorf("parse current version %q: %w", version, err)
	}

	latest, found, err := selfupdate.DetectLatest(releaseRepo)
	if err != nil {
		return err
	}
	if !found {
		fmt.Println("✓ upgrade check: no published releases found")
		return nil
	}

	if !latest.Version.GT(v) {
		fmt.Printf("✓ upgrade check: running latest (v%s)\n", v)
		return nil
	}

	fmt.Printf("✗ upgrade available: v%s (current v%s)\n", latest.Version, v)
	fmt.Println("  run `codo upgrade` to install")
	return nil
}
