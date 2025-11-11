package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/doctor"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment requirements",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		summary := doctor.Collect(root)
		fmt.Println("Dev tools checklist:")
		for _, item := range summary.Items {
			fmt.Printf(" - %s: %s\n", item.Label, item.Detail)
		}
		return nil
	},
}
