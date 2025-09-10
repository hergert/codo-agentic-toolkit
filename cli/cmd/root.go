package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// Set via -ldflags "-X github.com/hergert/codo-agentic-toolkit/cli/cmd.version=vX.Y.Z"
var version = "dev"

func Execute() { cobra.CheckErr(rootCmd.Execute()) }

var rootCmd = &cobra.Command{
    Use:   "codo",
    Short: "Manage the Codo Agentic Toolkit in any repo",
    Long:  "Install, update, remove, and check status of the Codo toolkit with safe conflict handling.",
}

func init() {
    rootCmd.Version = version
    rootCmd.AddCommand(initCmd, updateCmd, removeCmd, statusCmd, doctorCmd, selfUpdateCmd)
}

func abortIf(cond bool, msg string) {
    if cond {
        fmt.Fprintln(os.Stderr, msg)
        os.Exit(1)
    }
}

