package tui

import (
	"testing"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/filter"
	"github.com/kristyancarvalho/aurview/internal/ranking"
)

func TestModelFiltersUseAllResultsWithoutReranking(t *testing.T) {
	maint := "alice"
	model := New(Options{})
	model.lastQuery = "pa"
	model.allResults = []ranking.RankedPackage{
		{Package: aur.Package{Source: "aur", Name: "paru", Maintainer: &maint, NumVotes: 50, Popularity: 4}, Score: 20},
		{Package: aur.Package{Source: "aur", Name: "orphan", NumVotes: 50, Popularity: 4}, Score: 19},
		{Package: aur.Package{Source: "custom", Name: "pacseek", Maintainer: &maint, NumVotes: 50, Popularity: 4}, Score: 18},
	}
	model.filterState = filter.State{Source: "aur", Maintainer: filter.MaintainerMaintained, Match: filter.MatchPrefix}

	model.applyFilters()

	if len(model.results) != 1 || model.results[0].Package.Name != "paru" {
		t.Fatalf("filtered results = %#v, want only paru", model.results)
	}
	if len(model.allResults) != 3 {
		t.Fatalf("allResults changed, got %d entries", len(model.allResults))
	}
}

func TestCycleSelectedFilterUpdatesResults(t *testing.T) {
	model := New(Options{})
	model.lastQuery = "pkg"
	model.allResults = []ranking.RankedPackage{
		{Package: aur.Package{Source: "aur", Name: "pkg-a"}, Score: 2},
		{Package: aur.Package{Source: "custom", Name: "pkg-b"}, Score: 1},
	}
	model.focus = focusFilters
	model.filterIndex = 0

	_ = model.cycleSelectedFilter()

	if model.filterState.Source != "aur" {
		t.Fatalf("source filter = %q, want aur", model.filterState.Source)
	}
	if len(model.results) != 1 || model.results[0].Package.Name != "pkg-a" {
		t.Fatalf("filtered results = %#v, want only pkg-a", model.results)
	}
}

func TestDeveloperQueryStatus(t *testing.T) {
	model := New(Options{})
	model.results = []ranking.RankedPackage{{Package: aur.Package{Name: "pkg"}}}

	model.setFilteredStatus("dev:alice")

	if got, want := model.status, `1 packages matched developer "alice"`; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}
