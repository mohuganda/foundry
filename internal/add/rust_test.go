package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func rustWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Cargo.toml"), []byte("[package]\nname = \"demo\"\nedition = \"2024\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "src", "features"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src", "features", "mod.rs"), []byte("//! Existing feature registry.\npub mod greeting;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestRustFeatureDryRunAndExplicitRegistration(t *testing.T) {
	root := rustWorkspace(t)
	registry := filepath.Join(root, "src", "features", "mod.rs")
	before, err := os.ReadFile(registry)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Feature(root, "Order History", Options{DryRun: true})
	if err != nil || len(r.Files) != 4 {
		t.Fatalf("%+v, %v", r, err)
	}
	base := filepath.Join(root, "src", "features", "order_history")
	if _, err := os.Stat(base); !os.IsNotExist(err) {
		t.Fatal("dry-run created module")
	}
	r, err = Feature(root, "order-history")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"mod.rs", "domain.rs", "service.rs", "README.md"} {
		if _, err := os.Stat(filepath.Join(base, file)); err != nil {
			t.Fatal(err)
		}
	}
	if after, _ := os.ReadFile(registry); string(after) != string(before) {
		t.Fatal("existing module registry changed")
	}
	if !strings.Contains(strings.Join(r.Next, " "), "pub mod order_history;") {
		t.Fatal("module registration instructions missing")
	}
	if _, err := Feature(root, "order-history"); err == nil {
		t.Fatal("overwrote existing Rust module")
	}
}

func TestRustFeatureRejectsKeywordsAndWrongWorkspace(t *testing.T) {
	root := rustWorkspace(t)
	for _, name := range []string{"self", "Self", "super", "mod", "type", "async", "gen", "crate", "../other", "1orders"} {
		if _, err := Feature(root, name); err == nil {
			t.Fatalf("accepted %q", name)
		}
	}
	if _, err := RustFeature(webWorkspace(t, false), "orders"); err == nil {
		t.Fatal("Rust generator accepted web workspace")
	}
	if _, err := WebFeature(root, "orders"); err == nil {
		t.Fatal("web generator accepted Rust workspace")
	}
	entries, err := os.ReadDir(filepath.Join(root, "src", "features"))
	if err != nil || len(entries) != 1 {
		t.Fatal("invalid generation changed workspace")
	}
}

func TestRustLateCollisionDoesNotWriteEarlierFiles(t *testing.T) {
	root := rustWorkspace(t)
	base := filepath.Join(root, "src", "features", "orders")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "service.rs"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Feature(root, "orders"); err == nil {
		t.Fatal("collision accepted")
	}
	entries, err := os.ReadDir(base)
	if err != nil || len(entries) != 1 {
		t.Fatal("collision left partial module")
	}
}
