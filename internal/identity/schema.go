// Package identity manages non-secret organisation identity and generated bindings.
package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const Source = "identity/identity.json"

type Organisation struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	ShortName string `json:"short_name"`
	Website   string `json:"website"`
}
type Application struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Department   string `json:"department"`
	Programme    string `json:"programme"`
	SupportEmail string `json:"support_email"`
	Environment  string `json:"environment"`
}
type Assets struct {
	LogoLight       string `json:"logo_light"`
	LogoDark        string `json:"logo_dark"`
	CoatOfArms      string `json:"coat_of_arms"`
	ApplicationIcon string `json:"application_icon"`
	Favicon         string `json:"favicon"`
}
type Palette struct {
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	Muted      string `json:"muted"`
	Primary    string `json:"primary"`
	OnPrimary  string `json:"on_primary"`
	Focus      string `json:"focus"`
	Border     string `json:"border"`
	Success    string `json:"success"`
	Error      string `json:"error"`
}
type Tokens struct {
	Light       Palette `json:"light"`
	Dark        Palette `json:"dark"`
	FontFamily  string  `json:"font_family"`
	FontSize    int     `json:"font_size"`
	SpaceSmall  int     `json:"space_small"`
	SpaceMedium int     `json:"space_medium"`
	SpaceLarge  int     `json:"space_large"`
	Radius      int     `json:"radius"`
	BorderWidth int     `json:"border_width"`
}
type Approval struct {
	Status     string `json:"status"`
	Reference  string `json:"reference"`
	ApprovedBy string `json:"approved_by"`
	ApprovedOn string `json:"approved_on"`
}
type Branding struct {
	Assets   Assets   `json:"assets"`
	Tokens   Tokens   `json:"tokens"`
	Approval Approval `json:"approval"`
}
type Documents struct {
	TitlePrefix    string `json:"title_prefix"`
	FooterText     string `json:"footer_text"`
	DefaultVersion string `json:"default_version"`
	DateFormat     string `json:"date_format"`
}
type Links struct {
	Privacy       string `json:"privacy"`
	Accessibility string `json:"accessibility"`
	Support       string `json:"support"`
}
type Link struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}
type Ownership struct {
	TechnicalOwner      string   `json:"technical_owner"`
	ServiceOwner        string   `json:"service_owner"`
	RepositoryURL       string   `json:"repository_url"`
	OperationalContacts []string `json:"operational_contacts"`
}
type Config struct {
	Schema       int          `json:"schema"`
	Profile      string       `json:"profile"`
	Organisation Organisation `json:"organisation"`
	Application  Application  `json:"application"`
	Branding     Branding     `json:"branding"`
	Documents    Documents    `json:"documents"`
	Links        Links        `json:"links"`
	Navigation   []Link       `json:"navigation"`
	Ownership    *Ownership   `json:"ownership,omitempty"`
}

func Default(profile, name string) (Config, error) {
	if profile == "" {
		profile = "generic"
	}
	if profile != "generic" && profile != "moh-uganda" {
		return Config{}, fmt.Errorf("unknown identity profile %q; use generic or moh-uganda", profile)
	}
	c := Config{Schema: 1, Profile: profile, Organisation: Organisation{Name: "Your organisation"}, Application: Application{Name: name, Description: "Application built with Foundry", Environment: "development"}, Navigation: []Link{}, Ownership: &Ownership{OperationalContacts: []string{}}, Documents: Documents{DefaultVersion: "draft", DateFormat: "yyyy-MM-dd"}}
	c.Branding.Approval.Status = "pending"
	c.Branding.Tokens = Tokens{
		Light:      Palette{"#f4f5f1", "#ffffff", "#182723", "#52645c", "#1c624a", "#ffffff", "#9c4700", "#9aa89f", "#12683b", "#a31515"},
		Dark:       Palette{"#111a17", "#192620", "#f2f7f4", "#b7c9bf", "#a4e0bb", "#10251a", "#ffd18a", "#74877b", "#a4e0bb", "#ffb4b4"},
		FontFamily: "Arial, Helvetica, sans-serif", FontSize: 16, SpaceSmall: 8, SpaceMedium: 16, SpaceLarge: 24, Radius: 8, BorderWidth: 1,
	}
	if profile == "moh-uganda" {
		c.Organisation = Organisation{"Ministry of Health", "Republic of Uganda", "MoH", "https://health.go.ug"}
		c.Branding.Assets.CoatOfArms = "builtin:uganda-coat-of-arms"
	}
	return c, nil
}

// Resolve overlays a JSON document onto profile defaults and copies only declared image assets.
func Resolve(profile, file, name string) (Config, map[string][]byte, error) {
	var data []byte
	var err error
	if file != "" {
		st, statErr := os.Stat(file)
		if statErr != nil {
			return Config{}, nil, statErr
		}
		if !st.Mode().IsRegular() || st.Size() > 1024*1024 {
			return Config{}, nil, fmt.Errorf("identity configuration must be a regular JSON file no larger than 1 MiB")
		}
		data, err = os.ReadFile(file)
		if err != nil {
			return Config{}, nil, err
		}
		if len(data) > 1024*1024 {
			return Config{}, nil, fmt.Errorf("identity configuration exceeds 1 MiB")
		}
		if !bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			return Config{}, nil, fmt.Errorf("identity configuration must be a JSON object")
		}
		var hint struct {
			Profile string `json:"profile"`
		}
		if err = json.Unmarshal(data, &hint); err != nil {
			return Config{}, nil, err
		}
		if profile == "" {
			profile = hint.Profile
		} else if hint.Profile != "" && hint.Profile != profile {
			return Config{}, nil, fmt.Errorf("identity file profile conflicts with --profile")
		}
	}
	c, err := Default(profile, name)
	if err != nil {
		return Config{}, nil, err
	}
	if len(data) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&c); err != nil {
			return Config{}, nil, fmt.Errorf("identity configuration: %w", err)
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			return Config{}, nil, fmt.Errorf("identity configuration must contain one JSON object")
		}
	}
	if err := c.Validate(); err != nil {
		return Config{}, nil, err
	}
	assets, err := readAssets(&c, filepath.Dir(file))
	return c, assets, err
}

func (c Config) Public() Config {
	c.Ownership = nil
	c.Branding.Approval.ApprovedBy = ""
	c.Branding.Approval.Reference = ""
	return c
}
func JSON(c Config) []byte { data, _ := json.MarshalIndent(c, "", "  "); return append(data, '\n') }

var hex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
var font = regexp.MustCompile(`^[a-zA-Z0-9 ,'"_-]+$`)

func (c Config) Validate() error {
	if c.Schema != 1 || (c.Profile != "generic" && c.Profile != "moh-uganda") {
		return fmt.Errorf("unsupported identity schema or profile")
	}
	if strings.TrimSpace(c.Organisation.Name) == "" || strings.TrimSpace(c.Application.Name) == "" {
		return fmt.Errorf("organisation and application names are required")
	}
	if len(c.Navigation) > 12 {
		return fmt.Errorf("identity navigation supports at most 12 items")
	}
	for _, s := range []string{c.Application.Environment, c.Branding.Approval.Status} {
		if s == "" {
			return fmt.Errorf("environment and branding approval status are required")
		}
	}
	switch c.Application.Environment {
	case "development", "staging", "production", "test":
	default:
		return fmt.Errorf("environment must be development, staging, production, or test")
	}
	if c.Branding.Approval.Status != "pending" && c.Branding.Approval.Status != "approved" {
		return fmt.Errorf("branding approval status must be pending or approved")
	}
	if c.Branding.Approval.Status == "approved" && (c.Branding.Approval.Reference == "" || c.Branding.Approval.ApprovedBy == "" || c.Branding.Approval.ApprovedOn == "") {
		return fmt.Errorf("approved branding requires a reference, approver, and date")
	}
	if c.Branding.Approval.ApprovedOn != "" {
		if _, err := time.Parse("2006-01-02", c.Branding.Approval.ApprovedOn); err != nil {
			return fmt.Errorf("branding approved_on must use yyyy-MM-dd")
		}
	}
	if c.Documents.DateFormat != "yyyy-MM-dd" {
		return fmt.Errorf("document date_format currently supports yyyy-MM-dd")
	}
	if c.Application.SupportEmail != "" {
		a, err := mail.ParseAddress(c.Application.SupportEmail)
		if err != nil || a.Address != c.Application.SupportEmail {
			return fmt.Errorf("support_email must be a plain email address")
		}
	}
	for _, href := range []string{c.Organisation.Website, c.Links.Privacy, c.Links.Accessibility, c.Links.Support} {
		if err := validLink(href); err != nil {
			return err
		}
	}
	for _, link := range c.Navigation {
		if strings.TrimSpace(link.Label) == "" || link.Href == "" {
			return fmt.Errorf("navigation requires label and href")
		}
		if err := validLink(link.Href); err != nil {
			return err
		}
	}
	for _, p := range []Palette{c.Branding.Tokens.Light, c.Branding.Tokens.Dark} {
		for _, colour := range []string{p.Background, p.Surface, p.Text, p.Muted, p.Primary, p.OnPrimary, p.Focus, p.Border, p.Success, p.Error} {
			if !hex.MatchString(colour) {
				return fmt.Errorf("identity colours must use #RRGGBB")
			}
		}
	}
	t := c.Branding.Tokens
	if !font.MatchString(t.FontFamily) || t.FontSize < 12 || t.FontSize > 24 || t.SpaceSmall < 0 || t.SpaceMedium < 0 || t.SpaceLarge < 0 || t.SpaceLarge > 128 || t.SpaceMedium > 128 || t.SpaceSmall > 128 || t.Radius < 0 || t.Radius > 64 || t.BorderWidth < 1 || t.BorderWidth > 4 {
		return fmt.Errorf("invalid identity typography, spacing, radius, or border tokens")
	}
	// Config is plain metadata: reject control characters, which could corrupt generated code/headers.
	var tree any
	_ = json.Unmarshal(JSON(c), &tree)
	return plainText(tree)
}

func plainText(value any) error {
	switch value := value.(type) {
	case string:
		if len(value) > 4096 || strings.IndexFunc(value, func(r rune) bool { return r < 32 || r == 127 }) >= 0 {
			return fmt.Errorf("identity values must be plain text without control characters (maximum 4096 bytes)")
		}
	case map[string]any:
		for _, v := range value {
			if err := plainText(v); err != nil {
				return err
			}
		}
	case []any:
		for _, v := range value {
			if err := plainText(v); err != nil {
				return err
			}
		}
	}
	return nil
}
func validLink(href string) error {
	if href == "" {
		return nil
	}
	if strings.ContainsAny(href, "\\ \r\n\t") || strings.HasPrefix(href, "//") {
		return fmt.Errorf("invalid identity link %q", href)
	}
	u, err := url.Parse(href)
	if err != nil || u.User != nil {
		return fmt.Errorf("invalid identity link")
	}
	if (strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#")) && u.Host == "" && u.Scheme == "" {
		return nil
	}
	if u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("identity links must use HTTPS, /paths, or #fragments")
	}
	return nil
}
