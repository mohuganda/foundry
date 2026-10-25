package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jabahum/foundry/internal/workspace"
)

type Check struct {
	Name, Detail, Fix string
	OK, Required      bool
}

type probe func(string, ...string) (string, error)

func command(dir, bin string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "RUSTUP_AUTO_INSTALL=0", "GOTOOLCHAIN=local", "CI=true")
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return strings.TrimSpace(string(out)), err
}

// Inspect probes only the tools relevant to this workspace. No installs are run.
func Inspect(project workspace.Project) []Check {
	return inspect(project, func(bin string, args ...string) (string, error) { return command(project.Root, bin, args...) })
}

func inspect(p workspace.Project, run probe) []Check {
	var checks []Check
	tool := func(name, bin, minimum, fix string, required bool, args ...string) {
		out, err := run(bin, args...)
		ok := err == nil
		detail := strings.SplitN(out, "\n", 2)[0]
		if err != nil {
			detail = "unavailable: " + err.Error()
		} else if minimum != "" {
			ok = atLeast(out, minimum)
			detail += " (requires >= " + minimum + ")"
		}
		checks = append(checks, Check{name, detail, fix, ok, required})
	}
	switch p.Kind {
	case "go", "foundry":
		minimum := manifestVersion(filepath.Join(p.Root, "go.mod"), `(?m)^go\s+(\d+\.\d+(?:\.\d+)?)`, "1.23")
		tool("Go", "go", minimum, "Install/update Go, then run go mod download at the project root.", true, "version")
		if p.Kind == "go" {
			tool("Buf", "buf", "", "Run make tools and add GOBIN (or GOPATH/bin) to PATH for code generation.", false, "--version")
			tool("Docker", "docker", "", "Install/start Docker with Compose if using make docker-up.", false, "version", "--format", "{{.Server.Version}}")
		}
	case "javascript", "react":
		tool("Node.js", "node", "", "Install Node 22.18+, 24, or a supported newer even-numbered release; select it with nvm use.", true, "--version")
		c := &checks[len(checks)-1]
		if c.OK {
			v := numbers(c.Detail)
			c.OK = v[0] == 22 && atLeast(c.Detail, "22.18.0") || v[0] >= 24 && v[0]%2 == 0
			c.Detail += " (requires 22.18+ on 22.x, or an even release >= 24)"
		}
		tool("npm", "npm", "", "Install npm with Node, then run npm ci at the project root.", true, "--version")
	case "rust":
		minimum := manifestVersion(filepath.Join(p.Root, "Cargo.toml"), `(?m)^rust-version\s*=\s*"([0-9.]+)"`, "1.85")
		tool("Rust", "rustc", minimum, "Install/update Rust using rustup; the project selects its toolchain via rust-toolchain.toml.", true, "--version")
		tool("Cargo", "cargo", "", "Install Cargo with Rust using rustup.", true, "--version")
		tool("Clippy", "cargo", "", "Run rustup component add clippy for the selected toolchain.", true, "clippy", "--version")
		tool("rustfmt", "rustfmt", "", "Run rustup component add rustfmt for the selected toolchain.", true, "--version")
	case "flutter":
		tool("FVM", "fvm", "", "Install FVM, then run bash scripts/setup.sh at the workspace root.", true, "--version")
		// Read the selected SDK's cached version instead of invoking FVM/Flutter,
		// which can install or initialize SDKs when called for the first time.
		var pin struct {
			Flutter string `json:"flutter"`
		}
		var sdk struct {
			Version string `json:"frameworkVersion"`
		}
		pinErr := readJSON(filepath.Join(p.Root, ".fvmrc"), &pin)
		sdkErr := readJSON(filepath.Join(p.Root, ".fvm", "flutter_sdk", "bin", "cache", "flutter.version.json"), &sdk)
		ok := pinErr == nil && sdkErr == nil && pin.Flutter != "" && sdk.Version == pin.Flutter
		detail := "selected SDK unavailable or version differs from .fvmrc"
		if ok {
			detail = sdk.Version + " (matches .fvmrc)"
		}
		checks = append(checks, Check{"Flutter SDK", detail, "Run fvm install and bash scripts/setup.sh; doctor does not install or initialize SDKs.", ok, true})
	default:
		for _, bin := range []string{"git", "go", "node", "npm", "rustc", "cargo", "docker", "fvm"} {
			args := []string{"--version"}
			if bin == "go" {
				args = []string{"version"}
			}
			tool(bin, bin, "", "Install this tool only if needed for your chosen template.", false, args...)
		}
	}
	return checks
}

func Print(w io.Writer, p workspace.Project, checks []Check) bool {
	if p.Kind == "" {
		fmt.Fprintln(w, "No workspace detected; showing optional tool availability. Use --project <root> to select one.")
	} else {
		fmt.Fprintf(w, "Workspace: %s (%s)\n", p.Root, p.Kind)
	}
	ok := true
	for _, c := range checks {
		mark, scope := "OK", "required"
		if !c.Required {
			scope = "optional"
		}
		if !c.OK {
			mark = "WARN"
			if c.Required {
				mark, ok = "FAIL", false
			}
		}
		fmt.Fprintf(w, "%s %-12s [%s] %s\n", mark, c.Name, scope, c.Detail)
		if !c.OK {
			fmt.Fprintln(w, "  Fix:", c.Fix)
		}
	}
	return ok
}

var numericVersion = regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

func numbers(s string) [3]int {
	var v [3]int
	for i, part := range strings.Split(numericVersion.FindString(s), ".") {
		if i < len(v) {
			v[i], _ = strconv.Atoi(part)
		}
	}
	return v
}

func atLeast(actual, minimum string) bool {
	if numericVersion.FindString(actual) == "" {
		return false
	}
	a, b := numbers(actual), numbers(minimum)
	for i := range a {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return true
}

func manifestVersion(file, pattern, fallback string) string {
	data, _ := os.ReadFile(file)
	match := regexp.MustCompile(pattern).FindSubmatch(data)
	if len(match) == 2 {
		return string(match[1])
	}
	return fallback
}

func readJSON(file string, value any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
