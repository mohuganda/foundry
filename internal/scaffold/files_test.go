package scaffold

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPreflightAndDryRunNeverLeaveOutput(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "z.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	plan := map[string]File{"new/a.txt": {[]byte("new"), 0o644}, "z.txt": {[]byte("replace"), 0o644}}
	for _, dry := range []bool{false, true} {
		if _, err := Apply(root, plan, false, dry); err == nil {
			t.Fatal("collision accepted")
		}
		if _, err := os.Stat(filepath.Join(root, "new")); !os.IsNotExist(err) {
			t.Fatal("collision left partial output")
		}
	}
	absent := filepath.Join(t.TempDir(), "missing", "project")
	if _, err := Apply(absent, map[string]File{"a": {[]byte("new"), 0o644}}, false, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(absent)); !os.IsNotExist(err) {
		t.Fatal("dry-run created directories")
	}
}

func TestWriteFailureRollsBackNewAndReplacedFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	calls := 0
	_, err := apply(root, map[string]File{"a.txt": {[]byte("changed"), 0o644}, "new/z.txt": {[]byte("partial"), 0o644}}, true, false, func(f *os.File, data []byte) error {
		calls++
		if _, err := f.Write(data); err != nil {
			return err
		}
		if calls == 2 {
			return errors.New("simulated disk failure")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected write failure")
	}
	data, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(data) != "original" {
		t.Fatalf("original not restored: %s, %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, "new")); !os.IsNotExist(err) {
		t.Fatal("rollback left directories")
	}
	if runtime.GOOS != "windows" {
		st, _ := os.Stat(filepath.Join(root, "a.txt"))
		if st.Mode().Perm() != 0o600 {
			t.Fatal("original mode lost")
		}
	}
}

func TestRejectEscapingAndSymlinkPaths(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"../outside", filepath.Join(root, "absolute")} {
		if _, err := Apply(root, map[string]File{name: {[]byte("bad"), 0o644}}, true, false); err == nil {
			t.Fatalf("accepted %s", name)
		}
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skip("symlinks unavailable")
	}
	if _, err := Apply(root, map[string]File{"link/file": {[]byte("bad"), 0o644}}, true, false); err == nil {
		t.Fatal("followed symlink")
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Fatal("wrote outside root")
	}
}
