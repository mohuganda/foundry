package doctor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jabahum/foundry/internal/workspace"
)

func TestDoctorScopesToolsAndChecksNodeVersions(t *testing.T) {
	for _, tc := range []struct {
		version string
		ok      bool
	}{{"v22.17.0", false}, {"v22.18.0", true}, {"v23.0.0", false}, {"v24.1.0", true}, {"v26.0.0", true}, {"unknown", false}} {
		p := workspace.Project{Kind: "react"}
		checks := inspect(p, func(bin string, args ...string) (string, error) {
			switch bin {
			case "node":
				return tc.version, nil
			case "npm":
				return "10.9.0", nil
			default:
				t.Fatalf("unrelated tool probed: %s", bin)
				return "", nil
			}
		})
		var out bytes.Buffer
		if got := Print(&out, p, checks); got != tc.ok {
			t.Fatalf("Node %s: %s", tc.version, out.String())
		}
		if !tc.ok && !strings.Contains(out.String(), "Fix:") {
			t.Fatal("missing remediation")
		}
	}
}

func TestDoctorUsesGoManifestAndOptionalFailures(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module acme.org/api\ngo 1.25.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"go version go1.24.0 linux/amd64", "go version go1.25.1 linux/amd64"} {
		p := workspace.Project{Root: root, Kind: "go"}
		checks := inspect(p, func(bin string, args ...string) (string, error) {
			if bin == "go" {
				return version, nil
			}
			return "", errors.New("not installed")
		})
		var out bytes.Buffer
		want := strings.Contains(version, "1.25.1")
		if got := Print(&out, p, checks); got != want {
			t.Fatal(out.String())
		}
	}
}

func TestFlutterDoctorReadsPinWithoutLaunchingSDK(t *testing.T) {
	root := t.TempDir()
	cache := filepath.Join(root, ".fvm", "flutter_sdk", "bin", "cache")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	for file, data := range map[string]string{filepath.Join(root, ".fvmrc"): `{"flutter":"3.44.8"}`, filepath.Join(cache, "flutter.version.json"): `{"frameworkVersion":"3.44.8"}`} {
		if err := os.WriteFile(file, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	p := workspace.Project{Root: root, Kind: "flutter"}
	run := func(bin string, args ...string) (string, error) {
		if bin != "fvm" || strings.Join(args, " ") != "--version" {
			t.Fatal("doctor must not launch SDK/install commands")
		}
		return "4.0.0", nil
	}
	var out bytes.Buffer
	if !Print(&out, p, inspect(p, run)) {
		t.Fatal(out.String())
	}
	if err := os.Remove(filepath.Join(cache, "flutter.version.json")); err != nil {
		t.Fatal(err)
	}
	if Print(&out, p, inspect(p, run)) {
		t.Fatal("missing SDK was accepted")
	}
}

func TestProbeUsesWorkspaceAndDisablesToolchainDownloads(t *testing.T) {
	t.Setenv("FOUNDRY_DOCTOR_TEST_HELPER", "1")
	root := t.TempDir()
	out, err := command(root, os.Args[0], "-test.run=^TestDoctorProbeProcess$")
	if err != nil {
		t.Fatalf("probe: %q, %v", out, err)
	}
	got, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(got, want) {
		t.Fatalf("probe ran in %q instead of %q", out, root)
	}

}

func TestDoctorProbeProcess(t *testing.T) {
	if os.Getenv("FOUNDRY_DOCTOR_TEST_HELPER") != "1" {
		return
	}
	if os.Getenv("RUSTUP_AUTO_INSTALL") != "0" || os.Getenv("GOTOOLCHAIN") != "local" {
		os.Exit(2)
	}
	dir, err := os.Getwd()
	if err != nil {
		os.Exit(3)
	}
	fmt.Print(dir)
	os.Exit(0)
}
