package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jabahum/foundry/internal/authconfig"
	"github.com/jabahum/foundry/internal/scaffold"
	"github.com/jabahum/foundry/internal/workspace"
)

func (a App) authCmd(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: foundry auth example\n       foundry auth validate --file auth.json --template go|react|vanilla-js|flutter\n       foundry auth configure --project path [--file auth.json] [--dry-run]")
		return 0
	}
	f := flag.NewFlagSet("auth "+args[0], flag.ContinueOnError)
	file := f.String("file", "", "non-secret authentication JSON")
	project := f.String("project", ".", "generated root")
	kind := f.String("template", "react", "validation target: go, react, vanilla-js, flutter")
	dry := f.Bool("dry-run", false, "preview writes")
	if e := f.Parse(args[1:]); e != nil {
		if e == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if f.NArg() != 0 {
		return 2
	}
	fail := func(e error) int { fmt.Fprintln(os.Stderr, e); return 1 }
	switch args[0] {
	case "example":
		c := authconfig.Default()
		c.Mode = "oidc"
		c.Issuer = "https://auth.example.org/realms/moh"
		c.Audience = "programme-api"
		c.Web = authconfig.Client{ClientID: "programme-web", RedirectURI: "https://portal.example.org/auth/callback", PostLogoutRedirectURI: "https://portal.example.org/auth/signed-out"}
		c.Native = authconfig.Client{ClientID: "programme-mobile", RedirectURI: "ug.go.health.programme://auth/callback", PostLogoutRedirectURI: "ug.go.health.programme://auth/signed-out"}
		fmt.Print(string(authconfig.JSON(c)))
		return 0
	case "validate":
		if *file == "" {
			return fail(fmt.Errorf("--file is required"))
		}
		k := *kind
		if k == "vanilla-js" {
			k = "javascript"
		}
		_, e := authconfig.Resolve("", *file, k)
		if e != nil {
			return fail(e)
		}
		fmt.Println("Authentication configuration is valid (provider connectivity not checked).")
		return 0
	case "configure":
		p, e := workspace.At(*project)
		if e != nil {
			return fail(e)
		}
		// Require integration installed by this Foundry version; do not silently retrofit app code.
		marker := map[string]string{"go": "authconfig/config.go", "react": "src/shared/auth/client.js", "javascript": "src/shared/auth/client.js", "flutter": "packages/foundation_auth/lib/src/oidc_auth.dart"}[p.Kind]
		if marker == "" {
			return fail(fmt.Errorf("authentication is not implemented for %s", p.Kind))
		}
		if _, e = os.Stat(filepath.Join(p.Root, marker)); e != nil {
			return fail(fmt.Errorf("auth integration missing; generate a fresh foundation or follow the migration guide"))
		}
		input := *file
		if input == "" {
			input = filepath.Join(p.Root, "auth/auth.json")
		}
		c, e := authconfig.Resolve("", input, p.Kind)
		if e != nil {
			return fail(e)
		}
		files, e := authconfig.Render(a.FS, c, p.Kind)
		if e != nil {
			return fail(e)
		}
		paths, e := scaffold.Apply(p.Root, files, true, *dry)
		if e != nil {
			return fail(e)
		}
		for _, p := range paths {
			fmt.Println("  +", p)
		}
		if *dry {
			fmt.Println("Preview only; no files changed.")
		} else {
			fmt.Println("Authentication configuration updated. Format, test and rebuild the application.")
		}
		return 0
	default:
		return fail(fmt.Errorf("unknown auth command; use foundry auth --help"))
	}
}
