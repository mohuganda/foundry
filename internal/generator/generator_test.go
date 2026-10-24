package generator

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jabahum/foundry"
	"github.com/jabahum/foundry/internal/config"
	"github.com/jabahum/foundry/internal/registry"
)

func TestBundledTemplatesGenerate(t *testing.T) {
	for _, name := range []string{"go", "flutter", "vanilla-js", "react"} {
		t.Run(name, func(t *testing.T) {
			tmpl, err := registry.Find(name)
			if err != nil {
				t.Fatal(err)
			}
			dest := t.TempDir()
			g := Generator{FS: foundry.Templates}
			if err := g.Generate(tmpl, Options{Name: "demo", Destination: dest, NoHooks: true}); err != nil {
				t.Fatal(err)
			}
			if err := filepath.WalkDir(dest, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(path, ".tmpl") {
					t.Errorf("unrestored template filename: %s", path)
				}
				if d.IsDir() && (d.Name() == ".dart_tool" || d.Name() == ".fvm" || d.Name() == ".mason" || d.Name() == "node_modules" || d.Name() == "dist") {
					t.Errorf("local tool cache delivered: %s", path)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			if name == "go" {
				if _, err := os.Stat(filepath.Join(dest, "internal/application/auth/service_test.go")); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(filepath.Join(dest, "gen/proto/auth/v1/auth.pb.go")); err != nil {
					t.Fatal(err)
				}
			}
			if runtime.GOOS != "windows" && (name == "go" || name == "flutter") {
				script := "scripts/setup.sh"
				if name == "go" {
					script = "scripts/generate-dev-keys.sh"
				}
				info, err := os.Stat(filepath.Join(dest, script))
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode()&0o111 == 0 {
					t.Errorf("script is not executable: %s", script)
				}
			}
			file, expected := "pubspec.yaml", "name: flutter_foundation_workspace"
			if name == "go" {
				file, expected = "go.mod", "module demo\n"
			} else if name == "vanilla-js" || name == "react" {
				file, expected = "package.json", `"name": "demo"`
			}
			data, err := os.ReadFile(filepath.Join(dest, file))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), expected) {
				t.Fatalf("%s: expected %q, got %s", file, expected, data)
			}
		})
	}
}

func TestGenerateGoModule(t *testing.T) {
	for _, placeholder := range []string{"go-foundations", "example.com/go-foundations"} {
		for _, module := range []string{"", "example.org/acme/service", "example.org/acme/go-foundations"} {
			t.Run(placeholder+"/"+module, func(t *testing.T) {
				dest := t.TempDir()
				g := Generator{FS: fstest.MapFS{
					"template/go.mod.tmpl":  &fstest.MapFile{Data: []byte("module " + placeholder + "\n")},
					"template/main.go.tmpl": &fstest.MapFile{Data: []byte("package main\nimport _ \"" + placeholder + "/internal/app\"\n")},
				}}
				err := g.Generate(registry.Template{ID: "go", Kind: "go", EmbeddedDir: "template"}, Options{Name: "demo", Destination: dest, Module: module, NoHooks: true})
				if err != nil {
					t.Fatal(err)
				}
				want := defaultString(module, "demo")
				for file, expected := range map[string]string{
					"go.mod":  "module " + want + "\n",
					"main.go": "\"" + want + "/internal/app\"",
				} {
					data, err := os.ReadFile(filepath.Join(dest, file))
					if err != nil {
						t.Fatal(err)
					}
					if !strings.Contains(string(data), expected) {
						t.Fatalf("%s: expected module %q, got %s", file, want, data)
					}
				}
			})
		}
	}
}

func TestSanitize(t *testing.T) {
	cases := map[string]string{
		"Users Service": "users_service",
		"mediguide-app": "mediguide_app",
		"simple":        "simple",
	}
	for in, want := range cases {
		if got := sanitize(in); got != want {
			t.Fatalf("sanitize(%q)=%q want %q", in, got, want)
		}
	}
}

func TestWebTemplatesSubstituteProjectSettings(t *testing.T) {
	for _, name := range []string{"vanilla-js", "react"} {
		t.Run(name, func(t *testing.T) {
			tmpl, err := registry.Find(name)
			if err != nil {
				t.Fatal(err)
			}
			dest := t.TempDir()
			if err := (Generator{FS: foundry.Templates}).Generate(tmpl, Options{Name: "Customer Portal", Destination: dest, NoHooks: true}); err != nil {
				t.Fatal(err)
			}
			for _, file := range []string{"package.json", "package-lock.json"} {
				data, err := os.ReadFile(filepath.Join(dest, file))
				if err != nil {
					t.Fatal(err)
				}
				var manifest struct {
					Name string `json:"name"`
				}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatal(err)
				}
				if manifest.Name != "customer_portal" {
					t.Fatalf("%s: unexpected package name %q", file, manifest.Name)
				}
				if strings.Contains(string(data), "__PROJECT_NAME__") {
					t.Errorf("unreplaced name in %s", file)
				}
			}
			for _, file := range []string{"index.html", "src/shared/config/env.js", "README.md"} {
				data, err := os.ReadFile(filepath.Join(dest, file))
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(data), "customer_portal") || strings.Contains(string(data), "__PROJECT_NAME__") {
					t.Errorf("project name missing from %s", file)
				}
			}
			config, err := config.Read(dest)
			if err != nil {
				t.Fatal(err)
			}
			if config.Kind != tmpl.Kind || config.Template != tmpl.ID {
				t.Fatalf("wrong project metadata: %#v", config)
			}
		})
	}
}
