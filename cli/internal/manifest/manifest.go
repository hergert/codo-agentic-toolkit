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
    Stacks      []string `json:"stacks,omitempty"`
}

func Exists() bool {
    _, err := os.Stat(filePath)
    return err == nil
}

func Write(files []pack.File, version string) error {
    return WriteWithStacks(files, version, nil)
}

func WriteWithStacks(files []pack.File, version string, stacks []string) error {
    if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
        return err
    }
    var entries []Entry
    for _, f := range files {
        // use on-disk hash if placed; else hash of new content
        dst := f.RelPath
        var sum string
        if b, err := os.ReadFile(dst); err == nil {
            sum = fmt.Sprintf("%x", sha256.Sum256(b))
        } else {
            b, err := f.Read()
            if err != nil {
                return err
            }
            sum = fmt.Sprintf("%x", sha256.Sum256(b))
        }
        entries = append(entries, Entry{Path: dst, SHA256: sum})
    }
    m := Manifest{Version: version, InstalledAt: "", Files: entries, Stacks: stacks}
    buf, _ := json.MarshalIndent(m, "", "  ")
    return os.WriteFile(filePath, buf, 0o644)
}

func Open() (Manifest, error) {
    var m Manifest
    f, err := os.Open(filePath)
    if err != nil {
        return m, err
    }
    defer f.Close()
    b, _ := io.ReadAll(f)
    err = json.Unmarshal(b, &m)
    return m, err
}
