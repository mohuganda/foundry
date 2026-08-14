package projectconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Schema          int
	ProjectName     string
	Template        string
	TemplateVersion string
	Kind            string
	Module          string
	Org             string
	Package         string
}

func Path(root string) string { return filepath.Join(root, ".foundry.yaml") }

func Write(root string, c Config) error {
	if c.Schema == 0 {
		c.Schema = 1
	}
	content := fmt.Sprintf("schema: %d\nproject_name: %s\ntemplate: %s\ntemplate_version: %s\nkind: %s\nmodule: %s\norg: %s\npackage: %s\n", c.Schema, quote(c.ProjectName), quote(c.Template), quote(c.TemplateVersion), quote(c.Kind), quote(c.Module), quote(c.Org), quote(c.Package))
	return os.WriteFile(Path(root), []byte(content), 0o644)
}

func Read(root string) (Config, error) {
	f, err := os.Open(Path(root))
	if err != nil {
		return Config{}, fmt.Errorf(".foundry.yaml not found: %w", err)
	}
	defer f.Close()
	c := Config{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), "\"")
		switch k {
		case "schema":
			fmt.Sscanf(v, "%d", &c.Schema)
		case "project_name":
			c.ProjectName = v
		case "template":
			c.Template = v
		case "template_version":
			c.TemplateVersion = v
		case "kind":
			c.Kind = v
		case "module":
			c.Module = v
		case "org":
			c.Org = v
		case "package":
			c.Package = v
		}
	}
	if err := s.Err(); err != nil {
		return Config{}, err
	}
	if c.Template == "" {
		return Config{}, fmt.Errorf("invalid .foundry.yaml: template is required")
	}
	return c, nil
}

func quote(s string) string { return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"` }
