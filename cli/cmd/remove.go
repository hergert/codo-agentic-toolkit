package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var removeDry bool

var removeCmd = &cobra.Command{
    Use:   "remove",
    Short: "Remove the toolkit (backup first)",
    RunE: func(cmd *cobra.Command, args []string) error {
        // TODO: implement backup to .codo-backup/<ts>/ and cleanup
        fmt.Println("Remove (skeleton): will backup files and clean manifest.")
        return nil
    },
}

func init() {
    removeCmd.Flags().BoolVar(&removeDry, "dry-run", false, "Preview only; do not write files")
}

