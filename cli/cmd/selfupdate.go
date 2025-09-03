package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
    Use:   "self-update",
    Short: "Update the codo binary to the latest release",
    RunE: func(cmd *cobra.Command, args []string) error {
        // TODO: integrate self-update library
        fmt.Println("Self-update (skeleton): would check GitHub releases and update.")
        return nil
    },
}

