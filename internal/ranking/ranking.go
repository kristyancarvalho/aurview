package ranking

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/kristyancarvalho/aurview/internal/aur"
)

type RankedPackage struct {
	Package aur.Package
	Score   float64
}

type Scorer struct {
	Now time.Time
}

func NewScorer(now time.Time) Scorer {
	if now.IsZero() {
		now = time.Now()
	}
	return Scorer{Now: now}
}

func (s Scorer) Rank(query string, pkgs []aur.Package) []RankedPackage {
	out := make([]RankedPackage, 0, len(pkgs))
	for _, pkg := range pkgs {
		out = append(out, RankedPackage{Package: pkg.Clone(), Score: s.Score(query, pkg)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Package.Name < out[j].Package.Name
		}
		return out[i].Score > out[j].Score
	})
	return out
}

func (s Scorer) Score(query string, pkg aur.Package) float64 {
	q := normalize(query)
	name := normalize(pkg.Name)
	desc := normalize(pkg.Description)

	score := 0.0
	switch {
	case q != "" && name == q:
		score += 70
	case q != "" && strings.HasPrefix(name, q):
		score += 48
	case q != "" && strings.Contains(name, q):
		score += 30
	}
	if q != "" && strings.Contains(desc, q) {
		score += 13
	}

	score += clamp(math.Log1p(float64(pkg.NumVotes))*4.5, 0, 28)
	score += clamp(math.Log1p(pkg.Popularity)*8, 0, 24)
	score += recencyScore(s.Now, pkg.LastModified)

	if !pkg.IsOutOfDate() {
		score += 8
	} else {
		score -= 14
	}
	if !pkg.IsOrphan() {
		score += 8
	} else {
		score -= 10
	}

	return math.Round(score*100) / 100
}

func recencyScore(now time.Time, unix int64) float64 {
	if unix <= 0 {
		return 0
	}
	modified := time.Unix(unix, 0)
	ageDays := now.Sub(modified).Hours() / 24
	switch {
	case ageDays < 0:
		return 9
	case ageDays <= 30:
		return 12
	case ageDays <= 180:
		return 9
	case ageDays <= 365:
		return 6
	case ageDays <= 365*3:
		return 3
	default:
		return 0
	}
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
