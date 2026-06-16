package ranking

import (
	"testing"
	"time"

	"github.com/kristyancarvalho/aurview/internal/aur"
)

func TestScorerRanksBalancedRelevance(t *testing.T) {
	maintainer := "alice"
	now := time.Unix(1_700_000_000, 0)
	outdated := now.Add(-24 * time.Hour).Unix()
	scorer := NewScorer(now)

	pkgs := []aur.Package{
		{Name: "libparu", Description: "library", NumVotes: 20, Popularity: 1, LastModified: now.Add(-10 * 24 * time.Hour).Unix(), Maintainer: &maintainer},
		{Name: "paru", Description: "AUR helper", NumVotes: 5, Popularity: 0.2, LastModified: now.Add(-900 * 24 * time.Hour).Unix(), Maintainer: &maintainer},
		{Name: "x", Description: "paru integration", NumVotes: 2000, Popularity: 50, LastModified: now.Add(-3 * 24 * time.Hour).Unix(), Maintainer: &maintainer, OutOfDate: &outdated},
	}

	ranked := scorer.Rank("paru", pkgs)
	if ranked[0].Package.Name != "paru" {
		t.Fatalf("top ranked = %s, want exact name match", ranked[0].Package.Name)
	}
	if ranked[2].Package.Name != "x" {
		t.Fatalf("last ranked = %s, want description-only out-of-date match", ranked[2].Package.Name)
	}
}

func TestScoreSignals(t *testing.T) {
	maintainer := "alice"
	now := time.Unix(1_700_000_000, 0)
	scorer := NewScorer(now)
	base := aur.Package{
		Name:         "paru",
		Description:  "AUR helper",
		NumVotes:     100,
		Popularity:   5,
		LastModified: now.Add(-7 * 24 * time.Hour).Unix(),
		Maintainer:   &maintainer,
	}

	tests := []struct {
		name   string
		mutate func(*aur.Package)
		lower  bool
	}{
		{name: "orphaned package loses score", mutate: func(p *aur.Package) { p.Maintainer = nil }, lower: true},
		{name: "outdated package loses score", mutate: func(p *aur.Package) { v := now.Unix(); p.OutOfDate = &v }, lower: true},
		{name: "old package loses recency score", mutate: func(p *aur.Package) { p.LastModified = now.Add(-5 * 365 * 24 * time.Hour).Unix() }, lower: true},
	}

	baseScore := scorer.Score("paru", base)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := base
			tt.mutate(&pkg)
			got := scorer.Score("paru", pkg)
			if tt.lower && got >= baseScore {
				t.Fatalf("score = %v, want lower than %v", got, baseScore)
			}
		})
	}
}
