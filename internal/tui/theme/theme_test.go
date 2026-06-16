package theme

import "testing"

func TestNamedThemes(t *testing.T) {
	for _, name := range []string{"arch", "mono", "dark", "light", "high-contrast"} {
		t.Run(name, func(t *testing.T) {
			got, ok := Named(name)
			if !ok {
				t.Fatalf("Named(%q) missing", name)
			}
			if got.Name != name {
				t.Fatalf("Name = %q, want %q", got.Name, name)
			}
		})
	}
}

func TestDetectRejectsUnknownTheme(t *testing.T) {
	if _, err := Detect("missing"); err == nil {
		t.Fatal("Detect() error = nil, want unknown theme error")
	}
}

func TestMonoThemeDoesNotEmitANSI(t *testing.T) {
	tm, ok := Named("mono")
	if !ok {
		t.Fatal("mono theme missing")
	}
	if got := tm.Accent("AUR"); got != "AUR" {
		t.Fatalf("Accent() = %q, want plain text", got)
	}
}
