package authconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOptionalAndStrictConfiguration(t *testing.T) {
	if _, e := Resolve("", "", "react"); e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(t.TempDir(), "auth.json")
	for _, input := range []string{`null`, `{"client_secret":"secret"}`, `{"mode":"oidc"}`, `{"mode":"saml"}`} {
		os.WriteFile(file, []byte(input), 0600)
		if _, e := Resolve("", file, "react"); e == nil {
			t.Fatalf("accepted %s", input)
		}
	}
	c := Default()
	c.Mode = "oidc"
	c.Issuer = "https://auth.example.org/realms/test"
	c.Audience = "api"
	c.Web = Client{"web", "https://portal.example.org/auth/callback", "https://portal.example.org/auth/signed-out"}
	c.Native = Client{"mobile", "ug.go.health.test://auth/callback", "ug.go.health.test://auth/signed-out"}
	for _, kind := range []string{"go", "react", "javascript", "flutter"} {
		if e := c.Validate(kind); e != nil {
			t.Fatalf("%s: %v", kind, e)
		}
	}
	if c.Validate("rust") == nil {
		t.Fatal("claimed unsupported Rust auth")
	}
	c.Web.ClientID = "api"
	if c.Validate("react") == nil {
		t.Fatal("accepted audience/client collision")
	}
	c.Web.ClientID = "web"
	c.Issuer = "http://auth.example.org"
	if c.Validate("go") == nil {
		t.Fatal("accepted cleartext issuer")
	}
}
