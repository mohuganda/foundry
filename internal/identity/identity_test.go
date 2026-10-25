package identity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jabahum/foundry"
)

func TestProfileDefaultsAndStrictOverrides(t *testing.T) {
	c, err := Default("moh-uganda", "Portal")
	if err != nil {
		t.Fatal(err)
	}
	if c.Organisation.Name != "Ministry of Health" || c.Branding.Approval.Status != "pending" || c.Branding.Assets.LogoLight != "" {
		t.Fatal("invented or missing profile defaults")
	}
	root := t.TempDir()
	file := filepath.Join(root, "identity.json")
	for _, input := range []string{`null`, `{"schema":2}`, `{"unexpected":true}`, `{"application":{"environment":"development staging"}}`, `{"links":{"privacy":"javascript:alert(1)"}}`, `{"branding":{"tokens":{"font_family":"x; color:red"}}}`} {
		if err := os.WriteFile(file, []byte(input), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := Resolve("moh-uganda", file, "Portal"); err == nil {
			t.Fatalf("accepted %s", input)
		}
	}
	if err := os.WriteFile(file, []byte(`{"application":{"programme":"Immunisation"},"links":{"privacy":"/privacy"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	c, _, err = Resolve("moh-uganda", file, "Portal")
	if err != nil {
		t.Fatal(err)
	}
	if c.Application.Name != "Portal" || c.Application.Programme != "Immunisation" || c.Organisation.ShortName != "MoH" {
		t.Fatal("partial overrides lost defaults")
	}
}

func TestPublicBindingsExcludeOwnershipAcrossAllTemplates(t *testing.T) {
	c, _ := Default("moh-uganda", "Portal")
	c.Ownership.TechnicalOwner = "PRIVATE-OWNER-REFERENCE"
	for _, kind := range []string{"go", "flutter", "javascript", "react", "rust"} {
		files, err := Render(foundry.Templates, c, nil, kind)
		if err != nil {
			t.Fatal(err)
		}
		for path, file := range files {
			if filepath.ToSlash(path) == Source {
				if !strings.Contains(string(file.Data), "PRIVATE-OWNER-REFERENCE") {
					t.Fatal("source ownership lost")
				}
				continue
			}
			if strings.Contains(string(file.Data), "PRIVATE-OWNER-REFERENCE") {
				t.Fatalf("private ownership leaked into %s/%s", kind, path)
			}
		}
	}
}

func TestAssetsAreCopiedAndActiveSVGIsRejected(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "identity.json")
	if err := os.WriteFile(file, []byte(`{"branding":{"assets":{"logo_light":"logo.svg"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	logo := filepath.Join(root, "logo.svg")
	if err := os.WriteFile(logo, []byte(`<svg xmlns="http://www.w3.org/2000/svg"><path d="M0 0"/></svg>`), 0o600); err != nil {
		t.Fatal(err)
	}
	c, assets, err := Resolve("moh-uganda", file, "Portal")
	if err != nil {
		t.Fatal(err)
	}
	if c.Branding.Assets.LogoLight != "assets/logo_light.svg" || len(assets) != 1 {
		t.Fatal("asset paths not normalized")
	}
	files, err := Render(foundry.Templates, c, assets, "react")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files[filepath.FromSlash("public/branding/assets/logo_light.svg")]; !ok {
		t.Fatal("asset not copied for browser")
	}
	for _, svg := range []string{`<svg><script/></svg>`, `<svg onload="alert(1)"/>`, `<svg><image href="https://example.org/tracker"/></svg>`, `<?xml-stylesheet href="external.css"?><svg/>`} {
		if err := os.WriteFile(logo, []byte(svg), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := Resolve("moh-uganda", file, "Portal"); err == nil {
			t.Fatal("active SVG accepted")
		}
	}
}

func TestRustStringsEscapeWithoutRawDelimiterLimits(t *testing.T) {
	c, _ := Default("moh-uganda", "Portal")
	c.Application.Name = "quoted\"" + strings.Repeat("#", 300) + `\path`
	output := rust(c)
	if strings.Contains(output, `r#"`) || !strings.Contains(output, `quoted\"`) || !strings.Contains(output, `\\path`) {
		t.Fatal("Rust string literals are not escaped")
	}
}
