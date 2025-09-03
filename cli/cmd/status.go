package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/youruser/codo/internal/manifest"
)

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
        fmt.Println("Installed version:", m.Version)
        fmt.Println("Files tracked:", len(m.Files))
        // TODO: compute drifts
        return nil
    },
}

