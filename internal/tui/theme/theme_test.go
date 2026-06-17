package theme

import "testing"

func TestNamedThemes(t *testing.T) {
	for _, name := range []string{"arch", "mono", "dark", "light", "high-contrast", "matugen"} {
		t.Run(name, func(t *testing.T) {
			got, ok := Named(name)
			if !ok {
				t.Fatalf("Named(%q) missing", name)
			}
			if got.Name != name {
				t.Fatalf("Name = %q, want %q", got.Name, name)
			}
		})
	}
}

func TestMatugenThemeParsesHexColors(t *testing.T) {
	tm := Matugen(ColorConfig{
		Accent:      "#112233",
		Good:        "#445566",
		Warn:        "#778899",
		Danger:      "#aabbcc",
		Muted:       "#010203",
		Dim:         "#040506",
		Focus:       "#070809",
		SelectedFG:  "#ffffff",
		SelectedBG:  "#000000",
		BadgeFG:     "#101112",
		BadgeBG:     "#131415",
		HeaderFG:    "#161718",
		HeaderBG:    "#191a1b",
		FilterFG:    "#202122",
		FilterBG:    "#232425",
		FilterOnFG:  "#262728",
		FilterOnBG:  "#292a2b",
		FilterHotFG: "#303132",
		FilterHotBG: "#333435",
	})

	tests := map[string]string{
		"accent":     tm.Accent("x"),
		"good":       tm.Good("x"),
		"warn":       tm.Warn("x"),
		"danger":     tm.Danger("x"),
		"muted":      tm.Muted("x"),
		"dim":        tm.Dim("x"),
		"focus":      tm.Focus("x"),
		"selected":   tm.Selected("x"),
		"badge":      tm.SourceBadge("x"),
		"header":     tm.Header("x"),
		"filter":     tm.FilterChip("x"),
		"filter-on":  tm.FilterActive("x"),
		"filter-hot": tm.FilterFocused("x"),
	}
	for name, got := range tests {
		if got == "x" {
			t.Fatalf("%s style did not emit ANSI color", name)
		}
	}
	if got := tm.Accent("x"); got != "\x1b[38;2;17;34;51mx\x1b[0m" {
		t.Fatalf("Accent() = %q, want true-color ANSI", got)
	}
	if got := tm.Selected("x"); got != "\x1b[38;2;255;255;255;48;2;0;0;0mx\x1b[0m" {
		t.Fatalf("Selected() = %q, want fg/bg true-color ANSI", got)
	}
}

func TestMatugenInvalidHexFallsBack(t *testing.T) {
	tm := Matugen(ColorConfig{
		Accent:     "not-a-color",
		SelectedFG: "#ffffff",
		SelectedBG: "invalid",
	})
	fallback, _ := Named("arch")

	if got, want := tm.Accent("x"), fallback.Accent("x"); got != want {
		t.Fatalf("invalid accent fallback = %q, want %q", got, want)
	}
	if got, want := tm.Selected("x"), fallback.Selected("x"); got != want {
		t.Fatalf("invalid selected fallback = %q, want %q", got, want)
	}
}

func TestMatugenMissingColorsFallBack(t *testing.T) {
	tm := Matugen(ColorConfig{Accent: "#112233"})
	fallback, _ := Named("arch")

	if got := tm.Accent("x"); got != "\x1b[38;2;17;34;51mx\x1b[0m" {
		t.Fatalf("Accent() = %q, want configured color", got)
	}
	if got, want := tm.Header("x"), fallback.Header("x"); got != want {
		t.Fatalf("missing header fallback = %q, want %q", got, want)
	}
	if got, want := tm.FilterActive("x"), fallback.FilterActive("x"); got != want {
		t.Fatalf("missing active filter fallback = %q, want %q", got, want)
	}
}

func TestDetectWithColorsKeepsBuiltInThemesCompatible(t *testing.T) {
	got, err := DetectWithColors("arch", ColorConfig{Accent: "#000000"})
	if err != nil {
		t.Fatalf("DetectWithColors() error = %v", err)
	}
	want, _ := Detect("arch")
	if got.AccentCode != want.AccentCode || got.Name != want.Name {
		t.Fatalf("DetectWithColors built-in = %#v, want %#v", got, want)
	}
}

func TestDetectRejectsUnknownTheme(t *testing.T) {
	if _, err := Detect("missing"); err == nil {
		t.Fatal("Detect() error = nil, want unknown theme error")
	}
}

func TestMonoThemeDoesNotEmitANSI(t *testing.T) {
	tm, ok := Named("mono")
	if !ok {
		t.Fatal("mono theme missing")
	}
	if got := tm.Accent("AUR"); got != "AUR" {
		t.Fatalf("Accent() = %q, want plain text", got)
	}
}
