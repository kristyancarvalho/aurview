package tui

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/ranking"
	"github.com/kristyancarvalho/aurview/internal/tui/components"
	"github.com/kristyancarvalho/aurview/internal/tui/theme"
)

func TestListLayoutAlignsHeaderAndRows(t *testing.T) {
	model := viewTestModel()
	layout := newListLayout(80)

	header := layout.header()
	row := model.renderRow(0, layout)

	if components.RuneLen(header) != 80 {
		t.Fatalf("header width = %d, want 80: %q", components.RuneLen(header), header)
	}
	if components.RuneLen(row) != 80 {
		t.Fatalf("row width = %d, want 80: %q", components.RuneLen(row), row)
	}
	if strings.Index(header, "package") != strings.Index(components.StripANSI(row), "a-very") {
		t.Fatalf("package column is misaligned\nheader: %q\nrow:    %q", header, row)
	}
}

func TestListLayoutHidesColumnsResponsively(t *testing.T) {
	layout := newListLayout(62)
	header := layout.header()

	for _, hidden := range []string{"maint", "votes", "pop"} {
		if strings.Contains(header, hidden) {
			t.Fatalf("header %q contains hidden column %q", header, hidden)
		}
	}
	for _, visible := range []string{"src", "package", "score", "updated"} {
		if !strings.Contains(header, visible) {
			t.Fatalf("header %q missing required column %q", header, visible)
		}
	}
}

func TestSelectedRowIsPaddedBeforeStyling(t *testing.T) {
	model := viewTestModel()
	model.focus = focusList
	model.theme = theme.Theme{Color: true, SelectedCode: "7"}
	layout := newListLayout(72)

	row := model.renderRow(0, layout)

	if !strings.HasPrefix(row, "\x1b[7m") || !strings.HasSuffix(row, "\x1b[0m") {
		t.Fatalf("selected row is not wrapped in selected style: %q", row)
	}
	if strings.Contains(row, "\x1b[0m\x1b[") {
		t.Fatalf("selected row contains nested ANSI reset/style: %q", row)
	}
	if components.RuneLen(row) != 72 {
		t.Fatalf("selected row width = %d, want 72: %q", components.RuneLen(row), row)
	}
}

func TestWideViewKeepsSelectedRowHighlight(t *testing.T) {
	model := viewTestModel()
	model.width = 120
	model.height = 24
	model.focus = focusList
	model.input = "test"
	model.theme = theme.Theme{
		Color:        true,
		SelectedCode: "7",
		PanelDivider: "|",
		Separator:    "-",
	}

	view := model.View()

	if !strings.Contains(view, "\x1b[7m>") {
		t.Fatalf("wide view lost selected-row highlight:\n%q", view)
	}
}

func TestDetailLinesUseSectionsAndDashEmptyValues(t *testing.T) {
	model := viewTestModel()
	pkg := model.results[0].Package
	pkg.URL = ""
	pkg.Depends = nil
	pkg.Keywords = []string{"runner", "tests"}

	lines := strings.Join(stripLines(model.detailLines(pkg, 56)), "\n")

	for _, want := range []string{"Identity", "Health", "Links", "Description", "Relations", "Keywords"} {
		if !strings.Contains(lines, want) {
			t.Fatalf("detail lines missing section %q:\n%s", want, lines)
		}
	}
	for _, want := range []string{"  upstream    -", "  depends     -", "  runner tests"} {
		if !strings.Contains(lines, want) {
			t.Fatalf("detail lines missing formatted value %q:\n%s", want, lines)
		}
	}
}

func TestDetailWrappingTruncatesLongURLs(t *testing.T) {
	lines := wrapDetailField("upstream", "https://example.invalid/"+strings.Repeat("segment", 20), 40, 11)

	if len(lines) == 0 {
		t.Fatal("wrapDetailField returned no lines")
	}
	for _, line := range lines {
		if components.RuneLen(line) != 40 {
			t.Fatalf("line width = %d, want 40: %q", components.RuneLen(line), line)
		}
	}
}

func viewTestModel() Model {
	maint := "long-maintainer-name"
	model := New(Options{})
	model.theme = theme.Theme{}
	model.focus = focusList
	model.selected = 0
	model.results = []ranking.RankedPackage{{
		Package: aur.Package{
			Source:       "aur",
			Name:         "a-very-long-package-name-that-needs-truncation",
			PackageBase:  "test-base",
			Version:      "1.2.3.r456789-1",
			Description:  "Universal test runner with auto-detection for many languages",
			URL:          "https://example.invalid/project",
			NumVotes:     42,
			Popularity:   12.34,
			LastModified: 1_700_000_000,
			Maintainer:   &maint,
			Depends:      []string{"glibc", "libgcc"},
			Conflicts:    []string{"testx"},
			Provides:     []string{"testx"},
			Keywords:     []string{"runner", "tests"},
		},
		Score: 93.8,
	}}
	return model
}

func stripLines(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = components.StripANSI(line)
	}
	return out
}
