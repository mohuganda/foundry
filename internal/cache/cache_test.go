package cache

import (
	"path/filepath"
	"testing"
)

func TestDirHonorsOverride(t *testing.T) {
	want := filepath.Join(t.TempDir(), "foundry-cache")
	t.Setenv("FOUNDRY_CACHE_DIR", want)

	got, err := Dir()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Dir() = %q, want %q", got, want)
	}
}
