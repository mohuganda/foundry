package generator

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jabahum/foundry/internal/projectconfig"
	"github.com/jabahum/foundry/internal/registry"
)

type Options struct {
	Name, Destination, Module, Org, Package string
	ArchivePath                             string
	Force, NoHooks                          bool
}

type Generator struct{ FS fs.FS }

func (g Generator) Generate(t registry.Template, o Options) error {
	if strings.TrimSpace(o.Name) == "" {
		return fmt.Errorf("project name is required")
	}
	if o.Destination == "" {
		o.Destination = o.Name
	}
	dest, err := filepath.Abs(o.Destination)
	if err != nil {
		return err
	}
	if st, err := os.Stat(dest); err == nil && st.IsDir() {
		entries, _ := os.ReadDir(dest)
		if len(entries) > 0 && !o.Force {
			return fmt.Errorf("destination %s is not empty; use --force", dest)
		}
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	if t.Kind == "flutter" && !o.NoHooks {
		if err := prepareFlutter(dest, o); err != nil {
			fmt.Fprintf(os.Stderr, "warning: flutter platform scaffold: %v\n", err)
		}
	}
	if o.ArchivePath != "" {
		if err := extractZipFile(o.ArchivePath, dest); err != nil {
			return err
		}
	} else if t.EmbeddedArchive != "" {
		if err := extractEmbeddedZip(g.FS, t.EmbeddedArchive, dest); err != nil {
			return err
		}
	} else {
		if err := copyTree(g.FS, t.EmbeddedDir, dest); err != nil {
			return err
		}
	}

	vars := map[string]string{
		"go-foundations":             defaultString(o.Module, sanitize(o.Name)),
		"example.com/go-foundations": defaultString(o.Module, sanitize(o.Name)),
		"__PROJECT_NAME__":           sanitize(o.Name),
		"__ORG__":                    defaultString(o.Org, "com.example"),
		"__PACKAGE__":                defaultString(o.Package, sanitize(o.Name)),
	}
	if err := replaceText(dest, vars); err != nil {
		return err
	}
	if !o.NoHooks {
		if err := postHooks(t, dest, o); err != nil {
			fmt.Fprintf(os.Stderr, "warning: post-generation hook: %v\n", err)
		}
	}
	if err := projectconfig.Write(dest, projectconfig.Config{ProjectName: sanitize(o.Name), Template: t.ID, TemplateVersion: t.Version, Kind: t.Kind, Module: o.Module, Org: o.Org, Package: o.Package}); err != nil {
		return fmt.Errorf("write .foundry.yaml: %w", err)
	}
	fmt.Printf("\nCreated %s using %s v%s\n", dest, t.Name, t.Version)
	return nil
}

func extractZipFile(archivePath, dst string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(f, st.Size())
	if err != nil {
		return fmt.Errorf("open template archive: %w", err)
	}
	for _, zf := range zr.File {
		clean := filepath.Clean(zf.Name)
		if clean == "." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.IsAbs(clean) {
			return fmt.Errorf("unsafe archive path %q", zf.Name)
		}
		target := filepath.Join(dst, clean)
		if zf.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		r, err := zf.Open()
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(r)
		closeErr := r.Close()
		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}
		mode := zf.Mode()
		if mode == 0 {
			mode = 0o644
		}
		if err := os.WriteFile(target, data, mode); err != nil {
			return err
		}
	}
	return nil
}

func extractEmbeddedZip(src fs.FS, archivePath, dst string) error {
	data, err := fs.ReadFile(src, archivePath)
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("open embedded template archive: %w", err)
	}
	for _, f := range zr.File {
		clean := filepath.Clean(f.Name)
		if clean == "." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.IsAbs(clean) {
			return fmt.Errorf("unsafe archive path %q", f.Name)
		}
		target := filepath.Join(dst, clean)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(r)
		closeErr := r.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		mode := f.Mode()
		if mode == 0 {
			mode = 0o644
		}
		if err := os.WriteFile(target, content, mode); err != nil {
			return err
		}
	}
	return nil
}

func copyTree(src fs.FS, root, dst string) error {
	return fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}
		mode := fs.FileMode(0o644)
		if info, e := d.Info(); e == nil && info.Mode()&0o111 != 0 {
			mode = 0o755
		}
		return os.WriteFile(target, data, mode)
	})
}

func replaceText(root string, vars map[string]string) error {
	exts := map[string]bool{".go": true, ".mod": true, ".sum": true, ".proto": true, ".yaml": true, ".yml": true, ".json": true, ".md": true, ".sql": true, ".env": true, ".example": true, ".dart": true, ".xml": true, ".gradle": true, ".properties": true, ".toml": true, ".makefile": true, ".tmpl": true}
	ordered := [][2]string{
		{"example.com/go-foundations", vars["example.com/go-foundations"]},
		{"go-foundations", vars["go-foundations"]},
		{"__PROJECT_NAME__", vars["__PROJECT_NAME__"]},
		{"__ORG__", vars["__ORG__"]},
		{"__PACKAGE__", vars["__PACKAGE__"]},
	}
	pairs := make([]string, 0, len(ordered)*2)
	for _, pair := range ordered {
		if pair[1] != "" {
			pairs = append(pairs, pair[0], pair[1])
		}
	}
	replacer := strings.NewReplacer(pairs...)
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := strings.ToLower(filepath.Base(path))
		ext := strings.ToLower(filepath.Ext(path))
		if !exts[ext] && base != "makefile" && !strings.HasPrefix(base, ".env") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := replacer.Replace(string(data))
		return os.WriteFile(path, []byte(s), 0o644)
	})
}

func postHooks(t registry.Template, dest string, o Options) error {
	if t.Kind == "go" {
		if _, err := exec.LookPath("go"); err == nil {
			cmd := exec.Command("go", "fmt", "./...")
			cmd.Dir = dest
			_ = cmd.Run()
		}
		return nil
	}
	if t.Kind == "flutter" {
		appDir := filepath.Join(dest, "apps", sanitize(o.Name))
		sourceApp := filepath.Join(dest, "apps", "app")
		if _, err := os.Stat(sourceApp); err == nil && sourceApp != appDir {
			if _, err := os.Stat(appDir); os.IsNotExist(err) {
				if err := os.Rename(sourceApp, appDir); err != nil {
					return err
				}
			}
		}
		if _, err := exec.LookPath("melos"); err == nil {
			cmd := exec.Command("melos", "bootstrap")
			cmd.Dir = dest
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}
	}
	return nil
}

func prepareFlutter(dest string, o Options) error {
	if _, err := exec.LookPath("flutter"); err != nil {
		return fmt.Errorf("flutter not found; platform folders can be generated later with flutter create")
	}
	appDir := filepath.Join(dest, "apps", "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return err
	}
	cmd := exec.Command("flutter", "create", "--project-name", sanitize(o.Name), "--org", defaultString(o.Org, "com.example"), ".")
	cmd.Dir = appDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else if r == '-' || r == ' ' {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "app"
	}
	return b.String()
}
func defaultString(v, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}

func Confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	s, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "y" || s == "yes"
}
