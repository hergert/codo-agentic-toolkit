# Codo CLI — Essentials Plan

## Goal
Provide a single static Go CLI (`codo`) to safely install, update, remove, inspect, and self‑update the Codo Agentic Toolkit in any repo. It must never overwrite user changes without review, offer an interactive TUI for stack selection and diffs, track installed files via a manifest with checksums, and distribute cross‑platform binaries (Homebrew/Scoop/Winget) with a simple curl installer.

## Acceptance Checks
- `brew install <tap>/codo` installs; `codo version` shows tagged version.
- `codo init` in a temp repo writes `.claude/.codo-manifest.json`, installs selected stacks, marks hooks executable, and logs conflicts without overwriting (creates `*.codo.new`).
- `codo update --to v1.1` only overwrites files unchanged since install; edits become `*.codo.new`; manifest updated.
- `codo remove` backs up installed files to `.codo-backup/<ts>/` and removes manifest; `codo status` shows “not installed”.
- `codo doctor` reports required tooling and permissions; `codo self-update` replaces the binary with a newer release after checksum verification.

## Files To Touch (exact list; intent)
- `cli/go.mod`: Go module init for the CLI.
- `cli/main.go`: entrypoint calling `cmd.Execute()`.
- `cli/cmd/root.go`: Cobra root command and global flags.
- `cli/cmd/init.go`: implements `codo init` (non‑interactive + TUI hook).
- `cli/cmd/update.go`: implements `codo update` (compare manifest, apply safe updates).
- `cli/cmd/remove.go`: implements `codo remove` (backup + cleanup).
- `cli/cmd/status.go`: implements `codo status` (manifest + drift report).
- `cli/cmd/doctor.go`: implements `codo doctor` (env/tool checks).
- `cli/cmd/selfupdate.go`: implements `codo self-update` (GitHub releases).
- `cli/internal/pack/pack.go`: embed versioned templates `templates/v1.1/**` via `embed.FS`.
- `cli/internal/fsops/fsops.go`: safe copy/move, conflict handling, chmod hooks.
- `cli/internal/manifest/manifest.go`: read/write `.claude/.codo-manifest.json` with SHA256s.
- `cli/internal/diff/diff.go`: render text diffs (using `go-diff`).
- `cli/internal/tui/*.go`: Bubble Tea models (stack picker, preview, diff, progress).
- `cli/internal/selfupdate/selfupdate.go`: wrapper around self‑update lib + checksum verify.
- `cli/internal/pack/templates/v1.1/**`: embedded toolkit files by stack.
- `.goreleaser.yaml`: builds, archives, Homebrew/Scoop/Winget config.
- `docs/CLI.md`: usage docs and install commands.

## Diff Outline (pseudo‑code)
- `main.main() -> cmd.Execute()`
- `cmd.rootCmd`: `Use:"codo"`; adds persistent flags; registers subcommands.
- `cmd.initCmd.RunE()`: discover pack -> TUI select stacks -> preview -> `fsops.ApplyInit()` -> `manifest.Write()` -> report conflicts.
- `cmd.updateCmd.RunE()`: `manifest.Load()` -> calculate checksums -> plan (auto vs conflict) -> apply -> update manifest.
- `cmd.removeCmd.RunE()`: backup to `.codo-backup/<ts>/` -> remove files -> delete manifest -> clean byproducts.
- `cmd.statusCmd.RunE()`: read manifest -> compare current checksums -> print drift and available pack version.
- `cmd.doctorCmd.RunE()`: check git/python3/permissions/hooks exec -> print pass/fail.
- `cmd.selfUpdateCmd.RunE()`: fetch latest release -> verify checksum -> atomic replace.
- `fsops.CopySafe(src,dst)`: if identical skip; if absent copy; if differs write `dst.codo.new` + log.
- `manifest.Struct{Version,InstalledAt,Files[]:Path,SHA256}`; `ComputeSHA256(path)`.
- `tui`: stack picker (checkbox), preview list, diff viewer, progress screen.

## Risks (and detection)
- Overwrite user edits: mitigated by `*.codo.new` policy; verify via update acceptance test and conflict report presence.
- Manifest drift or corruption: detect with `status` diff; fallback to re‑compute checksums; test corrupted file path handling.
- Platform packaging issues: verify CI release matrix; install smoke tests on macOS/Linux/Windows runners.
- Self‑update integrity: verify checksums/signatures; simulate MITM by altering checksum to ensure failure.
- Performance on large repos: scope writes to project root only; measure `init`/`update` wall time in CI sample repos.

## Test Strategy (contracts)
- `cli/internal/fsops/fsops_test.go`: CopySafe behaviors (identical/absent/different -> `*.codo.new`).
- `cli/internal/manifest/manifest_test.go`: read/write round‑trip; checksum correctness; drift detection.
- `cli/internal/diff/diff_test.go`: basic diff rendering contract (non‑empty, expected markers).
- `cli/cmd/init_update_integ_test.go`: temp repo end‑to‑end for `init` + `update` (no overwrite; manifest written/updated).
- `cli/cmd/remove_status_integ_test.go`: `remove` backup/cleanup; `status` pre/post.
- `cli/cmd/selfupdate_test.go`: dry‑run path validates release selection and checksum gate (mocked).

## Out of Scope (YAGNI)
- Managing project dependencies (`npm`, `go`, `pip`, etc.).
- Telemetry/analytics and remote template fetching.
- Plugin system or third‑party extension hooks.
- Executing project scripts beyond setting executable bits on hook files.
- Rich side‑by‑side syntax‑highlighted diffs (basic diff first).

## Starter (Ready-to-Paste)

Why these libs (essentials)
- Cobra: de‑facto CLI framework; clear subcommand patterns; widely adopted.
- Bubble Tea (+Bubbles): polished TUI with list/pickers/spinners for modern UX.
- embed: bundle templates in the binary for zero‑deps install.
- GoReleaser: single config to publish binaries and brew/scoop/winget manifests.

0) Repo layout (new cli/ module)
```
cli/
├── cmd/
│   ├── root.go
│   ├── init.go
│   ├── update.go
│   ├── remove.go
│   ├── status.go
│   ├── doctor.go
│   └── selfupdate.go
├── internal/
│   ├── tui/
│   │   └── init_wizard.go
│   ├── fsops/
│   │   └── fsops.go
│   ├── manifest/
│   │   └── manifest.go
│   ├── pack/
│   │   └── pack.go              # //go:embed v1.1 pack
│   └── diff/
│       └── diff.go              # tiny wrapper around go-diff (optional)
├── main.go
├── go.mod
├── go.sum
└── .goreleaser.yaml
```

1) go.mod
```
module github.com/youruser/codo

go 1.22

require (
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/bubbletea v0.26.0
    github.com/sergi/go-diff v1.3.1
    github.com/spf13/cobra v1.8.0
)
```

2) main + Cobra roots

cli/main.go
```
package main

import "github.com/youruser/codo/cmd"

func main() { cmd.Execute() }
```

cli/cmd/root.go
```
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    // Set at build time via -ldflags "-X github.com/youruser/codo/cmd.version=vX.Y.Z"
    version = "dev"
)

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
```

3) Subcommands (skeletons)

cli/cmd/init.go
```
package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    "github.com/youruser/codo/internal/fsops"
    "github.com/youruser/codo/internal/manifest"
    "github.com/youruser/codo/internal/pack"
    "github.com/youruser/codo/internal/tui"
)

var initVersion string
var initDryRun bool

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Install the toolkit into this repository (safe-by-default)",
    RunE: func(cmd *cobra.Command, args []string) error {
        root, _ := os.Getwd()
        ctx := context.Background()

        choices, err := tui.RunInitWizard(ctx)
        if err != nil { return err }
        if !choices.Confirmed {
            fmt.Println("aborted")
            return nil
        }

        files, err := pack.Files(initVersion, choices.Stacks)
        if err != nil { return err }

        reportDir := filepath.Join(".claude", ".codo-report")
        _ = os.MkdirAll(reportDir, 0o755)

        // Copy safely (or simulate with --dry-run). fsops prints +/=!/conflict lines.
        for _, f := range files {
            if err := fsops.CopySafe(f, root, initDryRun); err != nil { return err }
        }
        if !initDryRun {
            if err := fsops.ChmodHooks(); err != nil { return err }
            if err := manifest.Write(files, pack.VersionOrDefault(initVersion)); err != nil { return err }
        }
        fmt.Printf("\nCodo %s initialized. See .claude/.codo-report/conflicts.txt if any.\n", pack.VersionOrDefault(initVersion))
        return nil
    },
}

func init() {
    initCmd.Flags().StringVar(&initVersion, "version", "", "Pack version (default embedded)")
    initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview only; do not write files")
}
```

4) TUI (Bubble Tea) — labeled options, exact keys

cli/internal/tui/init_wizard.go
```
package tui

import (
    "context"
    "fmt"
    "sort"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

type InitResult struct {
    Stacks    []string // exact keys to feed into pack.Files()
    Confirmed bool
}

type option struct {
    Key   string // exact key for pack.Files()
    Label string // pretty label in the TUI
}

var opts = []option{
    {Key: "cloudflare-workers", Label: "Cloudflare Workers (wrangler)"},
    {Key: "supabase",           Label: "Supabase (DB & functions)"},
    {Key: "trigger.dev",        Label: "trigger.dev (background jobs)"},
    {Key: "go",                 Label: "Go"},
    {Key: "typescript",         Label: "TypeScript / Node"},
    {Key: "python",             Label: "Python"},
    {Key: "flutter",            Label: "Flutter"},
}

type model struct {
    cursor   int
    selected map[int]bool
    quit     bool
    confirm  bool
}

func (m model) Init() tea.Cmd { return tea.Tick(time.Millisecond*10, func(time.Time) tea.Msg { return nil }) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch key := msg.(type) {
    case tea.KeyMsg:
        switch key.String() {
        case "ctrl+c", "q":
            m.quit, m.confirm = true, false
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 { m.cursor-- }
        case "down", "j":
            if m.cursor < len(opts)-1 { m.cursor++ }
        case " ":
            if m.selected == nil { m.selected = map[int]bool{} }
            m.selected[m.cursor] = !m.selected[m.cursor]
        case "a": // select all
            if m.selected == nil { m.selected = map[int]bool{} }
            for i := range opts { m.selected[i] = true }
        case "n": // select none
            m.selected = map[int]bool{}
        case "enter":
            m.confirm, m.quit = true, true
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    var b strings.Builder
    fmt.Fprintln(&b, "Select stacks to include (space to toggle):\n")
    for i, o := range opts {
        cursor := " "
        if m.cursor == i { cursor = ">" }
        box := " "
        if m.selected[i] { box = "x" }
        fmt.Fprintf(&b, " %s [%s] %s\n", cursor, box, o.Label)
    }
    fmt.Fprintln(&b, "\n[↑/↓/j/k] move   [space] toggle   [a] all   [n] none   [enter] continue   [q] abort")
    return b.String()
}

func RunInitWizard(_ context.Context) (InitResult, error) {
    m := model{selected: map[int]bool{}}
    pm := tea.NewProgram(m)
    res, err := pm.Run()
    if err != nil { return InitResult{}, err }

    out := res.(model)
    if out.quit && !out.confirm {
        return InitResult{Stacks: nil, Confirmed: false}, nil
    }
    // Build stable list of selected keys
    var keys []string
    for i, on := range out.selected {
        if on { keys = append(keys, opts[i].Key) }
    }
    sort.Strings(keys)
    return InitResult{Stacks: keys, Confirmed: true}, nil
}
```

5) Embed pack

cli/internal/pack/pack.go
```
package pack

import (
    "embed"
    "fmt"
    "io/fs"
    "path/filepath"
)

var defaultVersion = "v1.1"

//go:embed templates/v1.1/**
var fsV11 embed.FS

type File struct {
    RelPath string
    Read    func() ([]byte, error)
}

func VersionOrDefault(v string) string {
    if v == "" { return defaultVersion }
    return v
}

// Files returns files for the (version, stacks). Stack filtering can be added later.
func Files(version string, stacks []string) ([]File, error) {
    switch VersionOrDefault(version) {
    case "v1.1":
        var out []File
        err := fs.WalkDir(fsV11, "templates/v1.1", func(p string, d fs.DirEntry, err error) error {
            if err != nil { return err }
            if d.IsDir() { return nil }
            rel := p[len("templates/v1.1/"):]
            rel = filepath.ToSlash(rel)
            out = append(out, File{
                RelPath: rel,
                Read: func() ([]byte, error) { return fs.ReadFile(fsV11, p) },
            })
            return nil
        })
        return out, err
    default:
        return nil, fmt.Errorf("unknown pack version %q", version)
    }
}
```

6) Safe fs ops + manifest (minimal)

cli/internal/fsops/fsops.go
```
package fsops

import (
    "crypto/sha256"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/youruser/codo/internal/pack"
)

func sha256File(path string) (string, error) {
    f, err := os.Open(path); if err != nil { return "", err }
    defer f.Close()
    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil { return "", err }
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func CopySafe(f pack.File, projectRoot string, dry bool) error {
    dst := filepath.Join(projectRoot, f.RelPath)
    if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil { return err }
    srcBytes, err := f.Read(); if err != nil { return err }

    if _, err := os.Stat(dst); err == nil {
        tmp := dst + ".codo.new"
        // Only write .codo.new if different
        curHash, _ := sha256File(dst)
        newHash := fmt.Sprintf("%x", sha256.Sum256(srcBytes))
        if curHash == newHash { fmt.Println("= " + f.RelPath); return nil }
        fmt.Println("! conflict → " + tmp)
        if dry { return nil }
        return os.WriteFile(tmp, srcBytes, 0o644)
    }

    fmt.Println("+ " + f.RelPath)
    if dry { return nil }
    return os.WriteFile(dst, srcBytes, 0o644)
}

func ChmodHooks() error {
    files := []string{
        filepath.Join(".claude", "hooks", "pre_tool_use.py"),
        filepath.Join(".claude", "hooks", "post_tool_use.py"),
        filepath.Join(".claude", "hooks", "user_prompt_submit.py"),
    }
    for _, p := range files {
        if _, err := os.Stat(p); err == nil {
            if err := os.Chmod(p, 0o755); err != nil { return err }
        }
    }
    return nil
}
```

cli/internal/manifest/manifest.go
```
package manifest

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/youruser/codo/internal/pack"
)

const filePath = ".claude/.codo-manifest.json"

type Entry struct {
    Path   string `json:"path"`
    SHA256 string `json:"sha256"`
}
type Manifest struct {
    Version     string  `json:"version"`
    InstalledAt string  `json:"installed_at"`
    Files       []Entry `json:"files"`
}

func Exists() bool {
    _, err := os.Stat(filePath)
    return err == nil
}

func Write(files []pack.File, version string) error {
    if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil { return err }
    var entries []Entry
    for _, f := range files {
        // use on-disk hash if placed; else hash of new content
        dst := f.RelPath
        var sum string
        if b, err := os.ReadFile(dst); err == nil {
            sum = fmt.Sprintf("%x", sha256.Sum256(b))
        } else {
            b, err := f.Read(); if err != nil { return err }
            sum = fmt.Sprintf("%x", sha256.Sum256(b))
        }
        entries = append(entries, Entry{Path: dst, SHA256: sum})
    }
    m := Manifest{Version: version, InstalledAt: "", Files: entries}
    buf, _ := json.MarshalIndent(m, "", "  ")
    return os.WriteFile(filePath, buf, 0o644)
}

func Open() (Manifest, error) {
    var m Manifest
    f, err := os.Open(filePath); if err != nil { return m, err }
    defer f.Close()
    b, _ := io.ReadAll(f)
    err = json.Unmarshal(b, &m)
    return m, err
}
```

7) GoReleaser (minimal)

cli/.goreleaser.yaml
```
project_name: codo

before:
  hooks:
    - go mod tidy

builds:
  - id: codo
    main: ./main.go
    env:
      - CGO_ENABLED=0
    flags:
      - -trimpath
    ldflags:
      - -s -w -X github.com/youruser/codo/cmd.version={{.Version}}
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]

archives:
  - builds: [codo]
    format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"

release:
  github:
    owner: youruser
    name: codo
  draft: false
  prerelease: false

brews:
  - name: codo
    repository:
      owner: youruser
      name: homebrew-tap
    commit_author:
      name: youruser
      email: you@example.com
    directory: Formula
    homepage: "https://github.com/youruser/codo"
    description: "Modern CLI to install/update/remove the Codo Agentic Toolkit."
    install: |
      bin.install "codo"

scoops:
  - name: codo
    bucket:
      owner: youruser
      name: scoop-bucket
    homepage: "https://github.com/youruser/codo"
    description: "Modern CLI to manage the Codo Agentic Toolkit."

winget:
  - name: youruser.codo
    publisher: "Your Name"
    short_description: "Modern CLI to manage the Codo Agentic Toolkit."
    homepage: "https://github.com/youruser/codo"
```

8) Quick build & run
```
# inside cli/
go mod tidy

# Run (dev)
go run ./ main.go --help
go run ./ init --dry-run

# Build
go build -ldflags "-s -w -X github.com/youruser/codo/cmd.version=v0.1.0" -o codo

# First release (after tagging v0.1.0):
# 1) Install goreleaser locally (or via CI)
# 2) export GITHUB_TOKEN=... (repo scope)
goreleaser release --clean
```

## Selective Stacks (Overlay Merge)

Directory convention (inside CLI repo)
```
cli/internal/pack/templates/v1.1/
├── core/                     # always included
│   ├── .claude/**            # full v1.1 toolkit (hooks, settings, commands…)
│   ├── CLAUDE.md
│   ├── README.md
│   └── docs/**               # knowledge-base, tasks README, etc.
└── stacks/                   # optional overlays; installed only if selected
    ├── cloudflare-workers/
    │   ├── .claude/snippets/hooks/cloudflare.prod-gate.json
    │   └── docs/stack-notes/cloudflare.md
    ├── supabase/
    │   ├── .claude/snippets/hooks/supabase.db-gate.json
    │   └── docs/stack-notes/supabase.md
    ├── trigger.dev/
    │   ├── .claude/snippets/hooks/trigger.dev.deploy-gate.json
    │   └── docs/stack-notes/trigger-dev.md
    ├── go/
    │   └── docs/stack-notes/go.md
    ├── typescript/
    │   └── docs/stack-notes/typescript.md
    ├── python/
    │   └── docs/stack-notes/python.md
    └── flutter/
        ├── .claude/snippets/hooks/flutter.release-gate.json
        └── docs/stack-notes/flutter.md
```

Keep `core/` as your full v1.1 pack. Place stack-specific additions/overrides under `stacks/<name>/...`. If a stack provides a file with the same relative path as core, the stack version wins (overlay replaces core).

Stack names returned by the wizard: `cloudflare-workers`, `supabase`, `trigger.dev`, `go`, `typescript`, `python`, `flutter`.

Overlay-aware `pack.Files()`
```
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
```

Why this works
- embed.FS: read-only virtual filesystem; walk with `fs.WalkDir` to enumerate files.
- Stable keys: `filepath.ToSlash` normalizes separators across OSes.
- Overlay merge: a map keyed by relative path lets stacks cleanly replace core files when needed.
- Versioning: duplicate the `//go:embed` and `filesVXYZ` helper per version and switch in `Files()`.

Wizard ↔ stacks mapping
- The TUI should emit any subset of: `cloudflare-workers`, `supabase`, `trigger.dev`, `go`, `typescript`, `python`, `flutter`.
- Pass directly to `pack.Files(version, stacks)`; unrecognized names are ignored.

Notes on JSON overrides
- Prefer additive overlays (e.g., add snippets under `.claude/snippets/hooks/…`).
- If you must change a JSON file, replace the whole file via overlay. Avoid partial merges in v1.
- A future enhancement could add a targeted JSON merge step if truly necessary.

Quick test (no TUI)
```
files, err := pack.Files("", []string{"cloudflare-workers", "go"})
if err != nil { panic(err) }
for _, f := range files { fmt.Println(f.RelPath) }
```
You should see all `core/` files plus any under `stacks/cloudflare-workers/**` and `stacks/go/**`; if a stack overrides `README.md`, its version appears.
