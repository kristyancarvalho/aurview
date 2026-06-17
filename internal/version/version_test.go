package version

import (
	"strings"
	"testing"
)

func TestInfoStringIncludesMetadata(t *testing.T) {
	info := Info{Version: "0.4.2", Commit: "abc1234", Date: "2026-06-17T05:00:00Z"}
	got := info.String("aurview")

	for _, want := range []string{"aurview 0.4.2", "commit: abc1234", "date: 2026-06-17T05:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Info.String() missing %q in:\n%s", want, got)
		}
	}
}

func TestInfoStringUsesFallbacks(t *testing.T) {
	got := (Info{}).String("aurview")

	for _, want := range []string{"aurview dev", "commit: none", "date: unknown"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Info.String() missing fallback %q in:\n%s", want, got)
		}
	}
}
