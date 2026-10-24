package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jabahum/foundry"
)

func copyEmbeddedScript(t *testing.T, template, script string) (string, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX helper test; native PowerShell gate is covered by CI")
	}
	root := t.TempDir()
	path := filepath.Join(root, "scripts", script)
	data, err := foundry.Templates.ReadFile("templates/" + template + "/scripts/" + script)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o755); err != nil {
		t.Fatal(err)
	}
	return root, path
}

func fakeFVM(t *testing.T) (string, string) {
	t.Helper()
	bin := t.TempDir()
	log := filepath.Join(t.TempDir(), "calls.log")
	data := "#!/bin/sh\nprintf '%s|%s\\n' \"$PWD\" \"$*\" >> \"$FOUNDRY_COMMAND_LOG\"\nif [ \"$*\" = \"${FOUNDRY_FAIL_COMMAND:-}\" ]; then exit 23; fi\n"
	if err := os.WriteFile(filepath.Join(bin, "fvm"), []byte(data), 0o755); err != nil {
		t.Fatal(err)
	}
	return bin, log
}

func TestFlutterSetupAnchorsItsRootAndStopsOnFailure(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[fail], func(t *testing.T) {
			root, script := copyEmbeddedScript(t, "flutter_foundations", "setup.sh")
			bin, log := fakeFVM(t)
			cmd := exec.Command("bash", script)
			cmd.Dir = t.TempDir()
			cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FOUNDRY_COMMAND_LOG="+log)
			if fail {
				cmd.Env = append(cmd.Env, "FOUNDRY_FAIL_COMMAND=dart pub global activate mason_cli")
			}
			output, err := cmd.CombinedOutput()
			if fail && err == nil {
				t.Fatal("setup continued after installer failure")
			}
			if !fail && err != nil {
				t.Fatalf("setup: %v\n%s", err, output)
			}
			data, err := os.ReadFile(log)
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			want := 4
			if fail {
				want = 2
			}
			if len(lines) != want {
				t.Fatalf("unexpected setup calls: %s", data)
			}
			resolved, err := filepath.EvalSymlinks(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, line := range lines {
				parts := strings.SplitN(line, "|", 2)
				actual, err := filepath.EvalSymlinks(parts[0])
				if err != nil {
					t.Fatal(err)
				}
				if actual != resolved {
					t.Fatalf("setup used the caller directory: %s", line)
				}
			}
		})
	}
}

func TestFlutterRunEnvironmentUsesTheAppAndValidatesArguments(t *testing.T) {
	root, script := copyEmbeddedScript(t, "flutter_foundations", "run_flavor.sh")
	for _, dir := range []string{"apps/demo", "env"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "env/development.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	bin, log := fakeFVM(t)
	for _, args := range [][]string{{"demo"}, {"demo", "unknown"}, {"../outside"}, {"missing"}} {
		cmd := exec.Command("bash", append([]string{script}, args...)...)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FOUNDRY_COMMAND_LOG="+log)
		output, err := cmd.CombinedOutput()
		if args[0] == "demo" && len(args) == 1 {
			if err != nil {
				t.Fatalf("run: %v\n%s", err, output)
			}
		} else if err == nil {
			t.Fatalf("accepted invalid args: %v", args)
		}
	}
	data, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	app, err := filepath.EvalSymlinks(filepath.Join(root, "apps/demo"))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(strings.TrimSpace(string(data)), "|", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected FVM invocation: %s", data)
	}
	actual, err := filepath.EvalSymlinks(parts[0])
	if err != nil {
		t.Fatal(err)
	}
	if actual != app || parts[1] != "flutter run --dart-define-from-file=../../env/development.json" {
		t.Fatalf("unexpected FVM invocation: %s", data)
	}
}

func TestDevelopmentKeysAreProjectLocalAndPreserved(t *testing.T) {
	root, script := copyEmbeddedScript(t, "go-foundations", "generate-dev-keys.sh")
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("OpenSSL not installed")
	}
	caller := t.TempDir()
	run := func() {
		cmd := exec.Command("sh", script)
		cmd.Dir = caller
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("keys: %v\n%s", err, output)
		}
	}
	run()
	private := filepath.Join(root, "secrets/jwt_private.pem")
	before, err := os.ReadFile(private)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(private)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("private mode %o", info.Mode().Perm())
	}
	run()
	after, err := os.ReadFile(private)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("existing private key was replaced")
	}
	if _, err := os.Stat(filepath.Join(caller, "secrets")); !os.IsNotExist(err) {
		t.Fatal("keys were written in the caller directory")
	}
}
