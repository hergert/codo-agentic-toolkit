package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/spf13/cobra"
    "github.com/youruser/codo/internal/manifest"
)

var removeDry bool

var removeCmd = &cobra.Command{
    Use:   "remove",
    Short: "Remove the toolkit (backup first)",
    RunE: func(cmd *cobra.Command, args []string) error {
        abortIf(!manifest.Exists(), "No manifest found. Nothing to remove.")
        m, err := manifest.Open()
        if err != nil { return err }
        ts := time.Now().UTC().Format("20060102-150405")
        backup := filepath.Join(".codo-backup", ts)
        if !removeDry { if err := os.MkdirAll(backup, 0o755); err != nil { return err } }
        for _, ent := range m.Files {
            if _, err := os.Stat(ent.Path); err == nil {
                fmt.Println("- " + ent.Path)
                if !removeDry {
                    dest := filepath.Join(backup, ent.Path)
                    if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil { return err }
                    if err := os.Rename(ent.Path, dest); err != nil { return err }
                    _ = os.Remove(ent.Path + ".codo.new")
                }
            }
        }
        if !removeDry {
            _ = os.Remove(".claude/.codo-manifest.json")
        }
        fmt.Println("Backup at", backup)
        return nil
    },
}

func init() {
    removeCmd.Flags().BoolVar(&removeDry, "dry-run", false, "Preview only; do not write files")
}
