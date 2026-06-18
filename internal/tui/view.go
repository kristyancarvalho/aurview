package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/filter"
	"github.com/kristyancarvalho/aurview/internal/platform"
	"github.com/kristyancarvalho/aurview/internal/ranking"
	"github.com/kristyancarvalho/aurview/internal/tui/components"
)

func (m Model) View() string {
	if m.width > 0 && (m.width < 64 || m.height < 12) {
		return m.theme.Danger("aurview needs at least 64x12") + "\n" +
			m.theme.Muted(fmt.Sprintf("current terminal: %dx%d", m.width, m.height))
	}
	if m.help {
		return m.renderHelp()
	}

	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteByte('\n')
	b.WriteString(m.renderSearch())
	b.WriteByte('\n')
	b.WriteString(m.renderFilterBar())
	b.WriteByte('\n')

	if m.width >= 110 {
		leftWidth := max(62, m.width*58/100)
		rightWidth := m.width - leftWidth - 1
		left := strings.Split(m.renderList(leftWidth, m.height-5), "\n")
		right := strings.Split(m.renderDetail(rightWidth, m.height-5), "\n")
		rows := max(len(left), len(right))
		for i := 0; i < rows; i++ {
			l, r := "", ""
			if i < len(left) {
				l = left[i]
			}
			if i < len(right) {
				r = right[i]
			}
			b.WriteString(components.PadRight(l, leftWidth))
			b.WriteString(m.theme.Muted(m.theme.PanelDivider))
			b.WriteString(r)
			if i < rows-1 {
				b.WriteByte('\n')
			}
		}
	} else {
		listHeight := max(3, (m.height-6)*2/3)
		b.WriteString(m.renderList(m.width, listHeight))
		b.WriteByte('\n')
		b.WriteString(m.theme.Muted(components.Repeat(m.theme.Separator, m.width)))
		b.WriteByte('\n')
		b.WriteString(m.renderDetail(m.width, m.height-listHeight-7))
	}

	b.WriteByte('\n')
	b.WriteString(m.renderStatus())
	return b.String()
}

func (m Model) renderHeader() string {
	width := max(m.width, 80)
	count := fmt.Sprintf("%d/%d pkgs", len(m.results), len(m.allResults))
	focus := "search"
	if m.focus == focusFilters {
		focus = "filters"
	} else if m.focus == focusList {
		focus = "list"
	} else if m.focus == focusDetail {
		focus = "detail"
	}
	active := ""
	if n := m.filterState.ActiveCount(); n > 0 {
		active = fmt.Sprintf(" // %d filters", n)
	}
	left := " AURVIEW // read-only rpc"
	right := count + active + " // " + focus + " // ? help "
	spacer := components.Repeat(" ", max(1, width-components.RuneLen("AURVIEW // read-only rpc")-components.RuneLen(count+" // "+focus+" // ? help")))
	line := left + spacer + right
	return m.theme.Header(components.PadRight(components.Truncate(line, width), width))
}

func (m Model) renderSearch() string {
	width := max(m.width, 80)
	marker := " "
	if m.focus == focusSearch {
		marker = m.theme.Focus(">")
	}
	input := m.input
	if m.focus == focusSearch {
		input += m.theme.Muted("_")
	}
	state := ""
	if m.loading {
		state = m.theme.Muted(" searching")
	} else if m.searchError != "" {
		state = m.theme.Danger(" error")
	} else if strings.TrimSpace(m.input) == "" {
		state = m.theme.Muted(" type to search")
	}
	line := fmt.Sprintf("%s / %s%s", marker, input, state)
	return components.Truncate(line, width)
}

func (m Model) renderFilterBar() string {
	width := max(m.width, 80)
	prefix := " filters "
	if m.focus == focusFilters {
		prefix = ">filters "
	}
	parts := []string{m.theme.TableHeader(prefix)}
	for _, chip := range m.filterChips() {
		text := " " + chip.Label + " "
		switch {
		case chip.Focused:
			parts = append(parts, m.theme.FilterFocused(text))
		case chip.Active:
			parts = append(parts, m.theme.FilterActive(text))
		default:
			parts = append(parts, m.theme.FilterChip(text))
		}
	}
	if parsed := filter.ParseQuery(m.input); parsed.HasDeveloper() {
		parts = append(parts, m.theme.FilterActive(" dev:"+parsed.DeveloperLabel()+" "))
	}
	line := strings.Join(parts, " ")
	return components.PadRight(components.Truncate(line, width), width)
}

func (m Model) renderList(width, height int) string {
	var b strings.Builder
	layout := newListLayout(width)
	b.WriteString(m.theme.TableHeader(layout.header()))

	if m.loading && len(m.results) == 0 {
		b.WriteByte('\n')
		b.WriteString(m.theme.Muted("  searching AUR RPC..."))
		return b.String()
	}
	if m.searchError != "" {
		b.WriteByte('\n')
		b.WriteString(m.theme.Danger("  " + components.Truncate(m.searchError, width-2)))
		return b.String()
	}
	if strings.TrimSpace(m.input) == "" {
		b.WriteByte('\n')
		b.WriteString(m.theme.Muted("  enter a package name or keyword"))
		return b.String()
	}
	if len(m.results) == 0 {
		b.WriteByte('\n')
		if len(m.allResults) > 0 && m.filterState.Active() {
			b.WriteString(m.theme.Muted("  no results after filters"))
		} else {
			b.WriteString(m.theme.Muted("  no results"))
		}
		return b.String()
	}

	visible := max(1, height-1)
	end := components.Clamp(m.scroll+visible, 0, len(m.results))
	for i := m.scroll; i < end; i++ {
		b.WriteByte('\n')
		b.WriteString(m.renderRow(i, layout))
	}
	return b.String()
}

func (m Model) renderRow(index int, layout listLayout) string {
	ranked := m.results[index]
	pkg := ranked.Package
	marker := " "
	if index == m.selected {
		marker = ">"
	}
	flag := "ok"
	if pkg.IsOutOfDate() {
		flag = "ood"
	}
	maint := pkg.MaintainerName()
	if pkg.IsOrphan() {
		maint = "orphan"
	}
	values := map[string]string{
		"marker":  marker,
		"source":  filter.SourceBadgeLabel(pkg.Source),
		"package": pkg.Name,
		"version": pkg.Version,
		"score":   fmt.Sprintf("%.1f", ranked.Score),
		"votes":   strconv.Itoa(pkg.NumVotes),
		"pop":     components.FormatPopularity(pkg.Popularity),
		"maint":   maint,
		"updated": platform.UnixDate(pkg.LastModified),
		"flag":    flag,
	}
	var line string
	if index == m.selected {
		line = layout.row(values, nil)
	} else {
		line = layout.row(values, func(s string) string {
			return m.theme.SourceBadgeFor(pkg.Source, s)
		})
	}
	line = components.PadRight(line, layout.width)
	if index == m.selected {
		if m.focus == focusList {
			return m.theme.Selected(line)
		}
		return m.theme.Focus(line)
	}
	if pkg.IsOutOfDate() || pkg.IsOrphan() {
		return m.theme.Warn(line)
	}
	return line
}

func (m Model) renderDetail(width, height int) string {
	var b strings.Builder
	pkg, ok := m.selectedPackage()
	if !ok {
		return m.theme.TableHeader(components.PadRight("detail // no package selected", width))
	}
	title := "detail // " + pkg.Name
	if m.detailLoading {
		title += " // loading"
	}
	title = components.PadRight(components.Truncate(title, width), width)
	if m.focus == focusDetail {
		title = m.theme.Header(title)
	} else {
		title = m.theme.TableHeader(title)
	}
	b.WriteString(title)

	lines := m.detailLines(pkg, width)
	if m.detailError != "" {
		lines = append([]string{m.theme.Danger("detail error: " + m.detailError)}, lines...)
	}
	limit := max(1, height-1)
	start := components.Clamp(m.detailScroll, 0, max(0, len(lines)-limit))
	for i := start; i < len(lines) && i < start+limit; i++ {
		b.WriteByte('\n')
		b.WriteString(components.PadRight(components.Truncate(lines[i], width), width))
	}
	return b.String()
}

func (m Model) detailLines(pkg aur.Package, width int) []string {
	labelWidth := 11
	out := []string{}
	addSection := func(title string) {
		if len(out) > 0 {
			out = append(out, "")
		}
		out = append(out, m.theme.Muted(title))
	}
	addField := func(label, value string) {
		out = append(out, wrapDetailField(label, value, width, labelWidth)...)
	}

	addSection("Identity")
	addField("source", pkg.DisplaySource())
	addField("base", pkg.PackageBase)
	addField("version", pkg.Version)
	addField("maintainer", pkg.MaintainerName())
	addField("license", join(pkg.License))

	addSection("Health")
	addField("score", selectedScore(m.results, m.selected))
	addField("votes", fmt.Sprintf("%d", pkg.NumVotes))
	addField("popularity", components.FormatPopularity(pkg.Popularity))
	addField("first", platform.UnixDate(pkg.FirstSubmitted))
	addField("updated", platform.UnixDate(pkg.LastModified))
	addField("out-of-date", platform.OptionalUnixDate(pkg.OutOfDate))

	addSection("Links")
	addField("upstream", pkg.URL)
	addField("aur", pkg.AURURL())

	addSection("Description")
	out = append(out, wrapDetailText(valueOrDash(pkg.Description), width)...)

	addSection("Relations")
	addField("depends", join(pkg.Depends))
	addField("make", join(pkg.MakeDepends))
	addField("check", join(pkg.CheckDepends))
	addField("optional", join(pkg.OptDepends))
	addField("conflicts", join(pkg.Conflicts))
	addField("provides", join(pkg.Provides))

	addSection("Keywords")
	out = append(out, wrapDetailText(join(pkg.Keywords), width)...)
	return out
}

func (m Model) renderStatus() string {
	width := max(m.width, 80)
	left := components.Truncate(m.status, max(10, width-18))
	right := "q quit  / search"
	gap := components.Repeat(" ", max(1, width-components.RuneLen(left)-components.RuneLen(right)))
	return m.theme.Status(m.statusKind, left) + m.theme.Muted(gap+right)
}

func (m Model) renderHelp() string {
	width := max(m.width, 80)
	lines := []string{
		m.theme.Focus("AURVIEW // keys"),
		"",
		"/          focus search",
		"j/k        move selection",
		"h/l        move focus",
		"gg / G     top / bottom",
		"ctrl+d/u   half page",
		"ctrl+f/b   page",
		"n / N      history newer / older",
		"f / tab    focus filters / next filter",
		"space      cycle focused filter",
		"r          reset filters when filters are focused",
		"enter      copy selected package name",
		"mouse      click rows/search, wheel list/detail",
		"esc        blur or close",
		"?          close help",
		"q          quit",
		"",
		"read-only: no install, clone, build, update or remove actions exist",
	}
	for i, line := range lines {
		lines[i] = components.Truncate(line, width)
	}
	return strings.Join(lines, "\n")
}

func kv(label, value string) string {
	return label + ": " + valueOrDash(value)
}

func join(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, " ")
}

func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

type columnAlign int

const (
	alignLeft columnAlign = iota
	alignRight
)

type listColumn struct {
	key       string
	header    string
	width     int
	minWidth  int
	grow      bool
	align     columnAlign
	hideOrder int
}

type listLayout struct {
	width   int
	columns []listColumn
}

func newListLayout(width int) listLayout {
	width = max(1, width)
	columns := []listColumn{
		{key: "marker", header: "", width: 1, minWidth: 1, align: alignLeft},
		{key: "source", header: "src", width: 7, minWidth: 3, align: alignLeft},
		{key: "package", header: "package", width: 18, minWidth: 8, grow: true, align: alignLeft},
		{key: "version", header: "version", width: 11, minWidth: 7, grow: true, align: alignLeft, hideOrder: 5},
		{key: "score", header: "score", width: 5, minWidth: 5, align: alignRight},
		{key: "votes", header: "votes", width: 5, minWidth: 5, align: alignRight, hideOrder: 2},
		{key: "pop", header: "pop", width: 7, minWidth: 5, align: alignRight, hideOrder: 3},
		{key: "maint", header: "maint", width: 10, minWidth: 7, grow: true, align: alignLeft, hideOrder: 1},
		{key: "updated", header: "updated", width: 10, minWidth: 10, align: alignLeft},
		{key: "flag", header: "flag", width: 4, minWidth: 3, align: alignLeft, hideOrder: 4},
	}
	for totalColumnWidth(columns) > width {
		index := nextHiddenColumn(columns)
		if index < 0 {
			break
		}
		columns = append(columns[:index], columns[index+1:]...)
	}
	for totalColumnWidth(columns) > width {
		index := widestGrowColumn(columns)
		if index < 0 || columns[index].width <= columns[index].minWidth {
			break
		}
		columns[index].width--
	}
	extra := width - totalColumnWidth(columns)
	for extra > 0 {
		index := growColumn(columns)
		if index < 0 {
			break
		}
		columns[index].width++
		extra--
	}
	return listLayout{width: width, columns: columns}
}

func (l listLayout) header() string {
	values := make(map[string]string, len(l.columns))
	for _, col := range l.columns {
		values[col.key] = col.header
	}
	return components.PadRight(l.format(values, nil), l.width)
}

func (l listLayout) row(values map[string]string, styleSource func(string) string) string {
	return components.PadRight(l.format(values, styleSource), l.width)
}

func (l listLayout) format(values map[string]string, styleSource func(string) string) string {
	parts := make([]string, 0, len(l.columns))
	for _, col := range l.columns {
		if col.key == "source" && styleSource != nil {
			value := components.Truncate(values[col.key], col.width)
			parts = append(parts, components.PadRight(styleSource(value), col.width))
			continue
		}
		value := formatCell(values[col.key], col.width, col.align)
		parts = append(parts, value)
	}
	return strings.Join(parts, " ")
}

func formatCell(value string, width int, align columnAlign) string {
	value = components.Truncate(value, width)
	if align == alignRight {
		return components.PadLeft(value, width)
	}
	return components.PadRight(value, width)
}

func totalColumnWidth(columns []listColumn) int {
	if len(columns) == 0 {
		return 0
	}
	total := len(columns) - 1
	for _, col := range columns {
		total += col.width
	}
	return total
}

func nextHiddenColumn(columns []listColumn) int {
	bestIndex := -1
	bestOrder := 999
	for i, col := range columns {
		if col.hideOrder <= 0 {
			continue
		}
		if col.hideOrder < bestOrder {
			bestOrder = col.hideOrder
			bestIndex = i
		}
	}
	return bestIndex
}

func widestGrowColumn(columns []listColumn) int {
	bestIndex := -1
	bestWidth := 0
	for i, col := range columns {
		if !col.grow || col.width <= col.minWidth {
			continue
		}
		if col.width > bestWidth {
			bestWidth = col.width
			bestIndex = i
		}
	}
	return bestIndex
}

func growColumn(columns []listColumn) int {
	for i, col := range columns {
		if col.key == "package" && col.grow {
			return i
		}
	}
	for i, col := range columns {
		if col.grow {
			return i
		}
	}
	return -1
}

func selectedScore(results []ranking.RankedPackage, selected int) string {
	if selected < 0 || selected >= len(results) {
		return "-"
	}
	return fmt.Sprintf("%.1f", results[selected].Score)
}

func wrapDetailField(label, value string, width, labelWidth int) []string {
	label = components.Truncate(label, labelWidth)
	prefix := "  " + components.PadRight(label, labelWidth) + " "
	return wrapWithPrefix(prefix, components.Repeat(" ", components.RuneLen(prefix)), valueOrDash(value), width)
}

func wrapDetailText(value string, width int) []string {
	return wrapWithPrefix("  ", "  ", valueOrDash(value), width)
}

func wrapWithPrefix(firstPrefix, nextPrefix, value string, width int) []string {
	if width <= 0 {
		return nil
	}
	if width <= components.RuneLen(firstPrefix)+3 {
		return []string{components.Truncate(firstPrefix+value, width)}
	}
	words := strings.Fields(valueOrDash(value))
	if len(words) == 0 {
		words = []string{"-"}
	}
	lines := []string{}
	prefix := firstPrefix
	line := prefix
	for _, word := range words {
		available := width - components.RuneLen(prefix)
		if components.RuneLen(word) > available {
			if strings.TrimSpace(line) != "" && line != prefix {
				lines = append(lines, components.PadRight(line, width))
			}
			lines = append(lines, components.PadRight(prefix+components.Truncate(word, available), width))
			prefix = nextPrefix
			line = prefix
			continue
		}
		separator := ""
		if line != prefix {
			separator = " "
		}
		if components.RuneLen(line)+components.RuneLen(separator)+components.RuneLen(word) > width {
			lines = append(lines, components.PadRight(line, width))
			prefix = nextPrefix
			line = prefix + word
			continue
		}
		line += separator + word
	}
	if strings.TrimSpace(line) == "" {
		line = prefix + "-"
	}
	lines = append(lines, components.PadRight(line, width))
	return lines
}
