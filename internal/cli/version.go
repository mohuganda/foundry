package cli

import (
	"runtime/debug"
	"strings"
)

// Version is optionally set by release builds using -ldflags -X.
// Source installations otherwise use the version recorded by the Go toolchain.
var Version string

func CurrentVersion() string {
	info, _ := debug.ReadBuildInfo()
	return buildVersion(Version, info)
}

func buildVersion(override string, info *debug.BuildInfo) string {
	if override != "" {
		return strings.TrimPrefix(override, "v")
	}
	if info == nil {
		return "dev"
	}
	if info.Main.Replace == nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	revision, dirty := "", false
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}
	version := "dev"
	if revision != "" {
		if len(revision) > 12 {
			revision = revision[:12]
		}
		version += "-" + revision
	}
	if dirty {
		version += "+dirty"
	}
	return version
}
