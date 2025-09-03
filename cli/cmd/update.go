package cmd

import (
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/youruser/codo/internal/manifest"
    "github.com/youruser/codo/internal/pack"
)

var updateTo string
var updateDry bool

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update the toolkit (only overwrite files unchanged since install)",
    RunE: func(cmd *cobra.Command, args []string) error {
        abortIf(!manifest.Exists(), "No manifest found. Run `codo init` first.")
        m, err := manifest.Open()
        if err != nil { return err }

        if _, err := os.Stat("dotclaude"); err != nil {
            return fmt.Errorf("dotclaude pack not found; expected ./dotclaude directory")
        }
        rootFS := os.DirFS("dotclaude")
        files, err := pack.FilesFromDotclaudeFS(rootFS, m.Stacks)
        if err != nil { return err }

        // Build map of new contents
        newMap := map[string][]byte{}
        for _, f := range files {
            b, err := f.Read(); if err == nil { newMap[f.RelPath] = b }
        }

        // Ensure report dir exists
        _ = os.MkdirAll(filepath.Join(".claude", ".codo-report"), 0o755)

        for _, ent := range m.Files {
            dst := ent.Path
            nb, ok := newMap[dst]
            if !ok {
                // Skip if new pack doesn't contain this file
                continue
            }
            cur, err := os.ReadFile(dst)
            if err != nil {
                // Missing → treat as clean overwrite
                fmt.Println("+ " + dst)
                if !updateDry { if err := os.WriteFile(dst, nb, 0o644); err != nil { return err } }
                continue
            }
            curHash := fmt.Sprintf("%x", sha256.Sum256(cur))
            if curHash == ent.SHA256 {
                // clean → overwrite
                fmt.Println("~ " + dst)
                if !updateDry { if err := os.WriteFile(dst, nb, 0o644); err != nil { return err } }
            } else {
                // diverged → write .codo.new
                out := dst + ".codo.new"
                fmt.Println("! conflict → " + out)
                if !updateDry { if err := os.WriteFile(out, nb, 0o644); err != nil { return err } }
            }
        }
        if !updateDry {
            newVersion := updateTo
            if newVersion == "" { newVersion = m.Version }
            if err := manifest.WriteWithStacks(files, newVersion, m.Stacks); err != nil { return err }
        }
        return nil
    },
}

func init() {
    updateCmd.Flags().StringVar(&updateTo, "to", "", "Version/tag to update to (e.g. v1.2.0)")
    updateCmd.Flags().BoolVar(&updateDry, "dry-run", false, "Preview only")
}
