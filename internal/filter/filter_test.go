package filter

import (
	"testing"
	"time"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/ranking"
)

func TestStateApplyFiltersPackageHealthAndMetrics(t *testing.T) {
	now := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	maint := "alice"
	ood := now.AddDate(0, 0, -2).Unix()
	results := []ranking.RankedPackage{
		{Package: aur.Package{Source: "aur", Name: "paru", Maintainer: &maint, NumVotes: 120, Popularity: 8, LastModified: now.AddDate(0, 0, -3).Unix()}, Score: 10},
		{Package: aur.Package{Source: "aur", Name: "old", Maintainer: &maint, NumVotes: 200, Popularity: 12, LastModified: now.AddDate(0, 0, -700).Unix()}, Score: 9},
		{Package: aur.Package{Source: "custom", Name: "orphan", NumVotes: 300, Popularity: 20, LastModified: now.Unix()}, Score: 8},
		{Package: aur.Package{Source: "aur", Name: "flagged", Maintainer: &maint, OutOfDate: &ood, NumVotes: 300, Popularity: 20, LastModified: now.Unix()}, Score: 7},
	}

	state := State{
		Source:        "aur",
		Maintainer:    MaintainerMaintained,
		Flag:          FlagClean,
		MinVotes:      100,
		MinPopularity: 5,
		RecentDays:    30,
		Match:         MatchPrefix,
	}
	got := state.Apply("pa", results, now)

	if len(got) != 1 || got[0].Package.Name != "paru" {
		t.Fatalf("Apply() = %#v, want only paru", got)
	}
}

func TestParseQueryDeveloperAliases(t *testing.T) {
	tests := map[string]struct {
		text string
		dev  string
	}{
		"dev:alice paru":             {text: "paru", dev: "alice"},
		"developer:alice paru":       {text: "paru", dev: "alice"},
		"maint:alice paru":           {text: "paru", dev: "alice"},
		"maintainer:alice paru":      {text: "paru", dev: "alice"},
		"paru dev:alice helper":      {text: "paru helper", dev: "alice"},
		"developer:Arch maint:Linux": {text: "", dev: "Arch,Linux"},
	}
	for query, want := range tests {
		got := ParseQuery(query)
		if got.Text != want.text || got.DeveloperLabel() != want.dev {
			t.Fatalf("ParseQuery(%q) = %#v, want text %q developer %q", query, got, want.text, want.dev)
		}
	}
}

func TestDeveloperSearchMatchesAURMaintainer(t *testing.T) {
	maint := "Alice"
	pkg := aur.Package{Source: "aur", Name: "paru", Maintainer: &maint}

	if !((State{}).MatchPackage("dev:ali", pkg, time.Time{})) {
		t.Fatal("dev:ali did not match AUR maintainer Alice")
	}
}

func TestDeveloperSearchMatchesLocalPackager(t *testing.T) {
	packager := "Arch Linux"
	pkg := aur.Package{Source: "core", Name: "pacman", Maintainer: &packager}

	if !((State{}).MatchPackage("developer:arch", pkg, time.Time{})) {
		t.Fatal("developer:arch did not match local packager Arch Linux")
	}
}

func TestDeveloperSearchIsCaseInsensitiveAndPartial(t *testing.T) {
	maint := "Alice Maintainer"
	pkg := aur.Package{Name: "tool", Maintainer: &maint}

	for _, query := range []string{"dev:alice", "dev:MAIN", "maintainer:tain"} {
		if !((State{}).MatchPackage(query, pkg, time.Time{})) {
			t.Fatalf("%q did not match developer %q", query, maint)
		}
	}
}

func TestDeveloperSearchNoResultBehavior(t *testing.T) {
	alice := "Alice"
	bob := "Bob"
	results := []ranking.RankedPackage{
		{Package: aur.Package{Name: "pkg-a", Maintainer: &alice}},
		{Package: aur.Package{Name: "pkg-b", Maintainer: &bob}},
		{Package: aur.Package{Name: "pkg-c"}},
	}

	got := (State{}).Apply("dev:carol", results, time.Time{})
	if len(got) != 0 {
		t.Fatalf("Apply(dev:carol) = %#v, want no results", got)
	}
}

func TestSourcesNormalizesAndSorts(t *testing.T) {
	results := []ranking.RankedPackage{
		{Package: aur.Package{Source: "custom"}},
		{Package: aur.Package{Source: ""}},
		{Package: aur.Package{Source: "aur"}},
	}

	got := Sources(results)
	want := []string{"aur", "custom"}
	if len(got) != len(want) {
		t.Fatalf("Sources() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Sources() = %#v, want %#v", got, want)
		}
	}
}

func TestSourceBadgeLabelKnownRepositories(t *testing.T) {
	tests := map[string]string{
		"":            "AUR",
		"aur":         "AUR",
		"AUR":         "AUR",
		"core":        "CORE",
		"extra":       "EXT",
		"multilib":    "MULTI",
		"chaotic-aur": "CHAOTIC",
	}
	for source, want := range tests {
		if got := SourceBadgeLabel(source); got != want {
			t.Fatalf("SourceBadgeLabel(%q) = %q, want %q", source, got, want)
		}
	}
}

func TestSourceBadgeLabelFallbacks(t *testing.T) {
	tests := map[string]string{
		"custom":         "CUSTOM",
		"local-testing":  "LT",
		"archlinuxcn":    "ARCHLINU",
		"repo.with.dots": "RWD",
	}
	for source, want := range tests {
		if got := SourceBadgeLabel(source); got != want {
			t.Fatalf("SourceBadgeLabel(%q) = %q, want %q", source, got, want)
		}
	}
}

func TestActiveCount(t *testing.T) {
	state := State{Source: "aur", MinVotes: 10, Match: MatchExact}
	if got := state.ActiveCount(); got != 3 {
		t.Fatalf("ActiveCount() = %d, want 3", got)
	}
}
