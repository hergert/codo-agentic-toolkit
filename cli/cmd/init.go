package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"

    "github.com/youruser/codo/internal/fsops"
    "github.com/youruser/codo/internal/manifest"
    "github.com/youruser/codo/internal/pack"
    "github.com/youruser/codo/internal/tui"
)

var initVersion string
var initDryRun bool
var initStacks string // comma-separated
var initNoTUI bool    // headless mode

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Install the toolkit into this repository (safe-by-default)",
    RunE: func(cmd *cobra.Command, args []string) error {
        root, _ := os.Getwd()
        ctx := context.Background()

        var choices tui.InitResult
        if initNoTUI || initStacks != "" {
            // Headless: parse stacks from flag
            keys := []string{}
            for _, s := range strings.Split(initStacks, ",") {
                ss := strings.TrimSpace(s)
                if ss != "" { keys = append(keys, ss) }
            }
            choices = tui.InitResult{Stacks: keys, Confirmed: true}
        } else {
            var err error
            choices, err = tui.RunInitWizard(ctx)
            if err != nil { return err }
            if !choices.Confirmed {
                fmt.Println("aborted")
                return nil
            }
        }

        // Prefer local pack at dotclaude/ for development; fallback to embedded pack.
        if _, err := os.Stat("dotclaude"); err != nil {
            return fmt.Errorf("dotclaude pack not found; expected ./dotclaude directory")
        }
        rootFS := os.DirFS("dotclaude")
        files, err := pack.FilesFromDotclaudeFS(rootFS, choices.Stacks)
        if err != nil { return err }

        reportDir := filepath.Join(".claude", ".codo-report")
        _ = os.MkdirAll(reportDir, 0o755)

        // Copy safely (or simulate with --dry-run). fsops prints +/=!/conflict lines.
        for _, f := range files {
            if err := fsops.CopySafe(f, root, initDryRun); err != nil { return err }
        }
        if !initDryRun {
            if err := fsops.ChmodHooks(); err != nil { return err }
            installedVersion := initVersion
            if installedVersion == "" { installedVersion = "local" }
            if err := manifest.WriteWithStacks(files, installedVersion, choices.Stacks); err != nil { return err }
        }
        installedVersion := initVersion
        if installedVersion == "" { installedVersion = "local" }
        fmt.Printf("\nCodo %s initialized. See .claude/.codo-report/conflicts.txt if any.\n", installedVersion)
        return nil
    },
}

func init() {
    initCmd.Flags().StringVar(&initVersion, "version", "", "Pack tag (e.g. v1.2.0) or 'local'")
    initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview only; do not write files")
    initCmd.Flags().StringVar(&initStacks, "stacks", "", "Comma-separated stacks (skip TUI)")
    initCmd.Flags().BoolVar(&initNoTUI, "no-tui", false, "Don't show the TUI wizard")
}
