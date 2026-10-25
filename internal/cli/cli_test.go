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

func TestCommandHelpSucceedsWithoutGenerating(t *testing.T) {
	app := App{FS: foundry.Templates}
	for _, args := range [][]string{{"new", "--help"}, {"add", "--help"}, {"add", "service", "--help"}, {"template", "--help"}, {"upgrade", "--help"}, {"mason", "--help"}, {"info", "--help"}} {
		if code := app.Run(args); code != 0 {
			t.Fatalf("%v: exit %d, want 0", args, code)
		}
	}
}

func TestUsecaseProjectFlagBeforeOrAfterName(t *testing.T) {
	caller := t.TempDir()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(caller); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Error(err)
		}
	})
	for _, order := range []string{"name-first", "flags-first", "default-name"} {
		t.Run(order, func(t *testing.T) {
			project := t.TempDir()
			if err := os.WriteFile(filepath.Join(project, "go.mod"), []byte("module example.org/service\n\ngo 1.23\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			args := []string{"usecase", "orders"}
			want := "create"
			switch order {
			case "name-first":
				args = append(args, "create", "--project", project)
			case "flags-first":
				args = append(args, "--project", project, "create")
			case "default-name":
				args = append(args, "--project", project)
				want = "execute"
			}
			if code := addCmd(args); code != 0 {
				t.Fatalf("add %v: exit %d", args, code)
			}
			if _, err := os.Stat(filepath.Join(project, "internal", "application", "orders", want+".go")); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(caller, "internal")); !os.IsNotExist(err) {
				t.Fatal("generator wrote to the caller instead of --project")
			}
		})
	}
}

func TestUnexpectedArgumentsAndWriteUpgradeModeAreRejected(t *testing.T) {
	app := App{FS: foundry.Templates}
	for _, args := range [][]string{{"list", "extra"}, {"template", "list", "extra"}, {"info", "go", "extra"}, {"upgrade", "--check=false"}, {"add", "migration", "name", "unexpected"}, {"new", "demo", "--template", "rust", "extra"}} {
		if code := app.Run(args); code != 2 {
			t.Fatalf("%v: exit %d, want 2", args, code)
		}
	}
}

func TestDoctorHelpAndInvalidProject(t *testing.T) {
	app := App{FS: foundry.Templates}
	if code := app.Run([]string{"doctor", "--help"}); code != 0 {
		t.Fatal(code)
	}
	if code := app.Run([]string{"doctor", "--project", t.TempDir()}); code != 1 {
		t.Fatalf("unknown project: %d", code)
	}
	if code := app.Run([]string{"doctor", "extra"}); code != 2 {
		t.Fatal(code)
	}
}

func TestCLIDryRunsDoNotWrite(t *testing.T) {
	app := App{FS: foundry.Templates}
	root := t.TempDir()
	dest := filepath.Join(root, "new")
	if code := app.Run([]string{"new", "demo", "--template", "go", "--destination", dest, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatal("new dry-run wrote destination")
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module acme.org/api\ngo 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := app.Run([]string{"add", "service", "orders", "--project", root, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	if _, err := os.Stat(filepath.Join(root, "api")); !os.IsNotExist(err) {
		t.Fatal("add dry-run wrote source")
	}
}
