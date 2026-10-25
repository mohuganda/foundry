// Package authconfig renders non-secret, optional authentication configuration.
package authconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jabahum/foundry/internal/scaffold"
)

type Client struct {
	ClientID              string `json:"client_id"`
	RedirectURI           string `json:"redirect_uri"`
	PostLogoutRedirectURI string `json:"post_logout_redirect_uri"`
}
type Config struct {
	Schema       int      `json:"schema"`
	Mode         string   `json:"mode"`
	Issuer       string   `json:"issuer"`
	Audience     string   `json:"audience"`
	TokenProfile string   `json:"token_profile"`
	Scopes       []string `json:"scopes"`
	Web          Client   `json:"web"`
	Native       Client   `json:"native"`
}

func Default() Config {
	return Config{Schema: 1, Mode: "none", TokenProfile: "keycloak", Scopes: []string{"openid", "profile", "email"}}
}
func Resolve(mode, file, kind string) (Config, error) {
	c := Default()
	if file != "" {
		st, err := os.Stat(file)
		if err != nil {
			return c, err
		}
		if !st.Mode().IsRegular() || st.Size() > 65536 {
			return c, fmt.Errorf("auth configuration must be a regular JSON file up to 64 KiB")
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return c, err
		}
		if !bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			return c, fmt.Errorf("auth configuration must be an object")
		}
		d := json.NewDecoder(bytes.NewReader(data))
		d.DisallowUnknownFields()
		if err = d.Decode(&c); err != nil {
			return c, err
		}
		if d.Decode(new(any)) != io.EOF {
			return c, fmt.Errorf("expected one auth object")
		}
	}
	if mode != "" {
		c.Mode = mode
	}
	return c, c.Validate(kind)
}
func secureURL(value string, native bool) bool {
	u, e := url.Parse(value)
	if e != nil || u.User != nil || u.Fragment != "" || u.RawQuery != "" || strings.ContainsAny(value, "\\ \t\n\r") {
		return false
	}
	if u.Scheme == "https" && u.Host != "" {
		return true
	}
	return native && u.Scheme != "" && u.Scheme != "http" && u.Scheme != "javascript" && u.Scheme != "data" && u.Host != "" && u.Path != ""
}
func (c Config) Validate(kind string) error {
	switch kind {
	case "go", "flutter", "react", "javascript", "rust":
	default:
		return fmt.Errorf("unsupported authentication workspace %s", kind)
	}
	if c.Schema != 1 || (c.Mode != "none" && c.Mode != "oidc") {
		return fmt.Errorf("auth schema must be 1; mode must be none or oidc")
	}
	if kind == "rust" && c.Mode != "none" {
		return fmt.Errorf("OIDC is supported for Go, web and Flutter; Rust CLI authentication is not implemented")
	}
	if c.TokenProfile != "keycloak" && c.TokenProfile != "rfc9068" {
		return fmt.Errorf("token_profile must be keycloak or rfc9068")
	}
	for _, s := range []string{c.Issuer, c.Audience, c.Web.ClientID, c.Native.ClientID} {
		if len(s) > 2048 || strings.IndexFunc(s, func(r rune) bool { return r < 32 || r == 127 }) >= 0 {
			return fmt.Errorf("invalid auth value")
		}
	}
	if c.Mode == "none" {
		return nil
	}
	if !secureURL(c.Issuer, false) || c.Audience == "" {
		return fmt.Errorf("OIDC requires an HTTPS issuer and API audience")
	}
	if len(c.Scopes) == 0 || len(c.Scopes) > 16 {
		return fmt.Errorf("configure OIDC scopes including openid")
	}
	openid := false
	for _, s := range c.Scopes {
		if s == "openid" {
			openid = true
		}
		if s == "" || len(s) > 128 || strings.ContainsAny(s, " \t\r\n\"") {
			return fmt.Errorf("invalid OIDC scope")
		}
	}
	if !openid {
		return fmt.Errorf("openid scope is required")
	}
	for _, client := range []struct {
		value    Client
		native   bool
		required bool
	}{{c.Web, false, kind == "javascript" || kind == "react"}, {c.Native, true, kind == "flutter"}} {
		v := client.value
		if !client.required && v.ClientID == "" && v.RedirectURI == "" && v.PostLogoutRedirectURI == "" {
			continue
		}
		if v.ClientID == "" || v.ClientID == c.Audience || !secureURL(v.RedirectURI, client.native) || !secureURL(v.PostLogoutRedirectURI, client.native) {
			return fmt.Errorf("client requires a distinct client_id, valid redirect_uri and post_logout_redirect_uri")
		}
		if v.RedirectURI == v.PostLogoutRedirectURI {
			return fmt.Errorf("sign-in and sign-out callback URLs must differ")
		}
		if !client.native {
			a, _ := url.Parse(v.RedirectURI)
			b, _ := url.Parse(v.PostLogoutRedirectURI)
			if a.Scheme != b.Scheme || a.Host != b.Host {
				return fmt.Errorf("web callbacks must share an origin")
			}
		}
	}
	return nil
}
func JSON(c Config) []byte {
	b, _ := json.MarshalIndent(c, "", "  ")
	// Keep short scope lists consistent with the bundled web formatter.
	scopes, _ := json.MarshalIndent(c.Scopes, "  ", "  ")
	entries := make([]string, len(c.Scopes))
	for i, s := range c.Scopes {
		v, _ := json.Marshal(s)
		entries[i] = string(v)
	}
	inline := "[" + strings.Join(entries, ", ") + "]"
	if len(inline)+12 <= 80 {
		b = bytes.Replace(b, scopes, []byte(inline), 1)
	}
	return append(b, '\n')
}

func Render(source fs.FS, c Config, kind string) (map[string]scaffold.File, error) {
	files := map[string]scaffold.File{}
	put := func(p string, b []byte) { files[filepath.FromSlash(p)] = scaffold.File{Data: b, Mode: 0644} }
	put("auth/auth.json", JSON(c))
	switch kind {
	case "go":
		put("authconfig/config.json", JSON(c))
		b, e := fs.ReadFile(source, "templates/shared/auth/go/config.go.tmpl")
		if e != nil {
			return nil, e
		}
		put("authconfig/config.go", b)
	case "javascript", "react":
		put("src/shared/auth/config.json", JSON(c))
	case "flutter":
		put("packages/foundation_auth/lib/src/generated_config.dart", []byte("// Generated by Foundry; edit auth/auth.json and run foundry auth configure.\n\nconst authConfigJson = r\"\"\"\n"+string(JSON(c))+"\"\"\";\n"))
	case "rust":
	default:
		return nil, fmt.Errorf("unsupported auth workspace %s", kind)
	}
	return files, nil
}
