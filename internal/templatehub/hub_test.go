package templatehub

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdateAndResolveRemoteRelease(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	var z bytes.Buffer
	zw := zip.NewWriter(&z)
	w, err := zw.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte("hello"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(z.Bytes())
	checksum := hex.EncodeToString(sum[:])
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/template.zip" {
			w.Write(z.Bytes())
			return
		}
		m := Manifest{Schema: 1, Templates: []TemplateRelease{{ID: "remote-go", Name: "Remote Go", Version: "2.0.0", Kind: "go", Aliases: []string{"rg"}, URL: server.URL + "/template.zip", SHA256: checksum}}}
		_ = json.NewEncoder(w).Encode(m)
	}))
	defer server.Close()
	if _, err := Update(server.URL + "/registry.json"); err != nil {
		t.Fatal(err)
	}
	resolved, err := Resolve("rg", "2.0.0", false)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.ArchivePath == "" || resolved.Source != "remote" {
		t.Fatalf("unexpected resolution: %#v", resolved)
	}
	if _, err := os.Stat(resolved.ArchivePath); err != nil {
		t.Fatal(err)
	}
	cached, err := Resolve("rg", "2.0.0", true)
	if err != nil {
		t.Fatal(err)
	}
	if cached.Source != "cache" {
		t.Fatalf("expected cache, got %s", cached.Source)
	}
}
