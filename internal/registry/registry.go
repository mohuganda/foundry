package registry

import (
	"fmt"
	"sort"
	"strings"
)

type Template struct {
	ID              string
	Name            string
	Version         string
	Description     string
	Aliases         []string
	EmbeddedDir     string
	EmbeddedArchive string
	Kind            string
}

var templates = []Template{
	{
		ID: "go-grpc-clean", Name: "Go gRPC Clean Architecture", Version: "1.0.0",
		Description: "Go + gRPC + Protobuf + PostgreSQL + local JWT + optional OIDC + local RBAC + observability",
		Aliases:     []string{"go", "grpc", "backend"}, EmbeddedArchive: "templates/go-grpc-clean.zip", Kind: "go",
	},
	{
		ID: "flutter-foundations", Name: "Flutter Foundations", Version: "1.0.0",
		Description: "Flutter + Riverpod + go_router + Dio + Melos + Mason-ready reusable foundation packages",
		Aliases:     []string{"flutter", "mobile", "app"}, EmbeddedDir: "templates/flutter-foundations", Kind: "flutter",
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
