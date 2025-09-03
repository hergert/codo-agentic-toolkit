package cmd

import (
    "fmt"

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
        to := pack.VersionOrDefault(updateTo)
        fmt.Println("Update planner (skeleton): target", to, "— fill in manifest comparison here.")
        return nil
    },
}

func init() {
    updateCmd.Flags().StringVar(&updateTo, "to", "", "Version to update to (default embedded)")
    updateCmd.Flags().BoolVar(&updateDry, "dry-run", false, "Preview only")
}

