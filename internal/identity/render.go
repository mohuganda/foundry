package identity

import (
	"encoding/json"
	"fmt"
	"go/format"
	"html"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jabahum/foundry/internal/scaffold"
)

// Render builds derived files only; UI source files are copied by project generation.
func Render(source fs.FS, c Config, assets map[string][]byte, kind string) (map[string]scaffold.File, error) {
	files := map[string]scaffold.File{}
	put := func(path string, data []byte) {
		files[filepath.FromSlash(path)] = scaffold.File{Data: data, Mode: 0o644}
	}
	design, err := fs.ReadFile(source, "templates/shared/identity/login.json")
	if err != nil {
		return nil, err
	}
	put("identity/login-design.json", design)
	var login struct {
		Eyebrow, Title, Description, Action, Notice string
		PendingAction                               string `json:"pending_action"`
		MaxWidth                                    int    `json:"max_width"`
		Padding                                     int
		ControlHeight                               int `json:"control_height"`
		TitleSize                                   int `json:"title_size"`
	}
	if err = json.Unmarshal(design, &login); err != nil {
		return nil, err
	}
	attribution, err := fs.ReadFile(source, "templates/shared/identity/assets/README.md")
	if err != nil {
		return nil, err
	}
	if c.Profile == "moh-uganda" {
		put("identity/ASSET-SOURCES.md", attribution)
	}
	put(Source, JSON(c))
	put("identity/public.json", JSON(c.Public()))
	for path, data := range assets {
		put(path, data)
	}
	themeRoot := "auth/keycloak-theme/foundry/login/"
	for _, file := range []struct{ source, target string }{{"theme.properties", "theme.properties"}, {"foundry.css", "resources/css/foundry.css"}} {
		data, e := fs.ReadFile(source, "templates/shared/identity/keycloak/"+file.source)
		if e != nil {
			return nil, e
		}
		put(themeRoot+file.target, data)
	}
	iconFiles, err := fs.ReadDir(source, "templates/shared/identity/keycloak/icons")
	if err != nil {
		return nil, err
	}
	for _, file := range iconFiles {
		data, err := fs.ReadFile(source, "templates/shared/identity/keycloak/icons/"+file.Name())
		if err != nil {
			return nil, err
		}
		put(themeRoot+"resources/icons/"+file.Name(), data)
	}
	tokens := css(c.Branding.Tokens) + fmt.Sprintf("\n:root {\n  --login-max-width: %dpx;\n  --login-padding: %dpx;\n  --login-control-height: %dpx;\n  --login-title-size: %dpx;\n}\n", login.MaxWidth, login.Padding, login.ControlHeight, login.TitleSize)
	lightLogo := c.Branding.Assets.LogoLight
	if lightLogo == "" {
		lightLogo = c.Branding.Assets.LogoDark
	}
	if lightLogo == "" {
		lightLogo = c.Branding.Assets.CoatOfArms
	}
	if lightLogo != "" {
		tokens += fmt.Sprintf("\n:root {\n  --foundry-logo: url('../%s');\n  --foundry-logo-height: 48px;\n  --foundry-logo-gap: 16px;\n}\n", lightLogo)
	}
	if c.Branding.Assets.LogoDark != "" {
		tokens += fmt.Sprintf("\n@media (prefers-color-scheme: dark) {\n  :root {\n    --foundry-logo: url('../%s');\n  }\n}\n", c.Branding.Assets.LogoDark)
	}
	put(themeRoot+"resources/css/foundry-tokens.css", []byte(tokens))
	property := func(value string) string {
		return strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\r", `\r`).Replace(value)
	}
	header := "<div class=\"kc-logo-text\"><span class=\"foundry-country\">" + html.EscapeString(c.Organisation.Country) + "</span><span class=\"foundry-organisation\">" + html.EscapeString(c.Organisation.Name) + "</span></div>"
	messages := "loginTitle=" + property(login.Title) + "\nloginTitleHtml=" + property(header) + "\nloginAccountTitle=" + property(login.Title) + "\ndoLogIn=" + property(login.Title) + "\n"
	put(themeRoot+"messages/messages_en.properties", []byte(messages))
	for path, data := range assets {
		put(themeRoot+"resources/"+strings.TrimPrefix(path, "identity/"), data)
	}

	switch kind {
	case "go":
		data, err := fs.ReadFile(source, "templates/shared/identity/go/profile.go.tmpl")
		if err != nil {
			return nil, err
		}
		data, err = format.Source(data)
		if err != nil {
			return nil, err
		}
		put("identity/profile.go", data)
	case "javascript", "react":
		public := c.Public()
		for _, value := range assetSlots(&public.Branding.Assets) {
			if *value != "" {
				*value = "/branding/" + *value
			}
		}
		put("src/shared/identity/config.json", JSON(public))
		put("src/shared/identity/login.json", design)
		put("src/shared/identity/tokens.css", []byte(css(c.Branding.Tokens)))
		for path, data := range assets {
			put("public/branding/"+strings.TrimPrefix(path, "identity/"), data)
		}
	case "rust":
		put("src/identity.rs", []byte(rust(c)))
	case "flutter":
		stringsJSON, _ := json.Marshal([]string{login.Eyebrow, login.Title, login.Description, login.Action, login.PendingAction, login.Notice})
		put("packages/foundation_identity/lib/src/login_design.dart", []byte(fmt.Sprintf("// Generated from the shared Foundry login design; do not edit.\n\nconst loginMaxWidth = %d.0;\nconst loginPadding = %d.0;\nconst loginControlHeight = %d.0;\nconst loginTitleSize = %d.0;\nconst loginCopyJson =\n    r\"\"\"%s\"\"\";\n", login.MaxWidth, login.Padding, login.ControlHeight, login.TitleSize, stringsJSON)))
		put("packages/foundation_identity/lib/src/generated_identity.dart", []byte("// Generated by foundry identity sync; edit identity/identity.json.\n\n/// Public identity snapshot without ownership metadata.\nconst identityJson = r\"\"\"\n"+string(JSON(c.Public()))+"\"\"\";\n"))
		for path, data := range assets {
			put("packages/foundation_identity/"+strings.TrimPrefix(path, "identity/"), data)
		}
	default:
		return nil, fmt.Errorf("identity does not support workspace kind %q", kind)
	}
	return files, nil
}
func css(t Tokens) string {
	palette := func(p Palette, indent string) string {
		var b strings.Builder
		for _, pair := range [][2]string{{"background", p.Background}, {"surface", p.Surface}, {"text", p.Text}, {"muted", p.Muted}, {"accent", p.Primary}, {"accent-hover", p.Primary}, {"on-primary", p.OnPrimary}, {"focus", p.Focus}, {"border", p.Border}, {"success", p.Success}, {"error", p.Error}} {
			fmt.Fprintf(&b, "%s--color-%s: %s;\n", indent, pair[0], strings.ToLower(pair[1]))
		}
		return b.String()
	}
	return fmt.Sprintf(":root {\n  color-scheme: light dark;\n%s  --font-family: %s;\n  --font-size: %dpx;\n  --space-small: %dpx;\n  --space-medium: %dpx;\n  --space-large: %dpx;\n  --radius: %dpx;\n  --border-width: %dpx;\n}\n\n@media (prefers-color-scheme: dark) {\n  :root {\n%s  }\n}\n", palette(t.Light, "  "), t.FontFamily, t.FontSize, t.SpaceSmall, t.SpaceMedium, t.SpaceLarge, t.Radius, t.BorderWidth, palette(t.Dark, "    "))
}
func rust(c Config) string {
	quote := func(value string) string {
		return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value) + `"`
	}
	return fmt.Sprintf(`//! Generated public identity. Edit identity/identity.json and run foundry identity sync.

/// Public metadata, excluding ownership and operational contacts.
pub const PUBLIC_JSON: &str = include_str!("../identity/public.json");
/// Application name.
pub const APPLICATION_NAME: &str = %s;
/// Organisation name.
pub const ORGANISATION: &str = %s;
/// Country or jurisdiction label.
pub const COUNTRY: &str = %s;
/// Display environment; separate from runtime deployment configuration.
pub const ENVIRONMENT: &str = %s;
/// Report-title prefix.
pub const TITLE_PREFIX: &str = %s;
/// Default report version.
pub const DEFAULT_VERSION: &str = %s;
/// Report footer.
pub const REPORT_FOOTER: &str = %s;

/// Build a plain-text report heading; dates/versions are explicit report inputs.
pub fn report_header(title: &str, version: &str, issued_on: &str) -> String {
    let version = if version.is_empty() {
        DEFAULT_VERSION
    } else {
        version
    };
    format!(
        "{COUNTRY}\n{ORGANISATION}\n{APPLICATION_NAME}\n{TITLE_PREFIX}{title}\nVersion: {version}\nDate: {issued_on}\n{REPORT_FOOTER}"
    )
}
`, quote(c.Application.Name), quote(c.Organisation.Name), quote(c.Organisation.Country), quote(c.Application.Environment), quote(c.Documents.TitlePrefix), quote(c.Documents.DefaultVersion), quote(c.Documents.FooterText))
}
