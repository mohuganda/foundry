// Package workspace identifies project roots without modifying the filesystem.
package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jabahum/foundry/internal/config"
	"github.com/jabahum/foundry/internal/registry"
)

var ErrNotFound = errors.New("no supported workspace found")

type Project struct {
	Root, Kind, Template string
}

// At identifies exactly root; Find also searches its parents.
func At(root string) (Project, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Project{}, err
	}
	if _, err := os.Stat(config.Path(root)); err == nil {
		c, err := config.Read(root)
		if err != nil {
			return Project{}, err
		}
		t, err := registry.Find(c.Template)
		if err != nil {
			return Project{}, err
		}
		if c.Schema > 1 || c.Schema < 0 || (c.Kind != "" && c.Kind != t.Kind) {
			return Project{}, fmt.Errorf("unsupported or inconsistent .foundry.yaml in %s", root)
		}
		return Project{root, t.Kind, t.ID}, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return Project{}, err
	}
	kind := ""
	switch {
	case exists(root, "assets.go") && exists(root, "cmd/foundry/main.go"):
		kind = "foundry"
	case exists(root, "pubspec.yaml") && directory(root, "apps"):
		kind = "flutter"
	case exists(root, "go.mod"):
		kind = "go"
	case exists(root, "Cargo.toml"):
		kind = "rust"
	case exists(root, "package.json"):
		kind = "javascript"
		var manifest struct{ Dependencies map[string]json.RawMessage }
		data, err := os.ReadFile(filepath.Join(root, "package.json"))
		if err != nil {
			return Project{}, err
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return Project{}, fmt.Errorf("invalid package.json: %w", err)
		}
		if _, ok := manifest.Dependencies["react"]; ok {
			kind = "react"
		}
	default:
		return Project{}, fmt.Errorf("%w at %s", ErrNotFound, root)
	}
	return Project{Root: root, Kind: kind}, nil
}

func Find(start string) (Project, error) {
	root, err := filepath.Abs(start)
	if err != nil {
		return Project{}, err
	}
	for {
		p, err := At(root)
		if err == nil || !errors.Is(err, ErrNotFound) {
			return p, err
		}
		parent := filepath.Dir(root)
		if parent == root {
			return Project{}, ErrNotFound
		}
		root = parent
	}
}

func exists(root, name string) bool {
	st, err := os.Stat(filepath.Join(root, filepath.FromSlash(name)))
	return err == nil && st.Mode().IsRegular()
}

func directory(root, name string) bool {
	st, err := os.Stat(filepath.Join(root, name))
	return err == nil && st.IsDir()
}
