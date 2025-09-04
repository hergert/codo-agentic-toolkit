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
        latest, found, err := selfupdate.DetectLatest("youruser/codo")
        if err != nil { return err }
        v, err := semver.Parse(version)
        if err != nil { return err }
        if !found || !latest.Version.GT(v) {
            fmt.Println("codo is up to date:", version)
            return nil
        }
        // Note: For extra safety, verify checksum here using latest checksum info if available.
        if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0]); err != nil { return err }
        fmt.Println("Updated to", latest.Version)
        return nil
    },
}
