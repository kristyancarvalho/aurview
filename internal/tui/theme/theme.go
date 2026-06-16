package theme

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Theme struct {
	Name          string
	Color         bool
	AccentCode    string
	GoodCode      string
	WarnCode      string
	DangerCode    string
	MutedCode     string
	DimCode       string
	FocusCode     string
	SelectedCode  string
	BadgeCode     string
	Separator     string
	PanelDivider  string
	StatusDivider string
}

func Detect(name string) (Theme, error) {
	t, ok := Named(name)
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q; available themes: %s", name, strings.Join(Names(), ", "))
	}
	term := os.Getenv("TERM")
	t.Color = t.Color && os.Getenv("NO_COLOR") == "" && term != "" && term != "dumb"
	return t, nil
}

func Named(name string) (Theme, bool) {
	if strings.TrimSpace(name) == "" {
		name = "arch"
	}
	t, ok := themes[strings.ToLower(name)]
	return t, ok
}

func Names() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t Theme) Accent(s string) string      { return t.wrap(t.AccentCode, s) }
func (t Theme) Good(s string) string        { return t.wrap(t.GoodCode, s) }
func (t Theme) Warn(s string) string        { return t.wrap(t.WarnCode, s) }
func (t Theme) Danger(s string) string      { return t.wrap(t.DangerCode, s) }
func (t Theme) Muted(s string) string       { return t.wrap(t.MutedCode, s) }
func (t Theme) Dim(s string) string         { return t.wrap(t.DimCode, s) }
func (t Theme) Focus(s string) string       { return t.wrap(t.FocusCode, s) }
func (t Theme) Reverse(s string) string     { return t.Selected(s) }
func (t Theme) Selected(s string) string    { return t.wrap(t.SelectedCode, s) }
func (t Theme) SourceBadge(s string) string { return t.wrap(t.BadgeCode, s) }

func (t Theme) Status(kind, s string) string {
	switch kind {
	case "error":
		return t.Danger(s)
	case "warn":
		return t.Warn(s)
	case "ok":
		return t.Good(s)
	default:
		return t.Muted(s)
	}
}

func (t Theme) wrap(code, s string) string {
	if !t.Color || s == "" || code == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + len(code) + 8)
	b.WriteString("\x1b[")
	b.WriteString(code)
	b.WriteByte('m')
	b.WriteString(s)
	b.WriteString("\x1b[0m")
	return b.String()
}

var themes = map[string]Theme{
	"arch": {
		Name:          "arch",
		Color:         true,
		AccentCode:    "38;5;45",
		GoodCode:      "38;5;42",
		WarnCode:      "38;5;214",
		DangerCode:    "38;5;203",
		MutedCode:     "38;5;244",
		DimCode:       "2;38;5;245",
		FocusCode:     "1;38;5;51",
		SelectedCode:  "7;38;5;51",
		BadgeCode:     "1;38;5;39",
		Separator:     "─",
		PanelDivider:  "│",
		StatusDivider: "╾",
	},
	"mono": {
		Name:          "mono",
		Color:         false,
		Separator:     "-",
		PanelDivider:  "|",
		StatusDivider: "-",
		SelectedCode:  "",
	},
	"dark": {
		Name:          "dark",
		Color:         true,
		AccentCode:    "38;5;117",
		GoodCode:      "38;5;78",
		WarnCode:      "38;5;221",
		DangerCode:    "38;5;203",
		MutedCode:     "38;5;245",
		DimCode:       "2;38;5;242",
		FocusCode:     "1;38;5;159",
		SelectedCode:  "7;38;5;159",
		BadgeCode:     "1;38;5;111",
		Separator:     "─",
		PanelDivider:  "│",
		StatusDivider: "╾",
	},
	"light": {
		Name:          "light",
		Color:         true,
		AccentCode:    "38;5;25",
		GoodCode:      "38;5;28",
		WarnCode:      "38;5;130",
		DangerCode:    "38;5;160",
		MutedCode:     "38;5;240",
		DimCode:       "2;38;5;244",
		FocusCode:     "1;38;5;24",
		SelectedCode:  "7;38;5;24",
		BadgeCode:     "1;38;5;25",
		Separator:     "─",
		PanelDivider:  "│",
		StatusDivider: "╾",
	},
	"high-contrast": {
		Name:          "high-contrast",
		Color:         true,
		AccentCode:    "1;37",
		GoodCode:      "1;32",
		WarnCode:      "1;33",
		DangerCode:    "1;31",
		MutedCode:     "37",
		DimCode:       "2;37",
		FocusCode:     "1;36",
		SelectedCode:  "7;1;37",
		BadgeCode:     "1;35",
		Separator:     "=",
		PanelDivider:  "|",
		StatusDivider: "=",
	},
}
