package publish

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPack(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "x.zip")
	sha, err := Pack(dir, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(sha) != 64 {
		t.Fatalf("sha=%q", sha)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}
