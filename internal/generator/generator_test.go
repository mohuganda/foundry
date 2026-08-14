package generator

import "testing"

func TestSanitize(t *testing.T) {
	cases := map[string]string{
		"Users Service": "users_service",
		"mediguide-app": "mediguide_app",
		"simple":        "simple",
	}
	for in, want := range cases {
		if got := sanitize(in); got != want {
			t.Fatalf("sanitize(%q)=%q want %q", in, got, want)
		}
	}
}
