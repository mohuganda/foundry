package templatehub

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jabahum/foundry/internal/cache"
	"github.com/jabahum/foundry/internal/registry"
)

type Manifest struct {
	Schema    int               `json:"schema"`
	UpdatedAt string            `json:"updated_at,omitempty"`
	Templates []TemplateRelease `json:"templates"`
}

type TemplateRelease struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Kind        string   `json:"kind"`
	Aliases     []string `json:"aliases,omitempty"`
	URL         string   `json:"url"`
	SHA256      string   `json:"sha256"`
}

type Resolved struct {
	Template    registry.Template
	ArchivePath string
	Source      string
}

func registryPath() (string, error) {
	d, err := cache.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "registry.json"), nil
}

func LoadCachedManifest() (*Manifest, error) {
	p, err := registryPath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parse cached registry: %w", err)
	}
	return &m, nil
}

func Update(registryURL string) (*Manifest, error) {
	registryURL = strings.TrimSpace(registryURL)
	if registryURL == "" {
		registryURL = strings.TrimSpace(os.Getenv("FOUNDRY_REGISTRY_URL"))
	}
	if registryURL == "" {
		return nil, fmt.Errorf("registry URL required: use --registry or FOUNDRY_REGISTRY_URL")
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, fmt.Errorf("fetch registry: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch registry: HTTP %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("parse registry: %w", err)
	}
	if m.Schema == 0 {
		m.Schema = 1
	}
	for _, t := range m.Templates {
		if t.ID == "" || t.Version == "" || t.URL == "" || t.SHA256 == "" {
			return nil, fmt.Errorf("registry contains incomplete template release")
		}
	}
	p, err := registryPath()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(p, body, 0o644); err != nil {
		return nil, err
	}
	return &m, nil
}

func ListCached() ([]TemplateRelease, error) {
	m, err := LoadCachedManifest()
	if err != nil || m == nil {
		return nil, err
	}
	out := append([]TemplateRelease(nil), m.Templates...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].ID == out[j].ID {
			return out[i].Version > out[j].Version
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func Resolve(name, version string, offline bool) (Resolved, error) {
	builtin, builtinErr := registry.Find(name)
	if version == "" && builtinErr == nil {
		return Resolved{Template: builtin, Source: "embedded"}, nil
	}

	m, err := LoadCachedManifest()
	if err != nil {
		return Resolved{}, err
	}
	if m == nil {
		if builtinErr == nil && (version == "" || version == builtin.Version) {
			return Resolved{Template: builtin, Source: "embedded"}, nil
		}
		return Resolved{}, fmt.Errorf("template %q version %q is not cached; run foundry template update first", name, version)
	}

	var candidates []TemplateRelease
	q := strings.ToLower(strings.TrimSpace(name))
	for _, r := range m.Templates {
		match := q == strings.ToLower(r.ID) || q == strings.ToLower(r.Name)
		if !match {
			for _, a := range r.Aliases {
				if q == strings.ToLower(a) {
					match = true
					break
				}
			}
		}
		if match && (version == "" || r.Version == version) {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) == 0 {
		if builtinErr == nil && (version == "" || version == builtin.Version) {
			return Resolved{Template: builtin, Source: "embedded"}, nil
		}
		return Resolved{}, fmt.Errorf("template %q version %q not found in cached registry", name, version)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Version > candidates[j].Version })
	rel := candidates[0]
	path, err := archivePath(rel)
	if err != nil {
		return Resolved{}, err
	}
	if _, err := os.Stat(path); err == nil {
		if err := verify(path, rel.SHA256); err == nil {
			return Resolved{Template: toTemplate(rel), ArchivePath: path, Source: "cache"}, nil
		}
		_ = os.Remove(path)
	}
	if offline {
		return Resolved{}, fmt.Errorf("template %s@%s is not cached and --offline was set", rel.ID, rel.Version)
	}
	if err := download(rel, path); err != nil {
		return Resolved{}, err
	}
	return Resolved{Template: toTemplate(rel), ArchivePath: path, Source: "remote"}, nil
}

func toTemplate(r TemplateRelease) registry.Template {
	return registry.Template{ID: r.ID, Name: r.Name, Version: r.Version, Description: r.Description, Aliases: r.Aliases, Kind: r.Kind}
}

func archivePath(r TemplateRelease) (string, error) {
	d, err := cache.Dir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(d, "templates", r.ID, r.Version)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "template.zip"), nil
}

func download(r TemplateRelease, target string) error {
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(r.URL)
	if err != nil {
		return fmt.Errorf("download %s@%s: %w", r.ID, r.Version, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download %s@%s: HTTP %s", r.ID, r.Version, resp.Status)
	}
	tmp := target + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, io.LimitReader(resp.Body, 512<<20))
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}
	if err := verify(tmp, r.SHA256); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, target)
}

func verify(path, want string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	want = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(want, "sha256:")))
	if got != want {
		return fmt.Errorf("checksum mismatch for %s: got %s want %s", path, got, want)
	}
	return nil
}
