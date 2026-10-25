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
	for _, version := range []string{"", "1.1.0", "99.0.0"} {
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

func TestFlutterAppCLISelection(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "pubspec.yaml"), []byte("name: workspace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, app := range []string{"admin", "customer"} {
		dir := filepath.Join(root, "apps", app)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "pubspec.yaml"), []byte("name: "+app+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if code := addCmd([]string{"feature", "orders", "--project", root}); code != 1 {
		t.Fatal(code)
	}
	if code := addCmd([]string{"feature", "orders", "--app", "customer", "--project", root, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	if code := addCmd([]string{"feature", "orders", "--project", root, "--app", "customer"}); code != 0 {
		t.Fatal(code)
	}
	for _, kind := range []string{"package", "service", "migration"} {
		if code := addCmd([]string{kind, "oops", "--project", root, "--app", "customer"}); code != 2 {
			t.Fatalf("%s accepted --app: %d", kind, code)
		}
	}
}

func TestFeatureDispatchForWebProjects(t *testing.T) {
	for _, template := range []string{"vanilla-js", "react"} {
		root := filepath.Join(t.TempDir(), "project")
		app := App{FS: foundry.Templates}
		if code := app.Run([]string{"new", "demo", "--template", template, "--destination", root, "--no-hooks"}); code != 0 {
			t.Fatal(code)
		}
		if code := app.Run([]string{"add", "feature", "orders", "--project", root, "--dry-run"}); code != 0 {
			t.Fatal(code)
		}
		if _, err := os.Stat(filepath.Join(root, "src", "features", "orders")); !os.IsNotExist(err) {
			t.Fatal("dry-run wrote feature")
		}
		if code := app.Run([]string{"add", "feature", "orders", "--project", root}); code != 0 {
			t.Fatal(code)
		}
		if _, err := os.Stat(filepath.Join(root, "src", "features", "orders", "README.md")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFeatureDispatchForRustProject(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	app := App{FS: foundry.Templates}
	if code := app.Run([]string{"new", "demo", "--template", "rust", "--destination", root, "--no-hooks"}); code != 0 {
		t.Fatal(code)
	}
	if code := app.Run([]string{"add", "feature", "order-history", "--project", root, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	if code := app.Run([]string{"add", "feature", "order-history", "--project", root}); code != 0 {
		t.Fatal(code)
	}
	if _, err := os.Stat(filepath.Join(root, "src", "features", "order_history", "mod.rs")); err != nil {
		t.Fatal(err)
	}
	if code := app.Run([]string{"add", "feature", "self", "--project", root}); code != 1 {
		t.Fatal("Rust keyword accepted")
	}
	if code := app.Run([]string{"add", "feature", "orders", "--project", root, "--app", "customer"}); code != 1 {
		t.Fatal("Flutter option accepted for Rust")
	}
}

func TestIdentityGenerationAndSync(t *testing.T) {
	app := App{FS: foundry.Templates}
	root := filepath.Join(t.TempDir(), "portal")
	if code := app.Run([]string{"new", "Portal", "--template", "react", "--profile", "moh-uganda", "--destination", root, "--no-hooks"}); code != 0 {
		t.Fatal(code)
	}
	canonical := filepath.Join(root, "identity", "identity.json")
	data, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Ministry of Health") {
		t.Fatal("profile not applied")
	}
	data = []byte(strings.Replace(string(data), `"name": "Portal"`, `"name": "Changed Portal"`, 1))
	if err := os.WriteFile(canonical, data, 0o644); err != nil {
		t.Fatal(err)
	}
	public := filepath.Join(root, "src", "shared", "identity", "config.json")
	before, err := os.ReadFile(public)
	if err != nil {
		t.Fatal(err)
	}
	if code := app.Run([]string{"identity", "sync", "--project", root, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	if after, _ := os.ReadFile(public); string(after) != string(before) {
		t.Fatal("dry-run wrote bindings")
	}
	if code := app.Run([]string{"identity", "sync", "--project", root}); code != 0 {
		t.Fatal(code)
	}
	if after, _ := os.ReadFile(public); !strings.Contains(string(after), "Changed Portal") {
		t.Fatal("sync did not update bindings")
	}
	if code := app.Run([]string{"identity", "validate", "--file", canonical}); code != 0 {
		t.Fatal(code)
	}
}

func TestConfigurableAuthentication(t *testing.T) {
	app := App{FS: foundry.Templates}
	root := filepath.Join(t.TempDir(), "portal")
	if code := app.Run([]string{"new", "portal", "--template", "react", "--destination", root, "--no-hooks"}); code != 0 {
		t.Fatal(code)
	}
	source := filepath.Join(root, "auth", "auth.json")
	before, _ := os.ReadFile(source)
	if !strings.Contains(string(before), `"mode": "none"`) {
		t.Fatal("authentication enabled by default")
	}
	input := filepath.Join(t.TempDir(), "auth.json")
	data := `{"mode":"oidc","issuer":"https://auth.example.org/realms/test","audience":"api","web":{"client_id":"web","redirect_uri":"https://portal.example.org/auth/callback","post_logout_redirect_uri":"https://portal.example.org/auth/signed-out"}}`
	if err := os.WriteFile(input, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	if code := app.Run([]string{"auth", "configure", "--project", root, "--file", input, "--dry-run"}); code != 0 {
		t.Fatal(code)
	}
	unchanged, _ := os.ReadFile(source)
	if string(before) != string(unchanged) {
		t.Fatal("dry run changed configuration")
	}
	if code := app.Run([]string{"auth", "configure", "--project", root, "--file", input}); code != 0 {
		t.Fatal(code)
	}
	public, _ := os.ReadFile(filepath.Join(root, "src/shared/auth/config.json"))
	if !strings.Contains(string(public), `"mode": "oidc"`) {
		t.Fatal("OIDC not configured")
	}
	if code := app.Run([]string{"auth", "validate", "--file", input, "--template", "react"}); code != 0 {
		t.Fatal(code)
	}
	rust := filepath.Join(t.TempDir(), "rust")
	if code := app.Run([]string{"new", "tool", "--template", "rust", "--auth-config", input, "--destination", rust, "--no-hooks"}); code == 0 {
		t.Fatal("claimed Rust auth support")
	}
	if _, err := os.Stat(rust); !os.IsNotExist(err) {
		t.Fatal("unsupported auth wrote destination")
	}
}
