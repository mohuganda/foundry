package add

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jabahum/foundry/internal/scaffold"
	"github.com/jabahum/foundry/internal/workspace"
)

type Options struct {
	DryRun bool
	App    string
}

type Result struct {
	Files []string
	Next  []string
}

func GoService(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("service name is required")
	}
	module, err := goModule(root)
	if err != nil {
		return Result{}, err
	}
	pascal := pascalCase(name)
	files := map[string]string{
		filepath.Join("api", "proto", name, "v1", name+".proto"): fmt.Sprintf(`syntax = "proto3";

package %s.v1;

option go_package = "%s/gen/proto/%s/v1;%sv1";

service %sService {
  rpc Get%s(Get%sRequest) returns (Get%sResponse);
}

message %s {
  string id = 1;
  string name = 2;
}

message Get%sRequest { string id = 1; }
message Get%sResponse { %s %s = 1; }
`, name, module, name, name, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, name),
		filepath.Join("internal", "domain", name, "entity.go"): fmt.Sprintf(`package %s

type %s struct {
	ID   string
	Name string
}
`, name, pascal),
		filepath.Join("internal", "domain", name, "repository.go"): fmt.Sprintf(`package %s

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*%s, error)
}
`, name, pascal),
		filepath.Join("internal", "application", name, "service.go"): fmt.Sprintf(`package %s

import (
	"context"

	domain "%s/internal/domain/%s"
)

type Service struct { repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Get(ctx context.Context, id string) (*domain.%s, error) {
	return s.repo.GetByID(ctx, id)
}
`, name, module, name, pascal),
		filepath.Join("internal", "transport", "grpc", name+"_handler.go"): fmt.Sprintf(`package grpc

import (
	"context"

	pb "%s/gen/proto/%s/v1"
	application "%s/internal/application/%s"
)

type %sHandler struct {
	pb.Unimplemented%sServiceServer
	service *application.Service
}

func New%sHandler(service *application.Service) *%sHandler { return &%sHandler{service: service} }

func (h *%sHandler) Get%s(ctx context.Context, req *pb.Get%sRequest) (*pb.Get%sResponse, error) {
	entity, err := h.service.Get(ctx, req.GetId())
	if err != nil { return nil, err }
	return &pb.Get%sResponse{%s: &pb.%s{Id: entity.ID, Name: entity.Name}}, nil
}
`, module, name, module, name, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, name, pascal, pascal),
	}
	written, err := writeFiles(root, files, "go", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"implement a persistence adapter for internal/domain/" + name + ".Repository", "run: make codegen", "wire New" + pascal + "Handler into cmd/server/main.go"}}, nil
}

func GoMigration(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("migration name is required")
	}
	dir := filepath.Join(root, "migrations")
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return Result{}, fmt.Errorf("%s is not a Go Foundry project (migrations directory missing)", root)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Result{}, err
	}
	max := 0
	re := regexp.MustCompile(`^(\d+)_`)
	for _, e := range entries {
		m := re.FindStringSubmatch(e.Name())
		if len(m) == 2 {
			n, _ := strconv.Atoi(m[1])
			if n > max {
				max = n
			}
		}
	}
	prefix := fmt.Sprintf("%06d", max+1)
	up := filepath.Join("migrations", prefix+"_"+name+".up.sql")
	down := filepath.Join("migrations", prefix+"_"+name+".down.sql")
	written, err := writeFiles(root, map[string]string{up: "-- Write forward migration here.\n", down: "-- Write rollback migration here.\n"}, "go", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"edit both migration files", "run: make migrate-up"}}, nil
}

func FlutterFeature(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("feature name is required")
	}
	selected := ""
	if len(options) > 0 {
		selected = options[0].App
	}
	appDir, err := flutterAppDir(root, selected)
	if err != nil {
		return Result{}, err
	}
	pascal := pascalCase(name)
	base := filepath.Join(appDir, "lib", "features", name)
	files := map[string]string{
		filepath.Join(base, "domain", name+".dart"): fmt.Sprintf(`class %s {
  const %s({required this.id});
  final String id;
}
`, pascal, pascal),
		filepath.Join(base, "application", name+"_controller.dart"): fmt.Sprintf(`import 'package:flutter_riverpod/flutter_riverpod.dart';

final %sControllerProvider = Provider<%sController>((ref) => %sController());

class %sController {
  const %sController();
}
`, camel(name), pascal, pascal, pascal, pascal),
		filepath.Join(base, "presentation", name+"_page.dart"): fmt.Sprintf(`import 'package:flutter/material.dart';

class %sPage extends StatelessWidget {
  const %sPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: Text('%s')),
    );
  }
}
`, pascal, pascal, pascal),
		filepath.Join(base, "data", name+"_repository.dart"): fmt.Sprintf(`abstract interface class %sRepository {
  Future<void> refresh();
}
`, pascal),
	}
	written, err := writeFiles(root, files, "flutter", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"add the feature route to your go_router configuration", "run: fvm dart run melos run bootstrap", "add feature-specific tests"}}, nil
}

func FlutterPackage(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("package name is required")
	}
	base := filepath.Join("packages", name)
	files := map[string]string{
		filepath.Join(base, "pubspec.yaml"):      fmt.Sprintf("name: %s\ndescription: Generated by Foundry.\nversion: 0.1.0\nenvironment:\n  sdk: '>=3.5.0 <4.0.0'\n", name),
		filepath.Join(base, "lib", name+".dart"): fmt.Sprintf("library %s;\n", name),
		filepath.Join(base, "README.md"):         fmt.Sprintf("# %s\n\nGenerated by Foundry on %s.\n", name, time.Now().UTC().Format("2006-01-02")),
	}
	written, err := writeFiles(root, files, "flutter", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"run: fvm dart run melos run bootstrap"}}, nil
}

func writeFiles(root string, files map[string]string, kind string, options ...Options) ([]string, error) {
	p, err := workspace.At(root)
	if err != nil {
		return nil, err
	}
	if p.Kind != kind {
		return nil, fmt.Errorf("%s generator requires a %s workspace; found %s", kind, kind, p.Kind)
	}
	marker := "go.mod"
	if kind == "flutter" {
		marker = "pubspec.yaml"
	}
	if st, err := os.Stat(filepath.Join(root, marker)); err != nil || !st.Mode().IsRegular() {
		return nil, fmt.Errorf("%s missing from workspace", marker)
	}
	plan := make(map[string]scaffold.File, len(files))
	for rel, data := range files {
		plan[rel] = scaffold.File{Data: []byte(data), Mode: 0o644}
	}
	dry := len(options) > 0 && options[0].DryRun
	return scaffold.Apply(root, plan, false, dry)
}

func goModule(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("go.mod not found: %w", err)
	}
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) == 2 && f[0] == "module" {
			return f[1], nil
		}
	}
	return "", fmt.Errorf("module directive not found in go.mod")
}

func flutterAppDir(root, selected string) (string, error) {
	if selected != "" && (selected == "." || selected == ".." || strings.ContainsAny(selected, `/\`) || !filepath.IsLocal(selected)) {
		return "", fmt.Errorf("--app must be a single directory name under apps/")
	}
	entries, err := os.ReadDir(filepath.Join(root, "apps"))
	if err != nil {
		return "", fmt.Errorf("apps directory not found: %w", err)
	}
	var apps []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		st, err := os.Lstat(filepath.Join(root, "apps", e.Name(), "pubspec.yaml"))
		if err == nil && st.Mode().IsRegular() {
			apps = append(apps, e.Name())
		}
	}
	if selected != "" {
		for _, app := range apps {
			if app == selected {
				return filepath.Join("apps", app), nil
			}
		}
		return "", fmt.Errorf("Flutter app %q not found; available apps: %s", selected, strings.Join(apps, ", "))
	}
	if len(apps) == 0 {
		return "", fmt.Errorf("no Flutter app found under apps/; generate/register an app first")
	}
	if len(apps) > 1 {
		return "", fmt.Errorf("multiple Flutter apps found (%s); select one with --app <name>", strings.Join(apps, ", "))
	}
	return filepath.Join("apps", apps[0]), nil
}

func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' {
			b.WriteRune(r)
		} else if r == '-' || r == ' ' {
			b.WriteByte('_')
		}
	}
	return strings.Trim(b.String(), "_")
}
func pascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}
func camel(s string) string {
	p := pascalCase(s)
	if p == "" {
		return ""
	}
	return strings.ToLower(p[:1]) + p[1:]
}

func GoRepository(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("repository name is required")
	}
	module, err := goModule(root)
	if err != nil {
		return Result{}, err
	}
	pascal := pascalCase(name)
	files := map[string]string{
		filepath.Join("internal", "domain", name, "repository.go"): fmt.Sprintf(`package %s

import "context"

type Repository interface {
	Create(ctx context.Context, entity *%s) error
	GetByID(ctx context.Context, id string) (*%s, error)
	List(ctx context.Context, limit int) ([]*%s, error)
}
`, name, pascal, pascal, pascal),
		filepath.Join("internal", "infrastructure", "persistence", "postgres", name+"_repository.go"): fmt.Sprintf(`package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "%s/internal/domain/%s"
)

type %sRepository struct { db *pgxpool.Pool }

var _ domain.Repository = (*%sRepository)(nil)

func New%sRepository(db *pgxpool.Pool) *%sRepository { return &%sRepository{db: db} }

func (r *%sRepository) Create(ctx context.Context, entity *domain.%s) error { return nil }
func (r *%sRepository) GetByID(ctx context.Context, id string) (*domain.%s, error) { return nil, nil }
func (r *%sRepository) List(ctx context.Context, limit int) ([]*domain.%s, error) { return nil, nil }
`, module, name, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal),
	}
	written, err := writeFiles(root, files, "go", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"define the domain entity in internal/domain/" + name, "implement PostgreSQL queries and error mapping", "wire New" + pascal + "Repository into cmd/server/main.go"}}, nil
}

func GoUsecase(root, domainName, usecase string, options ...Options) (Result, error) {
	domainName = sanitize(domainName)
	usecase = sanitize(usecase)
	if domainName == "" || usecase == "" {
		return Result{}, fmt.Errorf("domain and usecase names are required")
	}
	module, err := goModule(root)
	if err != nil {
		return Result{}, err
	}
	pascal := pascalCase(usecase)
	entity := pascalCase(domainName)
	files := map[string]string{
		filepath.Join("internal", "application", domainName, usecase+".go"): fmt.Sprintf(`package %s

import (
	"context"
	domain "%s/internal/domain/%s"
)

type %sInput struct{}

type %sUseCase struct { repo domain.Repository }

func New%sUseCase(repo domain.Repository) *%sUseCase { return &%sUseCase{repo: repo} }

func (u *%sUseCase) Execute(ctx context.Context, input %sInput) (*domain.%s, error) {
	return nil, nil
}
`, domainName, module, domainName, pascal, pascal, pascal, pascal, pascal, pascal, pascal, entity),
	}
	written, err := writeFiles(root, files, "go", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"implement the use-case business flow", "add unit tests for application/" + domainName + "/" + usecase + ".go"}}, nil
}

func GoProto(root, name string, options ...Options) (Result, error) {
	name = sanitize(name)
	if name == "" {
		return Result{}, fmt.Errorf("proto name is required")
	}
	module, err := goModule(root)
	if err != nil {
		return Result{}, err
	}
	pascal := pascalCase(name)
	files := map[string]string{filepath.Join("api", "proto", name, "v1", name+".proto"): fmt.Sprintf(`syntax = "proto3";

package %s.v1;
option go_package = "%s/gen/proto/%s/v1;%sv1";

service %sService {
  rpc Get(Get%sRequest) returns (Get%sResponse);
}

message Get%sRequest { string id = 1; }
message Get%sResponse { string id = 1; }
`, name, module, name, name, pascal, pascal, pascal, pascal, pascal)}
	written, err := writeFiles(root, files, "go", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{"run: make codegen", "implement and register the generated service server"}}, nil
}
