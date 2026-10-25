package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	addgen "github.com/jabahum/foundry/internal/add"
	"github.com/jabahum/foundry/internal/doctor"
	"github.com/jabahum/foundry/internal/generator"
	"github.com/jabahum/foundry/internal/registry"
	"github.com/jabahum/foundry/internal/upgrade"
	"github.com/jabahum/foundry/internal/workspace"
)

type App struct{ FS fs.FS }

func (a App) Run(args []string) int {
	if len(args) > 1 {
		switch args[0] {
		case "list", "version", "--version", "-v", "help", "--help", "-h":
			fmt.Fprintln(os.Stderr, "unexpected arguments; use foundry help")
			return 2
		}
	}
	if len(args) == 0 {
		return a.interactive()
	}
	switch args[0] {
	case "new":
		return a.newCmd(args[1:])
	case "add":
		return addCmd(args[1:])
	case "upgrade":
		return upgradeCmd(args[1:])
	case "mason":
		return masonCmd(args[1:])
	case "template", "templates":
		return templateCmd(args[1:])
	case "list":
		printTemplates()
		return 0
	case "info":
		if len(args) == 2 && (args[1] == "--help" || args[1] == "-h") {
			fmt.Println("Usage: foundry info <template>")
			return 0
		}
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: foundry info <template>")
			return 2
		}
		t, e := registry.Find(args[1])
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			return 1
		}
		printInfo(t)
		return 0
	case "doctor":
		return doctorCmd(args[1:])
	case "version", "--version", "-v":
		fmt.Println("foundry", CurrentVersion())
		return 0
	case "help", "--help", "-h":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		usage()
		return 2
	}
}

func doctorCmd(args []string) int {
	f := flag.NewFlagSet("doctor", flag.ContinueOnError)
	project := f.String("project", "", "workspace root; default searches current directory and parents")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments; use foundry doctor --help")
		return 2
	}
	var p workspace.Project
	var err error
	if *project != "" {
		p, err = workspace.At(*project)
	} else {
		p, err = workspace.Find(".")
	}
	if err != nil && (*project != "" || !errors.Is(err, workspace.ErrNotFound)) {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if !doctor.Print(os.Stdout, p, doctor.Inspect(p)) {
		return 1
	}
	return 0
}

func (a App) interactive() int {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Foundry " + CurrentVersion() + " — project foundations")
	fmt.Print("Project name: ")
	name, _ := r.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Fprintln(os.Stderr, "project name required")
		return 2
	}
	fmt.Println("Template:\n  1) Go Foundations\n  2) Flutter Foundations\n  3) Vanilla JavaScript Foundations\n  4) React Foundations\n  5) Rust Foundations")
	fmt.Print("Select [1]: ")
	choice, _ := r.ReadString('\n')
	choice = strings.TrimSpace(choice)
	tmpl := "go"
	switch choice {
	case "", "1":
	case "2":
		tmpl = "flutter"
	case "3":
		tmpl = "vanilla-js"
	case "4":
		tmpl = "react"
	case "5":
		tmpl = "rust"
	default:
		fmt.Fprintln(os.Stderr, "select a template from 1 to 5")
		return 2
	}
	return a.newCmd([]string{name, "--template", tmpl})
}

func (a App) newCmd(args []string) int {
	name := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		name = args[0]
		args = args[1:]
	}
	f := flag.NewFlagSet("new", flag.ContinueOnError)
	f.SetOutput(os.Stderr)
	tmpl := f.String("template", "", "template: go, flutter, vanilla-js, react, or rust")
	version := f.String("version", "", "template version; defaults to bundled version")
	module := f.String("module", "", "Go module path")
	org := f.String("org", "", "Flutter organization, e.g. ug.go.health")
	pkg := f.String("package", "", "Flutter package/app id token")
	dest := f.String("destination", "", "output directory")
	force := f.Bool("force", false, "allow generation into non-empty directory")
	noHooks := f.Bool("no-hooks", false, "skip post-generation hooks")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	remaining := f.Args()
	if name == "" && len(remaining) > 0 {
		name, remaining = remaining[0], remaining[1:]
	}
	if len(remaining) > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments; use: foundry new <name> [flags]")
		return 2
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: foundry new <name> --template <go|flutter|vanilla-js|react|rust>")
		return 2
	}
	if *tmpl == "" {
		fmt.Fprintln(os.Stderr, "--template is required (go, flutter, vanilla-js, react, or rust)")
		return 2
	}

	selected, err := registry.Find(*tmpl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if *version != "" && *version != selected.Version {
		fmt.Fprintf(os.Stderr, "template %q version %q is not bundled; this Foundry build includes %s\n", selected.ID, *version, selected.Version)
		return 1
	}
	g := generator.Generator{FS: a.FS}
	err = g.Generate(selected, generator.Options{Name: name, Destination: *dest, Module: *module, Org: *org, Package: *pkg, Force: *force, NoHooks: *noHooks})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Println("Template source: bundled files")
	if selected.Kind == "javascript" || selected.Kind == "react" {
		fmt.Println("Next: run npm ci, then npm run dev in the generated project.")
	} else if selected.Kind == "rust" {
		fmt.Println("Next: run cargo run --locked -- greet World in the generated project.")
	}
	return 0
}

func addCmd(args []string) int {
	if (len(args) == 1 && (args[0] == "--help" || args[0] == "-h")) || (len(args) == 2 && (args[1] == "--help" || args[1] == "-h")) {
		fmt.Println(`Go project generators:
  foundry add <service|repository|proto|migration> <name> [--project path]
  foundry add usecase <domain> [name] [--project path]
Flutter workspace generators:
  foundry add <feature|package> <name> [--project path]
Run from the generated root, or select it with --project. No add generators for JavaScript, React, or Rust.`)
		return 0
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: foundry add <service|repository|usecase|proto|migration|feature|package> <name> [--project path]")
		return 2
	}
	kind, name := args[0], args[1]
	f := flag.NewFlagSet("add", flag.ContinueOnError)
	f.SetOutput(os.Stderr)
	project := f.String("project", ".", "project root")
	rest := args[2:]
	usecase := ""
	if kind == "usecase" && len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		usecase, rest = rest[0], rest[1:]
	}
	if err := f.Parse(rest); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	remaining := f.Args()
	if kind == "usecase" && usecase == "" && len(remaining) > 0 {
		usecase, remaining = remaining[0], remaining[1:]
	}
	if len(remaining) > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments; use foundry add --help")
		return 2
	}
	root, err := filepath.Abs(*project)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	var result addgen.Result
	switch kind {
	case "service", "grpc-service":
		result, err = addgen.GoService(root, name)
	case "repository":
		result, err = addgen.GoRepository(root, name)
	case "usecase":
		domain := name
		uc := usecase
		if uc == "" {
			uc = "execute"
		}
		result, err = addgen.GoUsecase(root, domain, uc)
	case "proto":
		result, err = addgen.GoProto(root, name)
	case "migration":
		result, err = addgen.GoMigration(root, name)
	case "feature":
		result, err = addgen.FlutterFeature(root, name)
	case "package":
		result, err = addgen.FlutterPackage(root, name)
	default:
		fmt.Fprintf(os.Stderr, "unknown add generator %q\n", kind)
		return 2
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Printf("Added %s %q\n", kind, name)
	for _, p := range result.Files {
		fmt.Println("  +", p)
	}
	if len(result.Next) > 0 {
		fmt.Println("Next:")
		for _, n := range result.Next {
			fmt.Println("  -", n)
		}
	}
	return 0
}

func templateCmd(args []string) int {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println("Usage: foundry template list (also: foundry list)")
		return 0
	}
	if len(args) == 0 || (len(args) == 1 && args[0] == "list") {
		printTemplates()
		return 0
	}
	fmt.Fprintln(os.Stderr, "usage: foundry template list (templates ship with Foundry; install a newer Foundry release for updates)")
	return 2
}

func printTemplates() {
	fmt.Printf("%-22s %-8s %s\n", "TEMPLATE", "VERSION", "DESCRIPTION")
	for _, t := range registry.All() {
		fmt.Printf("%-22s %-8s %s\n", t.ID, t.Version, t.Description)
	}
}
func printInfo(t registry.Template) {
	fmt.Printf("%s\nID: %s\nVersion: %s\nAliases: %s\n%s\n", t.Name, t.ID, t.Version, strings.Join(t.Aliases, ", "), t.Description)
}
func upgradeCmd(args []string) int {
	f := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	f.SetOutput(os.Stderr)
	project := f.String("project", ".", "project root")
	check := f.Bool("check", true, "check for a newer template version")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if !*check {
		fmt.Fprintln(os.Stderr, "upgrade supports read-only check mode only; omit --check or set it to true")
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments; use foundry upgrade --help")
		return 2
	}
	root, err := filepath.Abs(*project)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	report, err := upgrade.Check(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Print(upgrade.Format(report))
	if report.UpdateAvailable {
		fmt.Println("No files changed. Foundry v1.2 only reports upgrade availability safely.")
	}
	return 0
}

func masonCmd(args []string) int {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println("Usage: foundry mason <brick> [mason arguments...] (runs mason make from PATH)")
		return 0
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: foundry mason <brick> [mason arguments...]")
		return 2
	}
	if _, err := exec.LookPath("mason"); err != nil {
		fmt.Fprintln(os.Stderr, "mason is not installed or not in PATH")
		return 1
	}
	cmdArgs := append([]string{"make", args[0]}, args[1:]...)
	cmd := exec.Command("mason", cmdArgs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "mason failed:", err)
		return 1
	}
	return 0
}

func usage() {
	fmt.Println(`Foundry — reusable project foundations

Discover (any directory):
  foundry help                              Show this guide
  foundry list                              List bundled templates
  foundry template list                     Alias for list
  foundry info <template>                    Show template details and aliases
  foundry doctor [--project path]           Check workspace tools and SDK versions
  foundry version                            Show executable version

Create (run from the parent directory):
  foundry                                   Choose a name and template interactively
  foundry new <name> --template go [--module path]
  foundry new <name> --template flutter [--org com.example]
  foundry new <name> --template <vanilla-js|react|rust>
  Options: --destination path, --version version, --no-hooks, --force
  Run foundry new --help for all creation flags.

Extend Go (generated project root, or --project path):
  foundry add service <name> [--project path]
  foundry add repository <name> [--project path]
  foundry add usecase <domain> [name] [--project path]
  foundry add proto <name> [--project path]
  foundry add migration <name> [--project path]

Extend Flutter (generated workspace root):
  foundry add feature <name> [--project path]
  foundry add package <name> [--project path]
  foundry mason <brick> [args...]             Run mason make from PATH in this directory
  For multi-app workspaces, select the app with Mason's --app argument.
  JavaScript, React, and Rust have no foundry add generators yet.

Maintain (any generated project root, or --project path):
  foundry upgrade [--project path] [--check]  Compare template versions; no files changed

Develop generated projects with their native tools:
  Go: make help                             Flutter: fvm dart run melos run
  JavaScript / React: npm run                Rust: bash scripts/check.sh
  Rust on Windows: ./scripts/check.ps1

Guide: docs/workspaces.md; full reference: docs/commands.md (Foundry repository).`)
}
