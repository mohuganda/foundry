package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Dir() (string, error) {
	d := strings.TrimSpace(os.Getenv("FOUNDRY_CACHE_DIR"))
	if d == "" {
		base, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		d = filepath.Join(base, "foundry")
	} else {
		var err error
		d, err = filepath.Abs(d)
		if err != nil {
			return "", err
		}
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	return d, nil
}

func Clean() error {
	d, err := Dir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(d)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(d, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

func Summary() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (%d entries)", d, len(entries)), nil
}
