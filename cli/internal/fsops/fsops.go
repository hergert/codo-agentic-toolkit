package fsops

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hergert/codo-agentic-toolkit/cli/internal/pack"
)

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func CopySafe(f pack.File, projectRoot string, dry bool) error {
	dst := filepath.Join(projectRoot, f.RelPath)
	if !dry {
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
	}
	srcBytes, err := f.Read()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dst); err == nil {
		tmp := dst + ".codo.new"
		curHash, err := sha256File(dst)
		if err != nil {
			return fmt.Errorf("hash current file %s: %w", dst, err)
		}
		newHash := fmt.Sprintf("%x", sha256.Sum256(srcBytes))
		if curHash == newHash {
			fmt.Println("= " + f.RelPath)
			return nil
		}
		fmt.Println("! conflict → " + tmp)
		if dry {
			return nil
		}
		appendReport("conflicts.txt", tmp)
		return os.WriteFile(tmp, srcBytes, 0o644)
	}

	fmt.Println("+ " + f.RelPath)
	if dry {
		return nil
	}
	return os.WriteFile(dst, srcBytes, 0o644)
}

func ChmodHooks() error {
	files := []string{
		filepath.Join(".claude", "hooks", "pre_tool_use.py"),
		filepath.Join(".claude", "hooks", "post_tool_use.py"),
		filepath.Join(".claude", "hooks", "user_prompt_submit.py"),
		filepath.Join(".claude", "hooks", "experience_update.py"),
	}
	for _, p := range files {
		if _, err := os.Stat(p); err == nil {
			if err := os.Chmod(p, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}
