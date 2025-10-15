package zipbuild

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	epoch = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
)

// Build produces a deterministic zip archive of packDir.
func Build(packDir string) ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	type entry struct {
		rel  string
		full string
		info fs.FileInfo
	}

	var dirs []entry
	err := filepath.WalkDir(packDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(packDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		dirs = append(dirs, entry{rel: rel, full: path, info: info})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].rel < dirs[j].rel })

	for _, e := range dirs {
		header, err := zip.FileInfoHeader(e.info)
		if err != nil {
			return nil, err
		}
		header.Name = filepath.ToSlash(e.rel)
		header.Method = zip.Deflate
		header.Modified = epoch
		if e.info.IsDir() {
			if !strings.HasSuffix(header.Name, "/") {
				header.Name += "/"
			}
			header.SetMode(fs.ModeDir | 0o755)
			if _, err := zw.CreateHeader(header); err != nil {
				return nil, err
			}
			continue
		}
		header.SetMode(0o644)
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(e.full)
		if err != nil {
			return nil, err
		}
		if _, err := writer.Write(data); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
