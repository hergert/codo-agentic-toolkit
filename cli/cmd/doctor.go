package cmd

import (
    "fmt"
    "os"
    "os/exec"

    "github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check environment requirements",
    RunE: func(cmd *cobra.Command, args []string) error {
        checks := []struct{ name, bin string }{
            {"git", "git"},
            {"python3", "python3"},
        }
        for _, c := range checks {
            if _, err := exec.LookPath(c.bin); err != nil { fmt.Println("✗", c.name) } else { fmt.Println("✓", c.name) }
        }
        for _, h := range []string{".claude/hooks/pre_tool_use.py", ".claude/hooks/post_tool_use.py", ".claude/hooks/user_prompt_submit.py"} {
            if fi, err := os.Stat(h); err == nil && fi.Mode()&0o111 != 0 { fmt.Println("✓ exec", h) } else { fmt.Println("✗ exec", h) }
        }
        return nil
    },
}
