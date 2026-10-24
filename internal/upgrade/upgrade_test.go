package upgrade

import (
	"testing"

	"github.com/jabahum/foundry/internal/projectconfig"
)

func TestCheckUsesBundledVersion(t *testing.T) {
	for _, version := range []string{"0.9.0", "1.0.0", "2.0.0"} {
		t.Run(version, func(t *testing.T) {
			root := t.TempDir()
			if err := projectconfig.Write(root, projectconfig.Config{Template: "go-foundations", TemplateVersion: version}); err != nil {
				t.Fatal(err)
			}
			report, err := Check(root)
			if err != nil {
				t.Fatal(err)
			}
			if report.Latest != "1.0.0" || report.Source != "bundled" || report.UpdateAvailable != (version == "0.9.0") {
				t.Fatalf("unexpected upgrade report: %#v", report)
			}
		})
	}
}
