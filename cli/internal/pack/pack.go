package pack

import (
    "embed"
    "fmt"
    "io/fs"
    "path/filepath"
    "slices"
    "strings"
)

// Update this when you cut a new embedded pack.
var defaultVersion = "v1.1"

// Embed the whole version folder; we'll filter in code.
//go:embed templates/v1.1/**
var fsV11 embed.FS

type File struct {
    RelPath string
    Read    func() ([]byte, error)
}

// Only allow known stacks (avoid path traversal).
var allowedStacks = []string{
    "cloudflare-workers",
    "supabase",
    "trigger.dev",
    "go",
    "typescript",
    "python",
    "flutter",
}

func VersionOrDefault(v string) string {
    if v == "" {
        return defaultVersion
    }
    return v
}

// Files returns all core files + overlays for the selected stacks.
// Overlay rule: later sources override earlier ones (stack file replaces core if same RelPath).
func Files(version string, stacks []string) ([]File, error) {
    ver := VersionOrDefault(version)
    switch ver {
    case "v1.1":
        return filesV11(stacks)
    default:
        return nil, fmt.Errorf("unknown pack version %q", ver)
    }
}

func filesV11(stacks []string) ([]File, error) {
    root := "templates/" + defaultVersion
    coreRoot := root + "/core"
    stacksRoot := root + "/stacks"

    // Normalize requested stacks and keep only allowed.
    want := make([]string, 0, len(stacks))
    for _, s := range stacks {
        if slices.Contains(allowedStacks, s) {
            want = append(want, s)
        }
    }

    type entry struct{ p string } // path inside embed FS
    // Use a map keyed by relative path to allow overrides.
    index := map[string]entry{}

    // Helper to add files from a subtree into the index with precedence.
    addTree := func(treeRoot string) error {
        return fs.WalkDir(fsV11, treeRoot, func(p string, d fs.DirEntry, err error) error {
            if err != nil { return err }
            if d.IsDir() { return nil }
            // Compute project-relative path under core/ or stacks/<name>/
            rel := strings.TrimPrefix(p, treeRoot+"/")
            rel = filepath.ToSlash(rel) // stable slashes across OSes
            index[rel] = entry{p: p}    // later calls override earlier ones
            return nil
        })
    }

    // 1) add core
    if err := addTree(coreRoot); err != nil { return nil, err }
    // 2) add selected stacks (overlay precedence)
    for _, s := range want {
        if err := addTree(filepath.Join(stacksRoot, s)); err != nil { return nil, err }
    }

    // Build stable output slice
    out := make([]File, 0, len(index))
    for rel, e := range index {
        pLocal := e.p
        relLocal := rel
        out = append(out, File{
            RelPath: relLocal,
            Read:    func() ([]byte, error) { return fs.ReadFile(fsV11, pLocal) },
        })
    }
    // Optional: sort by RelPath for deterministic order
    slices.SortFunc(out, func(a, b File) int { return strings.Compare(a.RelPath, b.RelPath) })
    return out, nil
}

