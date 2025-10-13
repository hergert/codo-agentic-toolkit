package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/manifest"
	"github.com/hergert/codo-agentic-toolkit/cli/internal/statepath"
	"github.com/spf13/cobra"
)

var removeDry bool

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the toolkit (backup first)",
	RunE: func(cmd *cobra.Command, args []string) error {
		abortIf(!manifest.Exists(), "No manifest found. Nothing to remove.")
		m, err := manifest.Open()
		if err != nil {
			return err
		}
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		ts := time.Now().UTC().Format("20060102-150405")
		var backup string
		if !removeDry {
			dir, err := statepath.BackupDir(root, ts)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			backup = dir
		}
		for _, ent := range m.Files {
			if ent.Unmanaged {
				fmt.Println("~ skip unmanaged " + ent.Path)
				if !removeDry {
					_ = os.Remove(ent.Path + ".codo.new")
				}
				continue
			}
			if _, err := os.Stat(ent.Path); err == nil {
				fmt.Println("- " + ent.Path)
				if !removeDry {
					dest := filepath.Join(backup, ent.Path)
					if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
						return err
					}
					if err := moveWithCopyFallback(ent.Path, dest); err != nil {
						return err
					}
					_ = os.Remove(ent.Path + ".codo.new")
				}
			}
		}
		if !removeDry {
			manifest.Remove()
			fmt.Println("Backup at", backup)
		} else {
			fmt.Println("(dry-run) Removal would back up files outside the repo")
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVar(&removeDry, "dry-run", false, "Preview only; do not write files")
}

func moveWithCopyFallback(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		sin, serr := os.Open(src)
		if serr != nil {
			return serr
		}
		defer sin.Close()
		out, oerr := os.Create(dst)
		if oerr != nil {
			return oerr
		}
		if _, cerr := io.Copy(out, sin); cerr != nil {
			out.Close()
			return cerr
		}
		if cerr := out.Close(); cerr != nil {
			return cerr
		}
		if rerr := os.Remove(src); rerr != nil {
			return rerr
		}
		return nil
	}
	return nil
}
