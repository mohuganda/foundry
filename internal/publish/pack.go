package publish

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Pack(dir, out string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if out == "" {
		out = filepath.Base(dir) + ".zip"
	}
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	zw := zip.NewWriter(f)
	walkErr := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, e := filepath.Rel(dir, path)
		if e != nil {
			return e
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, ".git/") || rel == ".DS_Store" {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		h, e := zip.FileInfoHeader(info)
		if e != nil {
			return e
		}
		h.Name = rel
		h.Method = zip.Deflate
		w, e := zw.CreateHeader(h)
		if e != nil {
			return e
		}
		in, e := os.Open(path)
		if e != nil {
			return e
		}
		_, copyErr := io.Copy(w, in)
		closeErr := in.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
	closeZip := zw.Close()
	closeFile := f.Close()
	if walkErr != nil {
		return "", walkErr
	}
	if closeZip != nil {
		return "", closeZip
	}
	if closeFile != nil {
		return "", closeFile
	}
	data, err := os.ReadFile(out)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func RegistryEntry(id, version, kind, url, sha string) string {
	return fmt.Sprintf(`{"id":%q,"version":%q,"kind":%q,"url":%q,"sha256":%q}`, id, version, kind, url, sha)
}
