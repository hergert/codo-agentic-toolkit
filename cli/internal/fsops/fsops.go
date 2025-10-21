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

func CopySafe(f pack.File, projectRoot string, dry bool) (bool, error) {
	dst := filepath.Join(projectRoot, f.RelPath)
	if !dry {
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return false, err
		}
	}
	srcBytes, err := f.Read()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(dst); err == nil {
		tmp := dst + ".codo.new"
		curHash, err := sha256File(dst)
		if err != nil {
			return false, fmt.Errorf("hash current file %s: %w", dst, err)
		}
		newHash := fmt.Sprintf("%x", sha256.Sum256(srcBytes))
		if curHash == newHash {
			fmt.Println("= " + f.RelPath)
			return true, nil
		}
		fmt.Println("! conflict → " + tmp)
		if dry {
			return false, nil
		}
		if err := os.WriteFile(tmp, srcBytes, 0o644); err != nil {
			return false, err
		}
		return false, nil
	}

	fmt.Println("+ " + f.RelPath)
	if dry {
		return true, nil
	}
	if err := os.WriteFile(dst, srcBytes, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func ChmodHooks() error {
	files := []string{
		filepath.Join(".claude", "hooks", "pre_tool_use.py"),
		filepath.Join(".claude", "hooks", "post_tool_use.py"),
		filepath.Join(".claude", "hooks", "user_prompt_submit.py"),
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
