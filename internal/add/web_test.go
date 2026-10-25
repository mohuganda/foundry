package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func webWorkspace(t *testing.T, react bool) string {
	t.Helper()
	root := t.TempDir()
	manifest := `{"type":"module"}`
	if react {
		manifest = `{"type":"module","dependencies":{"react":"19"}}`
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestWebFeatureDispatchDryRunAndCollision(t *testing.T) {
	for _, react := range []bool{false, true} {
		root := webWorkspace(t, react)
		base := filepath.Join(root, "src", "features", "order_history")
		component := filepath.Join(base, "order_history.js")
		if react {
			component = filepath.Join(base, "components", "OrderHistory.jsx")
		}
		r, err := Feature(root, "Order History", Options{DryRun: true})
		if err != nil || len(r.Files) != 3 {
			t.Fatalf("%+v, %v", r, err)
		}
		if _, err := os.Stat(filepath.Join(root, "src")); !os.IsNotExist(err) {
			t.Fatal("dry-run wrote source")
		}
		if _, err := Feature(root, "order-history"); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(component)
		if err != nil || !strings.Contains(string(data), "Order History") {
			t.Fatalf("%s, %v", data, err)
		}
		if _, err := Feature(root, "order-history"); err == nil {
			t.Fatal("overwrote existing feature")
		}
		if got, _ := os.ReadFile(component); string(got) != string(data) {
			t.Fatal("collision changed source")
		}
	}
}

func TestWebFeatureRejectsInvalidNamesAndFlutterFlags(t *testing.T) {
	root := webWorkspace(t, false)
	for _, name := range []string{"", "../outside", "3orders", "orders!", "con", "nul", strings.Repeat("a", 65)} {
		if _, err := Feature(root, name); err == nil {
			t.Fatalf("accepted %q", name)
		}
	}
	if _, err := Feature(root, "orders", Options{App: "customer"}); err == nil {
		t.Fatal("accepted Flutter-only app flag")
	}
	if entries, _ := os.ReadDir(root); len(entries) != 1 {
		t.Fatal("invalid generation created files")
	}
}

func TestWebLateCollisionDoesNotCreateEarlierFiles(t *testing.T) {
	root := webWorkspace(t, true)
	collision := filepath.Join(root, "src", "features", "orders", "components", "Orders.test.jsx")
	if err := os.MkdirAll(filepath.Dir(collision), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(collision, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Feature(root, "orders"); err == nil {
		t.Fatal("collision accepted")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(collision), "Orders.jsx")); !os.IsNotExist(err) {
		t.Fatal("collision left component")
	}
	if _, err := os.Stat(filepath.Join(root, "src", "features", "orders", "README.md")); !os.IsNotExist(err) {
		t.Fatal("collision left README")
	}
}
