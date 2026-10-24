package config

import (
	"os"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	root := t.TempDir()
	want := Config{ProjectName: "demo", Template: "go-foundations", TemplateVersion: "1.2.0", Kind: "go", Module: "example.com/demo"}
	if err := Write(root, want); err != nil {
		t.Fatal(err)
	}
	got, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Template != want.Template || got.TemplateVersion != want.TemplateVersion || got.Module != want.Module {
		t.Fatalf("unexpected config: %#v", got)
	}
	if _, err := os.Stat(Path(root)); err != nil {
		t.Fatal(err)
	}
}
