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

type Query struct {
	Text           string
	DeveloperTerms []string
}

func ParseQuery(query string) Query {
	fields := strings.Fields(query)
	text := make([]string, 0, len(fields))
	developers := []string{}
	for _, field := range fields {
		key, value, ok := strings.Cut(field, ":")
		if !ok || !isDeveloperKey(key) {
			text = append(text, field)
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		developers = append(developers, value)
	}
	return Query{
		Text:           strings.Join(text, " "),
		DeveloperTerms: developers,
	}
}

func (q Query) HasDeveloper() bool {
	return len(q.DeveloperTerms) > 0
}

func (q Query) DeveloperLabel() string {
	return strings.Join(q.DeveloperTerms, ",")
}

func (q Query) SourceSearchText() string {
	if strings.TrimSpace(q.Text) != "" {
		return q.Text
	}
	if len(q.DeveloperTerms) > 0 {
		return q.DeveloperTerms[0]
	}
	return ""
}

func (q Query) MatchDeveloper(pkg aur.Package) bool {
	if len(q.DeveloperTerms) == 0 {
		return true
	}
	developer := packageDeveloper(pkg)
	if developer == "" {
		return false
	}
	developer = strings.ToLower(developer)
	for _, term := range q.DeveloperTerms {
		term = strings.ToLower(strings.TrimSpace(term))
		if term == "" {
			continue
		}
		if !strings.Contains(developer, term) {
			return false
		}
	}
	return true
}

func isDeveloperKey(key string) bool {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "dev", "developer", "maint", "maintainer":
		return true
	default:
		return false
	}
}

func packageDeveloper(pkg aur.Package) string {
	if pkg.Maintainer == nil {
		return ""
	}
	return strings.TrimSpace(*pkg.Maintainer)
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
	parsed := ParseQuery(query)
	if !parsed.MatchDeveloper(pkg) {
		return false
	}
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
	return matchName(parsed.Text, pkg, s.Match)
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

func SourceBadgeLabel(source string) string {
	label := SourceLabel(source)
	switch strings.ToLower(strings.TrimSpace(label)) {
	case "aur":
		return "AUR"
	case "core":
		return "CORE"
	case "extra":
		return "EXT"
	case "multilib":
		return "MULTI"
	case "chaotic-aur":
		return "CHAOTIC"
	default:
		return compactSourceBadgeLabel(label)
	}
}

func compactSourceBadgeLabel(source string) string {
	source = strings.TrimSpace(source)
	if source == "" {
		return "AUR"
	}
	label := strings.ToUpper(source)
	const maxRunes = 8
	if runeLen(label) <= maxRunes {
		return label
	}
	parts := strings.FieldsFunc(label, func(r rune) bool {
		return r == '-' || r == '_' || r == '/' || r == '.'
	})
	if len(parts) > 1 {
		var b strings.Builder
		for _, part := range parts {
			if part == "" {
				continue
			}
			b.WriteRune([]rune(part)[0])
		}
		if out := b.String(); out != "" && runeLen(out) <= maxRunes {
			return out
		}
	}
	return truncateRunes(label, maxRunes)
}

func runeLen(value string) int {
	return len([]rune(value))
}

func truncateRunes(value string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
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
