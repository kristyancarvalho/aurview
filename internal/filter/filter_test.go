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

func TestActiveCount(t *testing.T) {
	state := State{Source: "aur", MinVotes: 10, Match: MatchExact}
	if got := state.ActiveCount(); got != 3 {
		t.Fatalf("ActiveCount() = %d, want 3", got)
	}
}
