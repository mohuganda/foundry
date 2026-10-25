package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jabahum/foundry/internal/identity"
	"github.com/jabahum/foundry/internal/scaffold"
	"github.com/jabahum/foundry/internal/workspace"
)

func (a App) identityCmd(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: foundry identity example [--profile moh-uganda] [--name Portal]\n       foundry identity validate --file path.json\n       foundry identity sync [--project path] [--dry-run]")
		return 0
	}
	command := args[0]
	f := flag.NewFlagSet("identity "+command, flag.ContinueOnError)
	var profile, name, file, project *string
	var dry *bool
	switch command {
	case "example":
		profile = f.String("profile", "moh-uganda", "identity defaults")
		name = f.String("name", "Programme Portal", "application name")
	case "validate":
		file = f.String("file", "", "identity JSON to validate")
	case "sync":
		project = f.String("project", ".", "generated project root")
		dry = f.Bool("dry-run", false, "preview generated bindings without writing")
	default:
		fmt.Fprintln(os.Stderr, "unknown identity command; use foundry identity --help")
		return 2
	}
	if err := f.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected identity arguments")
		return 2
	}
	switch command {
	case "example":
		c, err := identity.Default(*profile, *name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Print(string(identity.JSON(c)))
		return 0
	case "validate":
		if *file == "" {
			fmt.Fprintln(os.Stderr, "--file is required")
			return 2
		}
		c, _, err := identity.Resolve("", *file, "")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		fmt.Printf("Identity is valid; branding approval declaration: %s (not verified by Foundry).\n", c.Branding.Approval.Status)
		return 0
	case "sync":
		p, err := workspace.At(*project)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		c, assets, err := identity.Resolve("", filepath.Join(p.Root, filepath.FromSlash(identity.Source)), "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		files, err := identity.Render(a.FS, c, assets, p.Kind)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		paths, err := scaffold.Apply(p.Root, files, true, *dry)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if *dry {
			fmt.Println("Would update identity bindings (no files changed):")
		} else {
			fmt.Println("Updated identity bindings:")
		}
		for _, path := range paths {
			fmt.Println("  +", path)
		}
		fmt.Println("Run the native formatter and quality gate. Ownership references are excluded from public bindings.")
		return 0
	}
	return 2
}
