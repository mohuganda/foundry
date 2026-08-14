package registry

import "testing"

func TestFindAliases(t *testing.T) {
	for _, name := range []string{"go", "grpc", "flutter", "mobile"} {
		if _, err := Find(name); err != nil {
			t.Fatalf("Find(%q): %v", name, err)
		}
	}
}

func TestFindUnknown(t *testing.T) {
	if _, err := Find("unknown"); err == nil {
		t.Fatal("expected unknown template error")
	}
}
