package cmd

import (
    "crypto/sha256"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/hergert/codo-agentic-toolkit/cli/internal/manifest"
)

var strictFlag bool

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show installed version and drift vs manifest",
    RunE: func(cmd *cobra.Command, args []string) error {
        if !manifest.Exists() {
            fmt.Println("codo: not installed")
            return nil
        }
        m, err := manifest.Open()
        if err != nil { return err }
        var drift []string
        for _, ent := range m.Files {
            b, err := os.ReadFile(ent.Path)
            if err != nil { drift = append(drift, "missing "+ent.Path); continue }
            if fmt.Sprintf("%x", sha256.Sum256(b)) != ent.SHA256 { drift = append(drift, "~ "+ent.Path) }
        }
        fmt.Println("Installed version:", m.Version)
        if len(drift) == 0 {
            fmt.Println("No drift")
        } else {
            fmt.Println("Drift:")
            for _, d := range drift { fmt.Println(" ", d) }
            if strictFlag {
                return fmt.Errorf("drift detected")
            }
        }
        return nil
    },
}

func init() {
    statusCmd.Flags().BoolVar(&strictFlag, "strict", false, "Exit non-zero if drift exists")
}
