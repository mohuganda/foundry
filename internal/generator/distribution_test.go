package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestHomebrewFormulaReleaseChecksums(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX release helper")
	}
	script, err := filepath.Abs("../../scripts/homebrew-formula.sh")
	if err != nil {
		t.Fatal(err)
	}
	var entries []string
	for i, platform := range []string{"darwin_arm64", "darwin_amd64", "linux_arm64", "linux_amd64"} {
		entries = append(entries, strings.Repeat(string(rune('a'+i)), 64)+"  foundry_1.3.0_"+platform+".tar.gz\n")
	}
	manifest := strings.Join(entries, "")
	for _, tc := range []struct {
		name, version, manifest string
		wantFailure             bool
	}{
		{"stable", "1.3.0", manifest, false},
		{"tag", "v1.3.0", manifest, false},
		{"missing", "1.3.0", strings.Join(entries[:3], ""), true},
		{"duplicate", "1.3.0", manifest + entries[0], true},
		{"invalid digest", "1.3.0", strings.Replace(manifest, "aaaa", "zzzz", 1), true},
		{"short digest", "1.3.0", strings.Replace(manifest, "aaaa", "a", 1), true},
		{"wrong release", "1.3.1", manifest, true},
		{"prerelease", "1.3.0-rc.1", manifest, true},
		{"invalid version", "1.3.0\"", manifest, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			checksums := filepath.Join(dir, "SHA256SUMS")
			if err := os.WriteFile(checksums, []byte(tc.manifest), 0o600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("sh", script, tc.version, checksums)
			cmd.Dir = dir
			output, err := cmd.Output()
			if tc.wantFailure {
				if err == nil || len(output) != 0 {
					t.Fatalf("invalid input emitted a formula: %s; error: %v", output, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			formula := string(output)
			for _, entry := range entries {
				fields := strings.Fields(entry)
				if !strings.Contains(formula, fields[0]) || !strings.Contains(formula, "/v1.3.0/"+fields[1]) {
					t.Fatalf("formula lost the release URL or checksum for %s", fields[1])
				}
			}
		})
	}
}
