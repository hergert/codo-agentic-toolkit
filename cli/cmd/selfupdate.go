package cmd

import (
    "fmt"
    "os"

    semver "github.com/Masterminds/semver/v3"
    "github.com/spf13/cobra"
    "github.com/rhysd/go-github-selfupdate/selfupdate"
)

var selfUpdateCmd = &cobra.Command{
    Use:   "self-update",
    Short: "Update the codo binary to the latest release",
    RunE: func(cmd *cobra.Command, args []string) error {
        latest, found, err := selfupdate.DetectLatest("youruser/codo")
        if err != nil { return err }
        v := semver.MustParse(version)
        if !found || !latest.Version.GreaterThan(v) {
            fmt.Println("codo is up to date:", version)
            return nil
        }
        // Note: For extra safety, verify checksum here using latest.AssetChecksum if provided.
        if err := selfupdate.UpdateTo(latest, os.Args[0]); err != nil { return err }
        fmt.Println("Updated to", latest.Version)
        return nil
    },
}
