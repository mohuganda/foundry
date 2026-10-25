package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jabahum/foundry/internal/config"
)

func TestFindNearestWorkspaceAndValidateMetadata(t *testing.T) {
	root := t.TempDir()
	if err := config.Write(root, config.Config{Template: "react-foundations", Kind: "react"}); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "src", "features")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := Find(child)
	if err != nil || p.Root != root || p.Kind != "react" {
		t.Fatalf("%+v, %v", p, err)
	}
	if _, err := At(child); !errors.Is(err, ErrNotFound) {
		t.Fatalf("At must not search parents: %v", err)
	}
	if err := config.Write(root, config.Config{Template: "react-foundations", Kind: "go"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Find(child); err == nil {
		t.Fatal("inconsistent metadata accepted")
	}
}

func TestNativeWorkspaceMarkers(t *testing.T) {
	for _, tc := range []struct{ file, contents, kind string }{
		{"go.mod", "module acme.org/api\n\ngo 1.23\n", "go"},
		{"Cargo.toml", "[package]\nname=\"demo\"\n", "rust"},
		{"package.json", `{"dependencies":{"react":"19"}}`, "react"},
		{"package.json", `{}`, "javascript"},
	} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.contents), 0o644); err != nil {
			t.Fatal(err)
		}
		p, err := At(root)
		if err != nil || p.Kind != tc.kind {
			t.Fatalf("%+v, %v", p, err)
		}
	}
}
