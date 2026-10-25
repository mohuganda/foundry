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
	if c.Organisation.Name != "Ministry of Health" || c.Branding.Approval.Status != "pending" || c.Branding.Assets.CoatOfArms != "builtin:uganda-coat-of-arms" {
		t.Fatal("invented or missing profile defaults")
	}
	if c.Branding.Tokens.FontFamily != "Arial, Helvetica, sans-serif" {
		t.Fatal("profile typography is not sans-serif")
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
	if c.Branding.Assets.LogoLight != "assets/logo_light.svg" || len(assets) != 2 {
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

func TestSharedLoginDesignAndTheme(t *testing.T) {
	c, _ := Default("moh-uganda", "Programme Portal")
	c.Organisation.Name = `Ministry <script>alert(1)</script>`
	for _, kind := range []string{"go", "javascript", "react", "flutter", "rust"} {
		files, err := Render(foundry.Templates, c, nil, kind)
		if err != nil {
			t.Fatal(err)
		}
		design := files[filepath.FromSlash("identity/login-design.json")].Data
		if !strings.Contains(string(design), `"max_width": 440`) {
			t.Fatal("missing shared design")
		}
		root := "auth/keycloak-theme/foundry/login/"
		props := string(files[filepath.FromSlash(root+"theme.properties")].Data)
		if !strings.Contains(props, "parent=keycloak.v2") {
			t.Fatal("provider forms are not inherited")
		}
		messages := string(files[filepath.FromSlash(root+"messages/messages_en.properties")].Data)
		if strings.Contains(messages, "<script>") || !strings.Contains(messages, "&lt;script&gt;") {
			t.Fatal("organisation text injected HTML")
		}
		if !strings.Contains(string(files[filepath.FromSlash(root+"resources/css/foundry-tokens.css")].Data), "--login-max-width: 440px") {
			t.Fatal("provider design drift")
		}
		if kind == "javascript" || kind == "react" {
			if string(files[filepath.FromSlash("src/shared/identity/login.json")].Data) != string(design) {
				t.Fatal("web design drift")
			}
		}
		if kind == "flutter" {
			fixture, err := foundry.Templates.ReadFile("templates/flutter_foundations/packages/foundation_identity/lib/src/login_design.dart")
			if err != nil {
				t.Fatal(err)
			}
			if string(fixture) != string(files[filepath.FromSlash("packages/foundation_identity/lib/src/login_design.dart")].Data) {
				t.Fatal("Flutter login fixture differs from shared contract; regenerate it")
			}
		}
	}
}

func TestBundledMinistryLogoSurvivesOfflineGenerationAndSync(t *testing.T) {
	c, assets, err := Resolve("moh-uganda", "", "Portal")
	if err != nil {
		t.Fatal(err)
	}
	if c.Branding.Assets.CoatOfArms != "assets/coat_of_arms.png" {
		t.Fatal("bundled asset not normalized")
	}
	generic, genericAssets, err := Resolve("generic", "", "Portal")
	if err != nil || generic.Branding.Assets.CoatOfArms != "" || len(genericAssets) != 0 {
		t.Fatal("Ministry emblem leaked into generic profile")
	}
	for _, kind := range []string{"go", "react", "javascript", "flutter", "rust"} {
		files, err := Render(foundry.Templates, c, assets, kind)
		if err != nil {
			t.Fatal(err)
		}
		root := t.TempDir()
		for path, file := range files {
			target := filepath.Join(root, path)
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(target, file.Data, file.Mode); err != nil {
				t.Fatal(err)
			}
		}
		again, copied, err := Resolve("", filepath.Join(root, Source), "")
		if err != nil || again.Branding.Assets.CoatOfArms != c.Branding.Assets.CoatOfArms || len(copied) != 1 {
			t.Fatalf("%s sync lost logo: %v", kind, err)
		}
		theme := "auth/keycloak-theme/foundry/login/resources/"
		if !strings.Contains(string(files[filepath.FromSlash(theme+"css/foundry-tokens.css")].Data), "../assets/coat_of_arms.png") {
			t.Fatal("provider logo missing")
		}
		if len(files[filepath.FromSlash(theme+"assets/coat_of_arms.png")].Data) == 0 || len(files[filepath.FromSlash(theme+"icons/eye.svg")].Data) == 0 {
			t.Fatal("provider assets not bundled")
		}
	}
}
