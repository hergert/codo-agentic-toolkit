package fsops

import (
	"fmt"
	"os"
	"path/filepath"
)

func appendReport(relPath, line string) {
	report := filepath.Join(".claude", ".codo-report", relPath)
	if err := os.MkdirAll(filepath.Dir(report), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(report, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintln(f, line)
}

// AppendReportHook exposes conflict reporting for other packages.
func AppendReportHook(path string) {
	appendReport("conflicts.txt", path)
}

// Sha256File exposes the shared hashing helper.
func Sha256File(path string) (string, error) {
	return sha256File(path)
}
