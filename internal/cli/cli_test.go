package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplatePackAcceptsDocumentedDirectoryFirstSyntax(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "template.zip")

	code := templateCmd([]string{
		"pack", dir,
		"--out", out,
		"--id", "example",
		"--version", "1.0.0",
		"--kind", "go",
	})
	if code != 0 {
		t.Fatalf("template pack exit code = %d, want 0", code)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("archive was not created at --out path: %v", err)
	}
}
