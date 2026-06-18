package theme

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const MatugenName = "matugen"

type Theme struct {
	Name               string
	Color              bool
	AccentCode         string
	GoodCode           string
	WarnCode           string
	DangerCode         string
	MutedCode          string
	DimCode            string
	FocusCode          string
	SelectedCode       string
	BadgeCode          string
	SourceAURCode      string
	SourceCoreCode     string
	SourceExtraCode    string
	SourceMultilibCode string
	SourceChaoticCode  string
	SourceUnknownCode  string
	HeaderCode         string
	TableCode          string
	FilterCode         string
	FilterOnCode       string
	FilterHotCode      string
	Separator          string
	PanelDivider       string
	StatusDivider      string
}

type ColorConfig struct {
	Accent      string
	Good        string
	Warn        string
	Danger      string
	Muted       string
	Dim         string
	Focus       string
	SelectedFG  string
	SelectedBG  string
	BadgeFG     string
	BadgeBG     string
	HeaderFG    string
	HeaderBG    string
	FilterFG    string
	FilterBG    string
	FilterOnFG  string
	FilterOnBG  string
	FilterHotFG string
	FilterHotBG string
}

func Detect(name string) (Theme, error) {
	return DetectWithColors(name, ColorConfig{})
}

func DetectWithColors(name string, colors ColorConfig) (Theme, error) {
	t, ok := namedWithColors(name, colors)
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q; available themes: %s", name, strings.Join(Names(), ", "))
	}
	return detectTerminalColor(t), nil
}

func detectTerminalColor(t Theme) Theme {
	term := os.Getenv("TERM")
	t.Color = t.Color && os.Getenv("NO_COLOR") == "" && term != "" && term != "dumb"
	return t
}

func Named(name string) (Theme, bool) {
	return namedWithColors(name, ColorConfig{})
}

func namedWithColors(name string, colors ColorConfig) (Theme, bool) {
	if strings.TrimSpace(name) == "" {
		name = "arch"
	}
	if strings.EqualFold(strings.TrimSpace(name), MatugenName) {
		return Matugen(colors), true
	}
	t, ok := themes[strings.ToLower(name)]
	return t, ok
}

func Names() []string {
	names := make([]string, 0, len(themes)+1)
	for name := range themes {
		names = append(names, name)
	}
	names = append(names, MatugenName)
	sort.Strings(names)
	return names
}

func Matugen(colors ColorConfig) Theme {
	base, _ := themes["arch"]
	base.Name = MatugenName
	base.AccentCode = colorCode(colors.Accent, base.AccentCode)
	base.GoodCode = colorCode(colors.Good, base.GoodCode)
	base.WarnCode = colorCode(colors.Warn, base.WarnCode)
	base.DangerCode = colorCode(colors.Danger, base.DangerCode)
	base.MutedCode = colorCode(colors.Muted, base.MutedCode)
	base.DimCode = dimColorCode(colors.Dim, base.DimCode)
	base.FocusCode = boldColorCode(colors.Focus, base.FocusCode)
	base.SelectedCode = colorPairCode(colors.SelectedFG, colors.SelectedBG, base.SelectedCode)
	base.BadgeCode = colorPairCode(colors.BadgeFG, colors.BadgeBG, base.BadgeCode)
	base.SourceAURCode = base.BadgeCode
	base.SourceCoreCode = base.GoodCode
	base.SourceExtraCode = base.AccentCode
	base.SourceMultilibCode = base.WarnCode
	base.SourceChaoticCode = base.FocusCode
	base.SourceUnknownCode = base.MutedCode
	base.HeaderCode = colorPairCode(colors.HeaderFG, colors.HeaderBG, base.HeaderCode)
	base.TableCode = colorPairCode(colors.FilterFG, colors.FilterBG, base.TableCode)
	base.FilterCode = colorPairCode(colors.FilterFG, colors.FilterBG, base.FilterCode)
	base.FilterOnCode = colorPairCode(colors.FilterOnFG, colors.FilterOnBG, base.FilterOnCode)
	base.FilterHotCode = colorPairCode(colors.FilterHotFG, colors.FilterHotBG, base.FilterHotCode)
	return base
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
func (t Theme) SourceBadgeFor(source, s string) string {
	return t.wrap(t.sourceBadgeCode(source), s)
}
func (t Theme) Header(s string) string      { return t.wrap(t.HeaderCode, s) }
func (t Theme) TableHeader(s string) string { return t.wrap(t.TableCode, s) }
func (t Theme) FilterChip(s string) string  { return t.wrap(t.FilterCode, s) }
func (t Theme) FilterActive(s string) string {
	return t.wrap(t.FilterOnCode, s)
}
func (t Theme) FilterFocused(s string) string {
	return t.wrap(t.FilterHotCode, s)
}

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

func (t Theme) sourceBadgeCode(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "", "aur":
		if t.SourceAURCode != "" {
			return t.SourceAURCode
		}
	case "core":
		if t.SourceCoreCode != "" {
			return t.SourceCoreCode
		}
	case "extra":
		if t.SourceExtraCode != "" {
			return t.SourceExtraCode
		}
	case "multilib":
		if t.SourceMultilibCode != "" {
			return t.SourceMultilibCode
		}
	case "chaotic-aur":
		if t.SourceChaoticCode != "" {
			return t.SourceChaoticCode
		}
	default:
		if t.SourceUnknownCode != "" {
			return t.SourceUnknownCode
		}
	}
	return t.BadgeCode
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

type rgb struct {
	r uint8
	g uint8
	b uint8
}

func colorCode(value, fallback string) string {
	color, ok := parseHexColor(value)
	if !ok {
		return fallback
	}
	return fgCode(color)
}

func boldColorCode(value, fallback string) string {
	color, ok := parseHexColor(value)
	if !ok {
		return fallback
	}
	return "1;" + fgCode(color)
}

func dimColorCode(value, fallback string) string {
	color, ok := parseHexColor(value)
	if !ok {
		return fallback
	}
	return "2;" + fgCode(color)
}

func colorPairCode(fg, bg, fallback string) string {
	foreground, fgOK := parseHexColor(fg)
	background, bgOK := parseHexColor(bg)
	if !fgOK || !bgOK {
		return fallback
	}
	return fgCode(foreground) + ";" + bgCode(background)
}

func fgCode(color rgb) string {
	return fmt.Sprintf("38;2;%d;%d;%d", color.r, color.g, color.b)
}

func bgCode(color rgb) string {
	return fmt.Sprintf("48;2;%d;%d;%d", color.r, color.g, color.b)
}

func parseHexColor(value string) (rgb, bool) {
	value = strings.TrimSpace(value)
	if len(value) != 7 || value[0] != '#' {
		return rgb{}, false
	}
	r, ok := parseHexByte(value[1:3])
	if !ok {
		return rgb{}, false
	}
	g, ok := parseHexByte(value[3:5])
	if !ok {
		return rgb{}, false
	}
	b, ok := parseHexByte(value[5:7])
	if !ok {
		return rgb{}, false
	}
	return rgb{r: r, g: g, b: b}, true
}

func parseHexByte(value string) (uint8, bool) {
	parsed, err := strconv.ParseUint(value, 16, 8)
	if err != nil {
		return 0, false
	}
	return uint8(parsed), true
}

var themes = map[string]Theme{
	"arch": {
		Name:               "arch",
		Color:              true,
		AccentCode:         "38;5;45",
		GoodCode:           "38;5;42",
		WarnCode:           "38;5;214",
		DangerCode:         "38;5;203",
		MutedCode:          "38;5;244",
		DimCode:            "2;38;5;245",
		FocusCode:          "1;38;5;51",
		SelectedCode:       "1;38;5;16;48;5;117",
		BadgeCode:          "1;38;5;16;48;5;171",
		SourceAURCode:      "1;38;5;16;48;5;171",
		SourceCoreCode:     "1;38;5;16;48;5;42",
		SourceExtraCode:    "1;38;5;16;48;5;45",
		SourceMultilibCode: "1;38;5;16;48;5;214",
		SourceChaoticCode:  "1;38;5;231;48;5;90",
		SourceUnknownCode:  "1;38;5;250;48;5;238",
		HeaderCode:         "1;38;5;16;48;5;45",
		TableCode:          "1;38;5;231;48;5;24",
		FilterCode:         "38;5;250;48;5;238",
		FilterOnCode:       "1;38;5;231;48;5;24",
		FilterHotCode:      "1;38;5;16;48;5;51",
		Separator:          "─",
		PanelDivider:       "│",
		StatusDivider:      "╾",
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
		Name:               "dark",
		Color:              true,
		AccentCode:         "38;5;117",
		GoodCode:           "38;5;78",
		WarnCode:           "38;5;221",
		DangerCode:         "38;5;203",
		MutedCode:          "38;5;245",
		DimCode:            "2;38;5;242",
		FocusCode:          "1;38;5;159",
		SelectedCode:       "1;38;5;231;48;5;60",
		BadgeCode:          "1;38;5;16;48;5;176",
		SourceAURCode:      "1;38;5;16;48;5;176",
		SourceCoreCode:     "1;38;5;16;48;5;78",
		SourceExtraCode:    "1;38;5;16;48;5;111",
		SourceMultilibCode: "1;38;5;16;48;5;221",
		SourceChaoticCode:  "1;38;5;231;48;5;90",
		SourceUnknownCode:  "1;38;5;250;48;5;238",
		HeaderCode:         "1;38;5;16;48;5;111",
		TableCode:          "1;38;5;231;48;5;24",
		FilterCode:         "38;5;250;48;5;238",
		FilterOnCode:       "1;38;5;231;48;5;67",
		FilterHotCode:      "1;38;5;16;48;5;159",
		Separator:          "─",
		PanelDivider:       "│",
		StatusDivider:      "╾",
	},
	"light": {
		Name:               "light",
		Color:              true,
		AccentCode:         "38;5;25",
		GoodCode:           "38;5;28",
		WarnCode:           "38;5;130",
		DangerCode:         "38;5;160",
		MutedCode:          "38;5;240",
		DimCode:            "2;38;5;244",
		FocusCode:          "1;38;5;24",
		SelectedCode:       "1;38;5;231;48;5;67",
		BadgeCode:          "1;38;5;231;48;5;125",
		SourceAURCode:      "1;38;5;231;48;5;125",
		SourceCoreCode:     "1;38;5;231;48;5;28",
		SourceExtraCode:    "1;38;5;231;48;5;25",
		SourceMultilibCode: "1;38;5;231;48;5;130",
		SourceChaoticCode:  "1;38;5;231;48;5;90",
		SourceUnknownCode:  "1;38;5;236;48;5;254",
		HeaderCode:         "1;38;5;231;48;5;25",
		TableCode:          "1;38;5;231;48;5;24",
		FilterCode:         "38;5;236;48;5;254",
		FilterOnCode:       "1;38;5;231;48;5;31",
		FilterHotCode:      "1;38;5;231;48;5;24",
		Separator:          "─",
		PanelDivider:       "│",
		StatusDivider:      "╾",
	},
	"high-contrast": {
		Name:               "high-contrast",
		Color:              true,
		AccentCode:         "1;37",
		GoodCode:           "1;32",
		WarnCode:           "1;33",
		DangerCode:         "1;31",
		MutedCode:          "37",
		DimCode:            "2;37",
		FocusCode:          "1;36",
		SelectedCode:       "1;37;45",
		BadgeCode:          "1;35",
		SourceAURCode:      "1;35",
		SourceCoreCode:     "1;32",
		SourceExtraCode:    "1;36",
		SourceMultilibCode: "1;33",
		SourceChaoticCode:  "1;37;45",
		SourceUnknownCode:  "37",
		HeaderCode:         "1;30;47",
		TableCode:          "1;37;44",
		FilterCode:         "37;40",
		FilterOnCode:       "1;30;46",
		FilterHotCode:      "1;30;43",
		Separator:          "=",
		PanelDivider:       "|",
		StatusDivider:      "=",
	},
}
