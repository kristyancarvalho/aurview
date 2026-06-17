package filter

import (
	"sort"
	"strings"
	"time"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/ranking"
)

type MaintainerMode int

const (
	MaintainerAny MaintainerMode = iota
	MaintainerMaintained
	MaintainerOrphaned
)

type FlagMode int

const (
	FlagAny FlagMode = iota
	FlagClean
	FlagOutOfDate
)

type MatchMode int

const (
	MatchSmart MatchMode = iota
	MatchName
	MatchPrefix
	MatchExact
)

type State struct {
	Source        string
	Maintainer    MaintainerMode
	Flag          FlagMode
	MinVotes      int
	MinPopularity float64
	RecentDays    int
	Match         MatchMode
}

func (s State) Apply(query string, results []ranking.RankedPackage, now time.Time) []ranking.RankedPackage {
	if now.IsZero() {
		now = time.Now()
	}
	out := make([]ranking.RankedPackage, 0, len(results))
	for _, result := range results {
		if s.MatchPackage(query, result.Package, now) {
			out = append(out, ranking.RankedPackage{
				Package: result.Package.Clone(),
				Score:   result.Score,
			})
		}
	}
	return out
}

func (s State) MatchPackage(query string, pkg aur.Package, now time.Time) bool {
	if s.Source != "" && !strings.EqualFold(sourceKey(pkg), s.Source) {
		return false
	}
	switch s.Maintainer {
	case MaintainerMaintained:
		if pkg.IsOrphan() {
			return false
		}
	case MaintainerOrphaned:
		if !pkg.IsOrphan() {
			return false
		}
	}
	switch s.Flag {
	case FlagClean:
		if pkg.IsOutOfDate() {
			return false
		}
	case FlagOutOfDate:
		if !pkg.IsOutOfDate() {
			return false
		}
	}
	if pkg.NumVotes < s.MinVotes {
		return false
	}
	if pkg.Popularity < s.MinPopularity {
		return false
	}
	if s.RecentDays > 0 {
		if pkg.LastModified <= 0 {
			return false
		}
		cutoff := now.AddDate(0, 0, -s.RecentDays)
		if time.Unix(pkg.LastModified, 0).Before(cutoff) {
			return false
		}
	}
	return matchName(query, pkg, s.Match)
}

func (s State) Active() bool {
	return s.Source != "" ||
		s.Maintainer != MaintainerAny ||
		s.Flag != FlagAny ||
		s.MinVotes > 0 ||
		s.MinPopularity > 0 ||
		s.RecentDays > 0 ||
		s.Match != MatchSmart
}

func (s State) ActiveCount() int {
	count := 0
	if s.Source != "" {
		count++
	}
	if s.Maintainer != MaintainerAny {
		count++
	}
	if s.Flag != FlagAny {
		count++
	}
	if s.MinVotes > 0 {
		count++
	}
	if s.MinPopularity > 0 {
		count++
	}
	if s.RecentDays > 0 {
		count++
	}
	if s.Match != MatchSmart {
		count++
	}
	return count
}

func Sources(results []ranking.RankedPackage) []string {
	seen := map[string]bool{}
	for _, result := range results {
		key := sourceKey(result.Package)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
	}
	out := make([]string, 0, len(seen))
	for source := range seen {
		out = append(out, source)
	}
	sort.Strings(out)
	return out
}

func SourceLabel(source string) string {
	if strings.TrimSpace(source) == "" || strings.EqualFold(source, "aur") {
		return "AUR"
	}
	return source
}

func sourceKey(pkg aur.Package) string {
	source := strings.TrimSpace(pkg.Source)
	if source == "" {
		source = "aur"
	}
	return strings.ToLower(source)
}

func matchName(query string, pkg aur.Package, mode MatchMode) bool {
	if mode == MatchSmart {
		return true
	}
	q := strings.ToLower(strings.TrimSpace(query))
	name := strings.ToLower(strings.TrimSpace(pkg.Name))
	if q == "" {
		return true
	}
	switch mode {
	case MatchName:
		return strings.Contains(name, q)
	case MatchPrefix:
		return strings.HasPrefix(name, q)
	case MatchExact:
		return name == q
	default:
		return true
	}
}
