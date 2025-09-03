package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    "github.com/youruser/codo/internal/fsops"
    "github.com/youruser/codo/internal/manifest"
    "github.com/youruser/codo/internal/pack"
    "github.com/youruser/codo/internal/tui"
)

var initVersion string
var initDryRun bool

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Install the toolkit into this repository (safe-by-default)",
    RunE: func(cmd *cobra.Command, args []string) error {
        root, _ := os.Getwd()
        ctx := context.Background()

        choices, err := tui.RunInitWizard(ctx)
        if err != nil { return err }
        if !choices.Confirmed {
            fmt.Println("aborted")
            return nil
        }

        files, err := pack.Files(initVersion, choices.Stacks)
        if err != nil { return err }

        reportDir := filepath.Join(".claude", ".codo-report")
        _ = os.MkdirAll(reportDir, 0o755)

        // Copy safely (or simulate with --dry-run). fsops prints +/=!/conflict lines.
        for _, f := range files {
            if err := fsops.CopySafe(f, root, initDryRun); err != nil { return err }
        }
        if !initDryRun {
            if err := fsops.ChmodHooks(); err != nil { return err }
            if err := manifest.Write(files, pack.VersionOrDefault(initVersion)); err != nil { return err }
        }
        fmt.Printf("\nCodo %s initialized. See .claude/.codo-report/conflicts.txt if any.\n", pack.VersionOrDefault(initVersion))
        return nil
    },
}

func init() {
    initCmd.Flags().StringVar(&initVersion, "version", "", "Pack version (default embedded)")
    initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview only; do not write files")
}

