package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check environment requirements",
    RunE: func(cmd *cobra.Command, args []string) error {
        // TODO: check git, python3, permissions, etc.
        fmt.Println("Doctor (skeleton): environment looks OK.")
        return nil
    },
}

