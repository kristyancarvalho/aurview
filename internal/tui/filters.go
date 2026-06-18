package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/aurview/internal/filter"
)

var (
	voteSteps       = []int{0, 10, 50, 100}
	popularitySteps = []float64{0, 1, 5, 10}
	recencySteps    = []int{0, 30, 180, 365}
)

type filterChip struct {
	Label   string
	Active  bool
	Focused bool
}

func (m *Model) applyFilters() {
	m.ensureSourceFilterValid()
	m.results = m.filterState.Apply(m.lastQuery, m.allResults, m.scorer.Now)
	if len(m.results) == 0 {
		m.selected = 0
		m.scroll = 0
		m.detailScroll = 0
		m.detailLoading = false
		return
	}
	m.selected = componentsClamp(m.selected, 0, len(m.results)-1)
	m.ensureSelectionVisible()
}

func (m *Model) setFilteredStatus(query string) {
	active := m.filterState.ActiveCount()
	parsed := filter.ParseQuery(query)
	developer := parsed.DeveloperLabel()
	if active == 0 {
		if developer != "" {
			if parsed.Text != "" {
				m.status = fmt.Sprintf("%d packages ranked for %q with developer %q", len(m.results), parsed.Text, developer)
				return
			}
			m.status = fmt.Sprintf("%d packages matched developer %q", len(m.results), developer)
			return
		}
		m.status = fmt.Sprintf("%d packages ranked for %q", len(m.results), query)
		return
	}
	if developer != "" {
		m.status = fmt.Sprintf("%d of %d packages after %d filters with developer %q", len(m.results), len(m.allResults), active, developer)
		return
	}
	m.status = fmt.Sprintf("%d of %d packages after %d filters", len(m.results), len(m.allResults), active)
}

func (m *Model) moveFilter(delta int) {
	m.focus = focusFilters
	chips := m.filterChips()
	if len(chips) == 0 {
		m.filterIndex = 0
		return
	}
	m.filterIndex = (m.filterIndex + delta + len(chips)) % len(chips)
	m.status = "filter " + chips[m.filterIndex].Label
	m.statusKind = "info"
}

func (m *Model) cycleSelectedFilter() tea.Cmd {
	m.focus = focusFilters
	m.ensureFilterIndex()
	switch m.filterIndex {
	case 0:
		m.cycleSourceFilter()
	case 1:
		m.filterState.Maintainer = (m.filterState.Maintainer + 1) % 3
	case 2:
		m.filterState.Flag = (m.filterState.Flag + 1) % 3
	case 3:
		m.filterState.MinVotes = nextIntStep(m.filterState.MinVotes, voteSteps)
	case 4:
		m.filterState.MinPopularity = nextFloatStep(m.filterState.MinPopularity, popularitySteps)
	case 5:
		m.filterState.RecentDays = nextIntStep(m.filterState.RecentDays, recencySteps)
	case 6:
		m.filterState.Match = (m.filterState.Match + 1) % 4
	}
	m.selected = 0
	m.scroll = 0
	m.detailScroll = 0
	m.applyFilters()
	m.setFilteredStatus(m.lastQuery)
	if len(m.results) == 0 {
		m.statusKind = "warn"
		return nil
	}
	m.statusKind = "ok"
	return m.fetchSelectedDetail()
}

func (m *Model) resetFilters() tea.Cmd {
	m.filterState = filter.State{}
	m.filterIndex = 0
	m.selected = 0
	m.scroll = 0
	m.detailScroll = 0
	m.applyFilters()
	if len(m.results) == 0 {
		m.status = "filters reset"
		m.statusKind = "info"
		return nil
	}
	m.setFilteredStatus(m.lastQuery)
	m.statusKind = "ok"
	return m.fetchSelectedDetail()
}

func (m *Model) ensureFilterIndex() {
	chips := m.filterChips()
	if len(chips) == 0 {
		m.filterIndex = 0
		return
	}
	if m.filterIndex < 0 || m.filterIndex >= len(chips) {
		m.filterIndex = 0
	}
}

func (m *Model) ensureSourceFilterValid() {
	if m.filterState.Source == "" {
		return
	}
	for _, source := range m.sourceOptions() {
		if source == m.filterState.Source {
			return
		}
	}
	m.filterState.Source = ""
}

func (m *Model) cycleSourceFilter() {
	sources := m.sourceOptions()
	if len(sources) == 0 {
		m.filterState.Source = ""
		return
	}
	if m.filterState.Source == "" {
		m.filterState.Source = sources[0]
		return
	}
	for i, source := range sources {
		if source == m.filterState.Source {
			if i == len(sources)-1 {
				m.filterState.Source = ""
				return
			}
			m.filterState.Source = sources[i+1]
			return
		}
	}
	m.filterState.Source = ""
}

func (m Model) sourceOptions() []string {
	return filter.Sources(m.allResults)
}

func (m Model) filterChips() []filterChip {
	chips := []filterChip{
		{Label: "src:" + sourceFilterLabel(m.filterState.Source), Active: m.filterState.Source != ""},
		{Label: "maint-state:" + maintainerFilterLabel(m.filterState.Maintainer), Active: m.filterState.Maintainer != filter.MaintainerAny},
		{Label: "flag:" + flagFilterLabel(m.filterState.Flag), Active: m.filterState.Flag != filter.FlagAny},
		{Label: votesFilterLabel(m.filterState.MinVotes), Active: m.filterState.MinVotes > 0},
		{Label: popularityFilterLabel(m.filterState.MinPopularity), Active: m.filterState.MinPopularity > 0},
		{Label: recencyFilterLabel(m.filterState.RecentDays), Active: m.filterState.RecentDays > 0},
		{Label: "match:" + matchFilterLabel(m.filterState.Match), Active: m.filterState.Match != filter.MatchSmart},
	}
	for i := range chips {
		chips[i].Focused = m.focus == focusFilters && i == m.filterIndex
	}
	return chips
}

func sourceFilterLabel(source string) string {
	if source == "" {
		return "all"
	}
	return filter.SourceLabel(source)
}

func maintainerFilterLabel(mode filter.MaintainerMode) string {
	switch mode {
	case filter.MaintainerMaintained:
		return "yes"
	case filter.MaintainerOrphaned:
		return "orphan"
	default:
		return "any"
	}
}

func flagFilterLabel(mode filter.FlagMode) string {
	switch mode {
	case filter.FlagClean:
		return "ok"
	case filter.FlagOutOfDate:
		return "ood"
	default:
		return "any"
	}
}

func votesFilterLabel(votes int) string {
	if votes <= 0 {
		return "votes:any"
	}
	return fmt.Sprintf("votes>=%d", votes)
}

func popularityFilterLabel(popularity float64) string {
	if popularity <= 0 {
		return "pop:any"
	}
	return fmt.Sprintf("pop>=%.0f", popularity)
}

func recencyFilterLabel(days int) string {
	if days <= 0 {
		return "upd:any"
	}
	return fmt.Sprintf("upd<=%dd", days)
}

func matchFilterLabel(mode filter.MatchMode) string {
	switch mode {
	case filter.MatchName:
		return "name"
	case filter.MatchPrefix:
		return "prefix"
	case filter.MatchExact:
		return "exact"
	default:
		return "smart"
	}
}

func nextIntStep(current int, steps []int) int {
	for i, step := range steps {
		if current == step {
			return steps[(i+1)%len(steps)]
		}
	}
	return steps[0]
}

func nextFloatStep(current float64, steps []float64) float64 {
	for i, step := range steps {
		if current == step {
			return steps[(i+1)%len(steps)]
		}
	}
	return steps[0]
}

func componentsClamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
