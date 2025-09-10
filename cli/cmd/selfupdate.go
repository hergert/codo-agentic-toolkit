package cmd

import (
    "fmt"
    "os"

    "github.com/blang/semver"
    "github.com/rhysd/go-github-selfupdate/selfupdate"
    "github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
    Use:   "self-update",
    Short: "Update the codo binary to the latest release",
    RunE: func(cmd *cobra.Command, args []string) error {
        latest, found, err := selfupdate.DetectLatest("hergert/codo-agentic-toolkit")
        if err != nil { return err }
        v, err := semver.Parse(version)
        if err != nil { return err }
        if !found || !latest.Version.GT(v) {
            fmt.Println("codo is up to date:", version)
            return nil
        }
        // Note: For extra safety, verify checksum here using latest checksum info if available.
        exe, err := os.Executable()
        if err != nil {
            exe = os.Args[0]
        }
        if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil { return err }
        fmt.Println("Updated to", latest.Version)
        return nil
    },
}
