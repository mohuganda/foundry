package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jabahum/foundry"
	"github.com/jabahum/foundry/internal/registry"
)

func TestBundledTemplatesGenerate(t *testing.T) {
	for _, name := range []string{"go", "flutter"} {
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
			file, expected := "pubspec.yaml", "name: flutter_foundation_workspace"
			if name == "go" {
				file, expected = "go.mod", "module demo\n"
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
	for _, placeholder := range []string{"go-foundations", "example.com/go-foundations", "example.com/grpc-clean-starter"} {
		for _, module := range []string{"", "example.org/acme/service", "example.org/acme/go-foundations"} {
			t.Run(placeholder+"/"+module, func(t *testing.T) {
				dest := t.TempDir()
				g := Generator{FS: fstest.MapFS{
					"template/go.mod":  &fstest.MapFile{Data: []byte("module " + placeholder + "\n")},
					"template/main.go": &fstest.MapFile{Data: []byte("package main\nimport _ \"" + placeholder + "/internal/app\"\n")},
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
