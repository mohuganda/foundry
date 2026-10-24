package cli

import (
	"bufio"
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
)

// Version is overridden from the release tag at build time. Keeping a useful
// fallback makes `go install` and local source builds self-describing.
var Version = "1.2.0"

type App struct{ FS fs.FS }

func (a App) Run(args []string) int {
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
		if len(args) < 2 {
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
		doctor.Print()
		return 0
	case "version", "--version", "-v":
		fmt.Println("foundry", Version)
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

func (a App) interactive() int {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Foundry v" + Version + " — project foundations")
	fmt.Print("Project name: ")
	name, _ := r.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Fprintln(os.Stderr, "project name required")
		return 2
	}
	fmt.Println("Template:\n  1) Go Foundations\n  2) Flutter Foundations")
	fmt.Print("Select [1]: ")
	choice, _ := r.ReadString('\n')
	choice = strings.TrimSpace(choice)
	tmpl := "go"
	if choice == "2" {
		tmpl = "flutter"
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
	tmpl := f.String("template", "", "template: go or flutter")
	version := f.String("version", "", "template version; defaults to bundled version")
	module := f.String("module", "", "Go module path")
	org := f.String("org", "", "Flutter organization, e.g. ug.go.health")
	pkg := f.String("package", "", "Flutter package/app id token")
	dest := f.String("destination", "", "output directory")
	force := f.Bool("force", false, "allow generation into non-empty directory")
	noHooks := f.Bool("no-hooks", false, "skip post-generation hooks")
	if err := f.Parse(args); err != nil {
		return 2
	}
	if name == "" && len(f.Args()) > 0 {
		name = f.Args()[0]
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: foundry new <name> --template <go|flutter>")
		return 2
	}
	if *tmpl == "" {
		fmt.Fprintln(os.Stderr, "--template is required (go or flutter)")
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
	return 0
}

func addCmd(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: foundry add <service|repository|usecase|proto|migration|feature|package> <name> [--project path]")
		return 2
	}
	kind, name := args[0], args[1]
	f := flag.NewFlagSet("add", flag.ContinueOnError)
	f.SetOutput(os.Stderr)
	project := f.String("project", ".", "project root")
	if err := f.Parse(args[2:]); err != nil {
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
		uc := "execute"
		if extra := f.Args(); len(extra) > 0 {
			uc = extra[0]
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
	if len(args) == 0 || args[0] == "list" {
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
		return 2
	}
	_ = check
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

Usage:
  foundry
  foundry new <name> --template go [--module path] [--version version]
  foundry new <name> --template flutter [--org com.example] [--version version]
  foundry add service <name> [--project path]
  foundry add repository <name> [--project path]
  foundry add usecase <domain> <name> [--project path]
  foundry add proto <name> [--project path]
  foundry add migration <name> [--project path]
  foundry add feature <name> [--project path]
  foundry add package <name> [--project path]
  foundry template list
  foundry upgrade [--project path] [--check]
  foundry mason <brick> [args...]
  foundry list
  foundry info <template>
  foundry doctor
  foundry version`)
}
