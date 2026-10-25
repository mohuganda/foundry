package generator

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jabahum/foundry/internal/config"
	"github.com/jabahum/foundry/internal/identity"
	"github.com/jabahum/foundry/internal/registry"
	"github.com/jabahum/foundry/internal/scaffold"
	"github.com/jabahum/foundry/internal/workspace"
)

type Options struct {
	Name, Destination, Module, Org, Package string
	Profile, IdentityFile                   string
	Force, NoHooks, DryRun                  bool
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
	if st, err := os.Lstat(dest); err == nil {
		if !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination must be a directory, not a file or symlink")
		}
		entries, err := os.ReadDir(dest)
		if err != nil {
			return err
		}
		if len(entries) > 0 && !o.Force {
			return fmt.Errorf("destination %s is not empty; use --force", dest)
		}
		if _, err := os.Stat(config.Path(dest)); err == nil {
			p, err := workspace.At(dest)
			if err != nil {
				return err
			}
			if p.Kind != t.Kind {
				return fmt.Errorf("destination is a %s workspace, not %s", p.Kind, t.Kind)
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	_, sharedErr := fs.Stat(g.FS, "templates/shared/identity")
	withIdentity := sharedErr == nil
	if !withIdentity && (o.Profile != "" || o.IdentityFile != "") {
		return fmt.Errorf("identity support is not bundled in this generator")
	}
	var identityConfig identity.Config
	var identityAssets map[string][]byte
	if withIdentity {
		identityConfig, identityAssets, err = identity.Resolve(o.Profile, o.IdentityFile, o.Name)
		if err != nil {
			return err
		}
	}
	stage, err := os.MkdirTemp("", "foundry-stage-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	output := dest
	dest = stage
	if o.DryRun {
		o.NoHooks = true
	}
	if t.Kind == "flutter" && !o.NoHooks {
		if err := prepareFlutter(dest, o); err != nil {
			fmt.Fprintf(os.Stderr, "warning: flutter platform scaffold: %v\n", err)
		}
	}
	if err := copyTree(g.FS, t.EmbeddedDir, dest); err != nil {
		return err
	}

	vars := map[string]string{
		"go-foundations":             defaultString(o.Module, sanitize(o.Name)),
		"example.com/go-foundations": defaultString(o.Module, sanitize(o.Name)),
		"__RUST_CRATE__":             rustCrateName(o.Name),
		"__PROJECT_NAME__":           sanitize(o.Name),
		"__ORG__":                    defaultString(o.Org, "com.example"),
		"__PACKAGE__":                defaultString(o.Package, sanitize(o.Name)),
	}
	if err := replaceText(dest, vars); err != nil {
		return err
	}
	if withIdentity {
		readme, err := fs.ReadFile(g.FS, "templates/shared/identity/README.md")
		if err != nil {
			return err
		}
		files, err := identity.Render(g.FS, identityConfig, identityAssets, t.Kind)
		if err != nil {
			return err
		}
		files[filepath.FromSlash("identity/README.md")] = scaffold.File{Data: readme, Mode: 0o644}
		if _, err := scaffold.Apply(stage, files, true, false); err != nil {
			return err
		}
		if t.Kind == "javascript" || t.Kind == "react" {
			for _, dir := range []string{"common", t.Kind} {
				path := "templates/shared/identity/web/" + dir
				if _, err := fs.Stat(g.FS, path); err == nil {
					if err := copyTree(g.FS, path, filepath.Join(stage, "src", "shared", "identity")); err != nil {
						return err
					}
				}
			}
		}
	}
	if !o.NoHooks {
		if err := postHooks(t, dest, o); err != nil {
			fmt.Fprintf(os.Stderr, "warning: post-generation hook: %v\n", err)
		}
	}
	if err := config.Write(dest, config.Config{ProjectName: sanitize(o.Name), Template: t.ID, TemplateVersion: t.Version, Kind: t.Kind, Module: o.Module, Org: o.Org, Package: o.Package}); err != nil {
		return fmt.Errorf("write .foundry.yaml: %w", err)
	}
	plan := map[string]scaffold.File{}
	if err := filepath.WalkDir(stage, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".dart_tool", ".fvm", ".mason", "node_modules", "target", "build":
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("generated output contains a symlink: %s", path)
		}
		rel, err := filepath.Rel(stage, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		st, err := d.Info()
		if err != nil {
			return err
		}
		plan[rel] = scaffold.File{Data: data, Mode: st.Mode().Perm()}
		return nil
	}); err != nil {
		return err
	}
	paths, err := scaffold.Apply(output, plan, o.Force, o.DryRun)
	if err != nil {
		return err
	}
	if o.DryRun {
		fmt.Printf("Would create %s using %s v%s (hooks skipped)\n", output, t.Name, t.Version)
		for _, path := range paths {
			fmt.Println("  +", path)
		}
		return nil
	}
	fmt.Printf("\nCreated %s using %s v%s\n", output, t.Name, t.Version)
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
		if !d.IsDir() {
			rel = strings.TrimSuffix(rel, ".tmpl")
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
		// embed.FS normalizes permissions, so restore runnable shell scripts.
		if filepath.Ext(rel) == ".sh" {
			mode = 0o755
		}
		if info, e := d.Info(); e == nil && info.Mode()&0o111 != 0 {
			mode = 0o755
		}
		return os.WriteFile(target, data, mode)
	})
}

func replaceText(root string, vars map[string]string) error {
	exts := map[string]bool{".rs": true, ".lock": true, ".html": true, ".js": true, ".jsx": true, ".mjs": true, ".css": true, ".go": true, ".mod": true, ".sum": true, ".proto": true, ".yaml": true, ".yml": true, ".json": true, ".md": true, ".sql": true, ".env": true, ".example": true, ".dart": true, ".xml": true, ".gradle": true, ".properties": true, ".toml": true, ".makefile": true, ".tmpl": true}
	ordered := [][2]string{
		{"example.com/go-foundations", vars["example.com/go-foundations"]},
		{"go-foundations", vars["go-foundations"]},
		{"__RUST_CRATE__", vars["__RUST_CRATE__"]},
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

// rustCrateName produces an identifier usable in package, library, and binary names.
func rustCrateName(name string) string {
	name = sanitize(name)
	reserved := " as async await break const continue crate dyn else enum extern false fn for if impl in let loop match mod move mut pub ref return self static struct super trait true type unsafe use where while abstract become box do final gen macro override priv typeof unsized virtual yield try union std core alloc test proc_macro build deps examples incremental con prn aux nul com1 com2 com3 com4 com5 com6 com7 com8 com9 lpt1 lpt2 lpt3 lpt4 lpt5 lpt6 lpt7 lpt8 lpt9 "
	if name == "_" || (name[0] >= '0' && name[0] <= '9') || strings.Contains(reserved, " "+name+" ") {
		return "app_" + name
	}
	return name
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
