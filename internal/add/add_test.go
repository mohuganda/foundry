package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoMigrationUsesNextSequence(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "migrations"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "migrations", "000003_seed.up.sql"), []byte("--"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := GoMigration(root, "add orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Files) != 2 || !strings.Contains(strings.Join(r.Files, " "), "000004_add_orders") {
		t.Fatalf("unexpected files: %#v", r.Files)
	}
}

func TestGoServiceUsesProjectModule(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.org/acme/service\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := GoService(root, "orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Files) == 0 {
		t.Fatal("no files created")
	}
	b, err := os.ReadFile(filepath.Join(root, "api", "proto", "orders", "v1", "orders.proto"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "example.org/acme/service/gen/proto/orders/v1") {
		t.Fatalf("module not rendered: %s", b)
	}
}

func TestFlutterFeatureFindsWorkspaceApp(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "apps", "demo")
	if err := os.MkdirAll(app, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "pubspec.yaml"), []byte("name: demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := FlutterFeature(root, "patient search")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Files) != 4 {
		t.Fatalf("expected 4 files, got %d", len(r.Files))
	}
	if _, err := os.Stat(filepath.Join(app, "lib", "features", "patient_search", "presentation", "patient_search_page.dart")); err != nil {
		t.Fatal(err)
	}
}
