package upgrade

import (
	"fmt"

	"github.com/jabahum/foundry/internal/config"
	"github.com/jabahum/foundry/internal/registry"
)

type Report struct {
	Current, Latest, Template, Source string
	UpdateAvailable                   bool
}

func Check(root string) (Report, error) {
	c, err := config.Read(root)
	if err != nil {
		return Report{}, err
	}
	builtin, err := registry.Find(c.Template)
	if err != nil {
		return Report{}, err
	}
	latest := builtin.Version
	source := "bundled"
	return Report{Current: c.TemplateVersion, Latest: latest, Template: builtin.ID, Source: source, UpdateAvailable: newer(latest, c.TemplateVersion)}, nil
}

func Format(r Report) string {
	status := "up to date"
	if r.UpdateAvailable {
		status = "update available"
	}
	return fmt.Sprintf("Template: %s\nCurrent: %s\nLatest: %s (%s)\nStatus: %s\n", r.Template, r.Current, r.Latest, r.Source, status)
}

func newer(a, b string) bool { return compare(a, b) > 0 }
func compare(a, b string) int {
	var A, B [3]int
	fmt.Sscanf(a, "%d.%d.%d", &A[0], &A[1], &A[2])
	fmt.Sscanf(b, "%d.%d.%d", &B[0], &B[1], &B[2])
	for i := 0; i < 3; i++ {
		if A[i] > B[i] {
			return 1
		}
		if A[i] < B[i] {
			return -1
		}
	}
	return 0
}
