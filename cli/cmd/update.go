package cmd

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/manifest"
	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
	"github.com/spf13/cobra"
)

var updateTo string
var updateDry bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the toolkit (only overwrite files unchanged since install)",
	RunE: func(cmd *cobra.Command, args []string) error {
		abortIf(!manifest.Exists(), "No manifest found. Run `codo init` first.")
		m, err := manifest.Open()
		if err != nil {
			return err
		}

		// Pack resolution for update (similar to init)
		var rootFS fs.FS
		var packSource string

		if _, err := os.Stat("pack"); err == nil {
			// Use local pack for development
			rootFS = os.DirFS("pack")
			packSource = "local"
			fmt.Println("Using local pack directory for update")
		} else {
			// Try to download pack from GitHub
			versionToFetch := updateTo
			if versionToFetch == "" {
				versionToFetch = "latest"
			}
			fmt.Printf("Downloading pack version: %s...\n", versionToFetch)

			packPath, err := pack.Resolve(versionToFetch)
			if err == nil {
				rootFS = os.DirFS(packPath)
				packSource = versionToFetch
				fmt.Printf("Downloaded pack from GitHub releases\n")
			} else {
				// Fall back to embedded base
				fmt.Printf("Download failed (%v), using embedded base pack\n", err)
				rootFS, err = pack.GetEmbeddedBaseFS()
				if err != nil {
					return fmt.Errorf("failed to load embedded pack: %w", err)
				}
				packSource = "embedded-base"
			}
		}

		files, err := pack.FilesFromDotclaudeFS(rootFS, m.Stacks)
		if err != nil {
			return err
		}

		// Build map of new contents
		newMap := map[string][]byte{}
		for _, f := range files {
			b, err := f.Read()
			if err == nil {
				newMap[f.RelPath] = b
			}
		}

		// Ensure report dir exists
		_ = os.MkdirAll(filepath.Join(".claude", ".codo-report"), 0o755)

		// Track which files exist in the old manifest for removal detection
		oldSet := map[string]manifest.Entry{}
		for _, ent := range m.Files {
			oldSet[ent.Path] = ent
		}

		// Process each file from the old manifest
		for _, ent := range m.Files {
			dst := ent.Path
			nb, ok := newMap[dst]
			if !ok {
				// File removed upstream - handle safely
				cur, err := os.ReadFile(dst)
				if err != nil {
					// File already gone, nothing to do
					continue
				}
				curHash := fmt.Sprintf("%x", sha256.Sum256(cur))
				if curHash == ent.SHA256 {
					// File is clean (unmodified) - safe to remove
					fmt.Println("- " + dst)
					if !updateDry {
						if err := os.Remove(dst); err != nil {
							return err
						}
					}
				} else {
					// File has local modifications - keep it and notify user
					note := dst + ".codo.removed.suggested"
					fmt.Println("! modified & removed upstream → " + note)
					if !updateDry {
						msg := []byte("Upstream removed this file, but you have local changes.\nConsider removing it manually if no longer needed.\n")
						if err := os.WriteFile(note, msg, 0o644); err != nil {
							return err
						}
					}
				}
				continue
			}

			// File exists in new pack - check if it needs updating
			cur, err := os.ReadFile(dst)
			if err != nil {
				// Missing → treat as clean overwrite
				fmt.Println("+ " + dst)
				if !updateDry {
					if err := os.WriteFile(dst, nb, 0o644); err != nil {
						return err
					}
				}
				continue
			}
			curHash := fmt.Sprintf("%x", sha256.Sum256(cur))
			if curHash == ent.SHA256 {
				// clean → overwrite
				fmt.Println("~ " + dst)
				if !updateDry {
					if err := os.WriteFile(dst, nb, 0o644); err != nil {
						return err
					}
				}
			} else {
				// diverged → write .codo.new
				out := dst + ".codo.new"
				fmt.Println("! conflict → " + out)
				if !updateDry {
					if err := os.WriteFile(out, nb, 0o644); err != nil {
						return err
					}
				}
			}
		}

		// Add any new files that weren't in the old manifest
		for path, content := range newMap {
			if _, exists := oldSet[path]; !exists {
				fmt.Println("+ " + path)
				if !updateDry {
					// Ensure directory exists
					if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
						return err
					}
					if err := os.WriteFile(path, content, 0o644); err != nil {
						return err
					}
				}
			}
		}
		if !updateDry {
			newVersion := packSource
			if err := manifest.WriteWithStacks(files, newVersion, m.Stacks); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateTo, "to", "", "Version/tag to update to (e.g. v1.2.0)")
	updateCmd.Flags().BoolVar(&updateDry, "dry-run", false, "Preview only")
}
