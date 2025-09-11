package cmd

import (
    "context"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"

    "github.com/hergert/codo-agentic-toolkit/cli/internal/fsops"
    "github.com/hergert/codo-agentic-toolkit/cli/internal/manifest"
    "github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
    "github.com/hergert/codo-agentic-toolkit/cli/internal/tui"
)

var initVersion string
var initDryRun bool
var initStacks string // comma-separated
var initNoTUI bool    // headless mode
var initOffline bool  // force embedded base pack only

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

        // Pack resolution order:
        // 1. Local ./pack directory (for development)
        // 2. Downloaded pack (unless --offline)
        // 3. Embedded base pack (fallback)
        var rootFS fs.FS
        var packSource string
        
        if _, err := os.Stat("pack"); err == nil {
            // Use local pack for development
            rootFS = os.DirFS("pack")
            packSource = "local"
            fmt.Println("Using local pack directory")
        } else if !initOffline {
            // Try to download pack from GitHub
            versionToFetch := initVersion
            if versionToFetch == "" {
                versionToFetch = "latest"
            }
            fmt.Printf("Downloading pack version: %s...\n", versionToFetch)
            
            packPath, err := pack.Resolve(versionToFetch)
            if err == nil {
                rootFS = os.DirFS(packPath)
                packSource = versionToFetch
                fmt.Printf("Downloaded pack from GitHub releases\n")
            } else {
                // Fall back to embedded base
                fmt.Printf("Download failed (%v), using embedded base pack\n", err)
                rootFS, err = pack.GetEmbeddedBaseFS()
                if err != nil {
                    return fmt.Errorf("failed to load embedded pack: %w", err)
                }
                packSource = "embedded-base"
            }
        } else {
            // --offline flag: use embedded base only
            fmt.Println("Using embedded base pack (offline mode)")
            rootFS, err = pack.GetEmbeddedBaseFS()
            if err != nil {
                return fmt.Errorf("failed to load embedded pack: %w", err)
            }
            packSource = "embedded-base"
        }
        
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
            installedVersion := packSource
            if initVersion != "" { installedVersion = initVersion }
            if err := manifest.WriteWithStacks(files, installedVersion, choices.Stacks); err != nil { return err }
        }
        installedVersion := packSource
        if initVersion != "" { installedVersion = initVersion }
        fmt.Printf("\nCodo %s initialized. See .claude/.codo-report/conflicts.txt if any.\n", installedVersion)
        return nil
    },
}

func init() {
    initCmd.Flags().StringVar(&initVersion, "version", "", "Pack version to download (e.g. v1.2.0)")
    initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview only; do not write files")
    initCmd.Flags().StringVar(&initStacks, "stacks", "", "Comma-separated stacks (skip TUI)")
    initCmd.Flags().BoolVar(&initNoTUI, "no-tui", false, "Don't show the TUI wizard")
    initCmd.Flags().BoolVar(&initOffline, "offline", false, "Use embedded base pack only (no download)")
}
