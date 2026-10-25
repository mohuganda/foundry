// Package scaffold validates file plans before writing and rolls back handled failures.
package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	Data []byte
	Mode os.FileMode
}

func Apply(root string, files map[string]File, overwrite, dryRun bool) ([]string, error) {
	return apply(root, files, overwrite, dryRun, func(f *os.File, data []byte) error {
		_, err := f.Write(data)
		return err
	})
}

func apply(root string, files map[string]File, overwrite, dryRun bool, write func(*os.File, []byte) error) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(files))
	for rel := range files {
		if !filepath.IsLocal(rel) || filepath.Clean(rel) != rel || rel == "." {
			return nil, fmt.Errorf("invalid scaffold path %q", rel)
		}
		keys = append(keys, rel)
	}
	sort.Strings(keys)
	for _, rel := range keys {
		if err := validate(root, rel, overwrite); err != nil {
			return nil, err
		}
		for parent := filepath.Dir(rel); parent != "."; parent = filepath.Dir(parent) {
			if _, ok := files[parent]; ok {
				return nil, fmt.Errorf("file is also a planned directory: %s", parent)
			}
		}
	}
	if dryRun {
		return keys, nil
	}
	type previous struct {
		path    string
		data    []byte
		mode    os.FileMode
		existed bool
	}
	var changed []previous
	var dirs []string
	rollback := func(cause error) ([]string, error) {
		for i := len(changed) - 1; i >= 0; i-- {
			p := changed[i]
			var err error
			if p.existed {
				err = os.WriteFile(p.path, p.data, p.mode)
				if err == nil {
					err = os.Chmod(p.path, p.mode)
				}
			} else {
				err = os.Remove(p.path)
			}
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				cause = errors.Join(cause, fmt.Errorf("rollback %s: %w", p.path, err))
			}
		}
		for i := len(dirs) - 1; i >= 0; i-- {
			if err := os.Remove(dirs[i]); err != nil {
				cause = errors.Join(cause, fmt.Errorf("rollback directory %s: %w", dirs[i], err))
			}
		}
		return nil, cause
	}
	var mkdir func(string) error
	mkdir = func(dir string) error {
		st, err := os.Stat(dir)
		if err == nil {
			if !st.IsDir() {
				return fmt.Errorf("not a directory: %s", dir)
			}
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := mkdir(filepath.Dir(dir)); err != nil {
			return err
		}
		if err := os.Mkdir(dir, 0o755); err != nil {
			return err
		}
		dirs = append(dirs, dir)
		return nil
	}
	for _, rel := range keys {
		if err := validate(root, rel, overwrite); err != nil {
			return rollback(err)
		}
		path := filepath.Join(root, rel)
		if err := mkdir(filepath.Dir(path)); err != nil {
			return rollback(err)
		}
		old := previous{path: path}
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		if st, err := os.Lstat(path); err == nil && overwrite {
			old.existed, old.mode = true, st.Mode().Perm()
			old.data, err = os.ReadFile(path)
			if err != nil {
				return rollback(err)
			}
			flags = os.O_WRONLY | os.O_TRUNC
		}
		file, err := os.OpenFile(path, flags, files[rel].Mode)
		if err != nil {
			return rollback(err)
		}
		changed = append(changed, old)
		err = write(file, files[rel].Data)
		closeErr := file.Close()
		if err == nil && closeErr == nil {
			err = os.Chmod(path, files[rel].Mode)
		}
		if err = errors.Join(err, closeErr); err != nil {
			return rollback(err)
		}
	}
	return keys, nil
}

func validate(root, rel string, overwrite bool) error {
	parts := append([]string{"."}, strings.Split(rel, string(filepath.Separator))...)
	path := root
	for i, part := range parts {
		if i > 0 {
			path = filepath.Join(path, part)
		}
		st, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing scaffold through symlink %s", path)
		}
		if i < len(parts)-1 {
			if !st.IsDir() {
				return fmt.Errorf("not a directory: %s", path)
			}
		} else if !st.Mode().IsRegular() || !overwrite {
			return fmt.Errorf("refusing to overwrite existing path %s", path)
		}
	}
	return nil
}
