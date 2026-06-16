package tui

import (
	"fmt"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/platform"
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

	if m.width >= 110 {
		leftWidth := max(62, m.width*58/100)
		rightWidth := m.width - leftWidth - 1
		left := strings.Split(m.renderList(leftWidth, m.height-4), "\n")
		right := strings.Split(m.renderDetail(rightWidth, m.height-4), "\n")
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
			b.WriteString(m.theme.Muted("|"))
			b.WriteString(r)
			if i < rows-1 {
				b.WriteByte('\n')
			}
		}
	} else {
		listHeight := max(3, (m.height-5)*2/3)
		b.WriteString(m.renderList(m.width, listHeight))
		b.WriteByte('\n')
		b.WriteString(m.theme.Muted(components.Repeat("-", m.width)))
		b.WriteByte('\n')
		b.WriteString(m.renderDetail(m.width, m.height-listHeight-6))
	}

	b.WriteByte('\n')
	b.WriteString(m.renderStatus())
	return b.String()
}

func (m Model) renderHeader() string {
	width := max(m.width, 80)
	count := fmt.Sprintf("%d pkgs", len(m.results))
	focus := "search"
	if m.focus == focusList {
		focus = "list"
	}
	if m.focus == focusDetail {
		focus = "detail"
	}
	left := m.theme.Focus("AURVIEW") + m.theme.Muted(" // read-only rpc")
	right := m.theme.Muted(count + " // " + focus + " // ? help")
	spacer := components.Repeat(" ", max(1, width-components.RuneLen("AURVIEW // read-only rpc")-components.RuneLen(count+" // "+focus+" // ? help")))
	return left + spacer + right
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

func (m Model) renderList(width, height int) string {
	var b strings.Builder
	header := fmt.Sprintf("  %-24s %-11s %5s %5s %7s %-10s %-10s %-4s %s", "package", "version", "score", "votes", "pop", "maint", "updated", "flag", "description")
	b.WriteString(m.theme.Muted(components.Truncate(header, width)))

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
		b.WriteString(m.theme.Muted("  no results"))
		return b.String()
	}

	visible := max(1, height-1)
	end := components.Clamp(m.scroll+visible, 0, len(m.results))
	for i := m.scroll; i < end; i++ {
		b.WriteByte('\n')
		b.WriteString(m.renderRow(i, width))
	}
	return b.String()
}

func (m Model) renderRow(index, width int) string {
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
	line := fmt.Sprintf("%s %-24s %-11s %5.1f %5d %7s %-10s %-10s %-4s %s",
		marker,
		components.Truncate(pkg.Name, 24),
		components.Truncate(pkg.Version, 11),
		ranked.Score,
		pkg.NumVotes,
		components.FormatPopularity(pkg.Popularity),
		components.Truncate(maint, 10),
		platform.UnixDate(pkg.LastModified),
		flag,
		pkg.Description,
	)
	line = components.Truncate(line, width)
	if index == m.selected {
		if m.focus == focusList {
			return m.theme.Reverse(line)
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
		return m.theme.Muted("detail // no package selected")
	}
	title := "detail // " + pkg.Name
	if m.detailLoading {
		title += " // loading"
	}
	if m.focus == focusDetail {
		title = m.theme.Focus(title)
	} else {
		title = m.theme.Muted(title)
	}
	b.WriteString(components.Truncate(title, width))

	lines := m.detailLines(pkg, width)
	if m.detailError != "" {
		lines = append([]string{m.theme.Danger("detail error: " + m.detailError)}, lines...)
	}
	limit := max(1, height-1)
	start := components.Clamp(m.detailScroll, 0, max(0, len(lines)-limit))
	for i := start; i < len(lines) && i < start+limit; i++ {
		b.WriteByte('\n')
		b.WriteString(components.Truncate(lines[i], width))
	}
	return b.String()
}

func (m Model) detailLines(pkg aur.Package, width int) []string {
	out := []string{
		kv("base", pkg.PackageBase),
		kv("version", pkg.Version),
		kv("maintainer", pkg.MaintainerName()),
		kv("votes", fmt.Sprintf("%d", pkg.NumVotes)),
		kv("popularity", components.FormatPopularity(pkg.Popularity)),
		kv("first", platform.UnixDate(pkg.FirstSubmitted)),
		kv("modified", platform.UnixDate(pkg.LastModified)),
		kv("out-of-date", platform.OptionalUnixDate(pkg.OutOfDate)),
	}
	if pkg.URL != "" {
		out = append(out, kv("upstream", pkg.URL))
	}
	out = append(out, kv("aur", pkg.AURURL()))
	out = append(out, components.WrapLine("desc", valueOrDash(pkg.Description), width)...)
	out = append(out, components.WrapLine("license", join(pkg.License), width)...)
	out = append(out, components.WrapLine("depends", join(pkg.Depends), width)...)
	out = append(out, components.WrapLine("make", join(pkg.MakeDepends), width)...)
	out = append(out, components.WrapLine("check", join(pkg.CheckDepends), width)...)
	out = append(out, components.WrapLine("optional", join(pkg.OptDepends), width)...)
	out = append(out, components.WrapLine("conflicts", join(pkg.Conflicts), width)...)
	out = append(out, components.WrapLine("provides", join(pkg.Provides), width)...)
	out = append(out, components.WrapLine("keywords", join(pkg.Keywords), width)...)
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
