package pack

type File struct {
    RelPath string
    Read    func() ([]byte, error)
}

// Allowed stack keys returned by the TUI/flags.
var allowedStacks = []string{
    "cloudflare-workers",
    "supabase",
    "trigger.dev",
    "go",
    "typescript",
    "python",
    "flutter",
}

