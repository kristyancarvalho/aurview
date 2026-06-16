package theme

import (
	"os"
	"strings"
)

type Theme struct {
	Color bool
}

func Detect() Theme {
	term := os.Getenv("TERM")
	return Theme{Color: os.Getenv("NO_COLOR") == "" && term != "" && term != "dumb"}
}

func (t Theme) Accent(s string) string  { return t.wrap("38;5;45", s) }
func (t Theme) Good(s string) string    { return t.wrap("38;5;42", s) }
func (t Theme) Warn(s string) string    { return t.wrap("38;5;214", s) }
func (t Theme) Danger(s string) string  { return t.wrap("38;5;203", s) }
func (t Theme) Muted(s string) string   { return t.wrap("38;5;244", s) }
func (t Theme) Dim(s string) string     { return t.wrap("2;38;5;245", s) }
func (t Theme) Focus(s string) string   { return t.wrap("1;38;5;51", s) }
func (t Theme) Reverse(s string) string { return t.wrap("7", s) }

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
	if !t.Color || s == "" {
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
