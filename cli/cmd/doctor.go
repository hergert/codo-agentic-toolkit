package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check environment requirements",
    RunE: func(cmd *cobra.Command, args []string) error {
        checks := []struct{ name, bin string }{
            {"git", "git"},
        }
        for _, c := range checks {
            if _, err := exec.LookPath(c.bin); err != nil { fmt.Println("✗", c.name) } else { fmt.Println("✓", c.name) }
        }
        
        // Check Python version
        if err := checkPython(); err != nil {
            fmt.Println("✗ Python 3:", err)
        }
        
        for _, h := range []string{".claude/hooks/pre_tool_use.py", ".claude/hooks/post_tool_use.py", ".claude/hooks/user_prompt_submit.py"} {
            if fi, err := os.Stat(h); err == nil && fi.Mode()&0o111 != 0 { fmt.Println("✓ exec", h) } else { fmt.Println("✗ exec", h) }
        }
        return nil
    },
}

func checkPython() error {
    cmd := exec.Command("python3", "--version")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("Python 3 not found (required for hooks)")
    }
    version := strings.TrimSpace(string(output))
    fmt.Printf("✓ %s\n", version)
    // Could parse version and warn if < 3.6, but basic check is sufficient
    return nil
}
