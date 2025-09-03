package pack

import (
    "io/fs"
    "path/filepath"
    "slices"
    "strings"
)

// FilesFromDotclaudeFS composes .claude/base + selected .claude/stacks/<stack>
// and includes top-level files (CLAUDE.md, docs/**) from the provided FS root.
// The returned RelPath is the project-relative destination path.
func FilesFromDotclaudeFS(root fs.FS, stacks []string) ([]File, error) {
    const baseRoot = ".claude/base"
    const stacksRoot = ".claude/stacks"

    index := map[string]string{} // rel -> FS path

    // 1) base contents
    if err := fs.WalkDir(root, baseRoot, func(p string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        rel := strings.TrimPrefix(p, baseRoot+"/")
        rel = filepath.ToSlash(filepath.Join(".claude", rel))
        index[rel] = p
        return nil
    }); err != nil { return nil, err }

    // Normalize stacks to allowed ones only
    want := make([]string, 0, len(stacks))
    for _, s := range stacks {
        if slices.Contains(allowedStacks, s) {
            want = append(want, s)
        }
    }

    // 2) overlays
    for _, s := range want {
        base := filepath.Join(stacksRoot, s)
        _ = fs.WalkDir(root, base, func(p string, d fs.DirEntry, err error) error {
            if err != nil { return err }
            if d.IsDir() { return nil }
            rel := strings.TrimPrefix(p, base+"/")
            rel = filepath.ToSlash(filepath.Join(".claude", rel))
            index[rel] = p // overlay wins
            return nil
        })
    }

    // 3) Top-level non-.claude files (e.g., CLAUDE.md, docs/**)
    _ = fs.WalkDir(root, ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil { return nil }
        if d.IsDir() { return nil }
        if strings.HasPrefix(p, ".claude/") { return nil }
        index[filepath.ToSlash(p)] = p
        return nil
    })

    out := make([]File, 0, len(index))
    for rel, p := range index {
        relLocal, pLocal := rel, p
        out = append(out, File{
            RelPath: relLocal,
            Read:    func() ([]byte, error) { return fs.ReadFile(root, pLocal) },
        })
    }
    slices.SortFunc(out, func(a, b File) int { return strings.Compare(a.RelPath, b.RelPath) })
    return out, nil
}

