package cli

import (
	"runtime/debug"
	"testing"
)

func TestBuildVersion(t *testing.T) {
	local := &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}, Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abcdef1234567890"}, {Key: "vcs.modified", Value: "true"}}}
	for _, tc := range []struct {
		name, override string
		info           *debug.BuildInfo
		want           string
	}{
		{"release override", "1.4.0", local, "1.4.0"},
		{"tag override", "v1.4.0", nil, "1.4.0"},
		{"make dev", "dev", local, "dev"},
		{"source install", "", &debug.BuildInfo{Main: debug.Module{Version: "v1.3.0"}}, "1.3.0"},
		{"pseudo version", "", &debug.BuildInfo{Main: debug.Module{Version: "v1.3.1-0.20261025000000-abcdef123456"}}, "1.3.1-0.20261025000000-abcdef123456"},
		{"local dirty checkout", "", local, "dev-abcdef123456+dirty"},
		{"replaced module", "", &debug.BuildInfo{Main: debug.Module{Version: "v1.3.0", Replace: &debug.Module{Path: "../local"}}}, "dev"},
		{"missing metadata", "", nil, "dev"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildVersion(tc.override, tc.info); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
