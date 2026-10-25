package add

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jabahum/foundry/internal/workspace"
)

// Feature chooses the scaffold matching the selected workspace.
func Feature(root, name string, options ...Options) (Result, error) {
	p, err := workspace.At(root)
	if err != nil {
		return Result{}, err
	}
	if p.Kind != "flutter" && len(options) > 0 && options[0].App != "" {
		return Result{}, fmt.Errorf("--app is only supported for Flutter features")
	}
	switch p.Kind {
	case "flutter":
		return FlutterFeature(root, name, options...)
	case "javascript", "react":
		return WebFeature(root, name, options...)
	default:
		return Result{}, fmt.Errorf("feature generation is not supported for %s workspaces; use its supported generators", p.Kind)
	}
}

var featureIdentifier = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*([_ -][a-zA-Z0-9]+)*$`)

func featureName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 || !featureIdentifier.MatchString(value) {
		return "", fmt.Errorf("feature name must start with a letter and contain letters/digits separated by single spaces, hyphens, or underscores (maximum 64 characters)")
	}
	name := sanitize(value)
	reserved := " con prn aux nul com1 com2 com3 com4 com5 com6 com7 com8 com9 lpt1 lpt2 lpt3 lpt4 lpt5 lpt6 lpt7 lpt8 lpt9 "
	if strings.Contains(reserved, " "+name+" ") {
		return "", fmt.Errorf("feature name %q is reserved on Windows; choose another name", name)
	}
	return name, nil
}

func featureTitle(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
