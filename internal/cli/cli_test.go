package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jabahum/foundry"
)

func TestBundledTemplateCommands(t *testing.T) {
	app := App{FS: foundry.Templates}
	for _, args := range [][]string{{"template", "list"}, {"list"}, {"info", "go-foundations"}} {
		if code := app.Run(args); code != 0 {
			t.Fatalf("%v: exit %d", args, code)
		}
	}
}

func TestNewUsesOnlyBundledVersion(t *testing.T) {
	t.Setenv("FOUNDRY_REGISTRY_URL", "http://127.0.0.1:1/registry.json")
	app := App{FS: foundry.Templates}
	for _, version := range []string{"", "1.0.0", "99.0.0"} {
		t.Run(version, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "project")
			args := []string{"new", "demo", "--template", "go", "--destination", dest, "--no-hooks"}
			if version != "" {
				args = append(args, "--version", version)
			}
			code := app.Run(args)
			if version == "99.0.0" {
				if code == 0 {
					t.Fatal("unbundled version unexpectedly succeeded")
				}
				if _, err := os.Stat(dest); !os.IsNotExist(err) {
					t.Fatal("unbundled version created output")
				}
				return
			}
			if code != 0 {
				t.Fatalf("new exit %d", code)
			}
			data, err := os.ReadFile(filepath.Join(dest, "go.mod"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "module demo\n") {
				t.Fatalf("unexpected module: %s", data)
			}
		})
	}
}

func TestRetiredDeliveryCommandsAreRejected(t *testing.T) {
	app := App{FS: foundry.Templates}
	for _, args := range [][]string{{"template", "pack"}, {"template", "update"}, {"cache", "list"}} {
		if code := app.Run(args); code != 2 {
			t.Fatalf("%v: exit %d, want 2", args, code)
		}
	}
}
