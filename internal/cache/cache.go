package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func Dir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "foundry")
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
