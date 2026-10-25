package identity

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func assetSlots(a *Assets) map[string]*string {
	return map[string]*string{"logo_light": &a.LogoLight, "logo_dark": &a.LogoDark, "coat_of_arms": &a.CoatOfArms, "application_icon": &a.ApplicationIcon, "favicon": &a.Favicon}
}
func readAssets(c *Config, base string) (map[string][]byte, error) {
	files := map[string][]byte{}
	for slot, value := range assetSlots(&c.Branding.Assets) {
		if *value == "" {
			continue
		}
		rel := filepath.FromSlash(*value)
		if !filepath.IsLocal(rel) {
			return nil, fmt.Errorf("asset %s must be relative to the identity file", slot)
		}
		path := base
		for _, part := range strings.Split(filepath.Clean(rel), string(filepath.Separator)) {
			path = filepath.Join(path, part)
			st, err := os.Lstat(path)
			if err != nil {
				return nil, err
			}
			if st.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("identity asset symlinks are not supported")
			}
		}
		st, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if !st.Mode().IsRegular() || st.Size() > 5*1024*1024 {
			return nil, fmt.Errorf("identity assets must be regular image files no larger than 5 MiB")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".svg":
			if err := safeSVG(data); err != nil {
				return nil, err
			}
		case ".png":
			if !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
				return nil, fmt.Errorf("invalid PNG asset")
			}
		case ".jpg", ".jpeg":
			if !bytes.HasPrefix(data, []byte{0xff, 0xd8, 0xff}) {
				return nil, fmt.Errorf("invalid JPEG asset")
			}
		case ".ico":
			if !bytes.HasPrefix(data, []byte{0, 0, 1, 0}) {
				return nil, fmt.Errorf("invalid icon asset")
			}
		default:
			return nil, fmt.Errorf("identity assets support SVG, PNG, JPEG, and ICO")
		}
		target := "assets/" + slot + ext
		files["identity/"+target] = data
		*value = target
	}
	return files, nil
}
func safeSVG(data []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	root := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("invalid SVG: %w", err)
		}
		switch t := token.(type) {
		case xml.ProcInst:
			if t.Target != "xml" {
				return fmt.Errorf("SVG processing instructions are not supported")
			}
		case xml.Directive:
			return fmt.Errorf("SVG directives are not supported")
		case xml.StartElement:
			if !root {
				if t.Name.Local != "svg" {
					return fmt.Errorf("asset is not an SVG")
				}
				root = true
			}
			name := strings.ToLower(t.Name.Local)
			if name == "script" || name == "foreignobject" || name == "style" || name == "animate" || name == "set" || name == "animatetransform" || name == "animatemotion" {
				return fmt.Errorf("SVG scripts, embedded HTML, and styles are not supported")
			}
			for _, a := range t.Attr {
				key := strings.ToLower(a.Name.Local)
				if strings.HasPrefix(key, "on") || key == "style" || strings.Contains(strings.ToLower(a.Value), "url(") {
					return fmt.Errorf("SVG active attributes are not supported")
				}
				if (key == "href" || key == "src") && !strings.HasPrefix(a.Value, "#") {
					return fmt.Errorf("SVG external references are not supported")
				}
			}
		}
	}
	if !root {
		return fmt.Errorf("empty SVG asset")
	}
	return nil
}
