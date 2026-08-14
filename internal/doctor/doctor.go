package doctor

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Check struct {
	Name, Detail string
	OK           bool
}

func Run() []Check {
	checks := []Check{{Name: "OS", Detail: runtime.GOOS + "/" + runtime.GOARCH, OK: true}}
	for _, c := range []struct{ Name, Bin string }{
		{"Git", "git"}, {"Go", "go"}, {"Docker", "docker"}, {"Buf", "buf"},
		{"grpcurl", "grpcurl"}, {"Flutter", "flutter"}, {"Dart", "dart"}, {"Melos", "melos"}, {"Mason", "mason"},
	} {
		p, err := exec.LookPath(c.Bin)
		checks = append(checks, Check{Name: c.Name, Detail: p, OK: err == nil})
	}
	return checks
}

func Print() {
	for _, c := range Run() {
		mark := "✓"
		detail := c.Detail
		if !c.OK {
			mark = "-"
			detail = "not found"
		}
		fmt.Printf("%s %-10s %s\n", mark, c.Name, detail)
	}
}
