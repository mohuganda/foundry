package registry

import (
	"fmt"
	"sort"
	"strings"
)

type Template struct {
	ID          string
	Name        string
	Version     string
	Description string
	Aliases     []string
	EmbeddedDir string
	Kind        string
}

var templates = []Template{
	{
		ID: "go-foundations", Name: "Go Foundations", Version: "1.1.0",
		Description: "Go + gRPC + Protobuf + PostgreSQL + local JWT + optional OIDC + local RBAC + observability",
		Aliases:     []string{"go", "grpc", "backend"}, EmbeddedDir: "templates/go-foundations", Kind: "go",
	},
	{
		ID: "flutter-foundations", Name: "Flutter Foundations", Version: "1.1.0",
		Description: "Flutter + Riverpod + go_router + Dio + Melos + Mason-ready reusable foundation packages",
		Aliases:     []string{"flutter", "mobile", "app"}, EmbeddedDir: "templates/flutter_foundations", Kind: "flutter",
	},
	{
		ID: "vanilla-js-foundations", Name: "Vanilla JavaScript Foundations", Version: "1.1.0",
		Description: "JavaScript + Vite + feature-first modules + ESLint + Prettier + Vitest + accessible DOM components",
		Aliases:     []string{"vanilla", "vanilla-js", "javascript", "js"}, EmbeddedDir: "templates/vanilla-js-foundations", Kind: "javascript",
	},
	{
		ID: "react-foundations", Name: "React Foundations", Version: "1.1.0",
		Description: "React + JavaScript + Vite + feature-first components + Hooks linting + Testing Library + Vitest",
		Aliases:     []string{"react"}, EmbeddedDir: "templates/react-foundations", Kind: "react",
	},
	{
		ID: "rust-foundations", Name: "Rust Foundations", Version: "1.1.0",
		Description: "Rust 2024 + reusable library + CLI adapter + feature modules + typed errors + Cargo tests + Clippy",
		Aliases:     []string{"rust", "rs"}, EmbeddedDir: "templates/rust-foundations", Kind: "rust",
	},
}

func All() []Template {
	out := append([]Template(nil), templates...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func Find(name string) (Template, error) {
	q := strings.ToLower(strings.TrimSpace(name))
	for _, t := range templates {
		if q == t.ID || q == strings.ToLower(t.Name) {
			return t, nil
		}
		for _, a := range t.Aliases {
			if q == a {
				return t, nil
			}
		}
	}
	return Template{}, fmt.Errorf("unknown template %q", name)
}
