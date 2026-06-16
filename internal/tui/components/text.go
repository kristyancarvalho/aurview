package components

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if RuneLen(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	runes := []rune(s)
	return string(runes[:width-1]) + "…"
}

func PadRight(s string, width int) string {
	l := RuneLen(s)
	if l >= width {
		return Truncate(s, width)
	}
	return s + strings.Repeat(" ", width-l)
}

func PadLeft(s string, width int) string {
	l := RuneLen(s)
	if l >= width {
		return Truncate(s, width)
	}
	return strings.Repeat(" ", width-l) + s
}

func Repeat(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return strings.Repeat(s, width)
}

func RuneLen(s string) int {
	return utf8.RuneCountInString(s)
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func FormatPopularity(value float64) string {
	switch {
	case value >= 100:
		return fmt.Sprintf("%.0f", value)
	case value >= 10:
		return fmt.Sprintf("%.1f", value)
	default:
		return fmt.Sprintf("%.2f", value)
	}
}

func WrapLine(label, value string, width int) []string {
	if width <= 0 {
		return nil
	}
	prefix := label + ": "
	available := width - RuneLen(prefix)
	if available < 4 {
		return []string{Truncate(prefix+value, width)}
	}
	words := strings.Fields(value)
	if len(words) == 0 {
		return []string{Truncate(prefix+"-", width)}
	}
	lines := []string{}
	current := prefix
	for _, word := range words {
		if RuneLen(current)+RuneLen(word)+1 > width {
			lines = append(lines, Truncate(current, width))
			current = strings.Repeat(" ", RuneLen(prefix)) + word
			continue
		}
		if strings.TrimSpace(current) != strings.TrimSpace(prefix) {
			current += " "
		}
		current += word
	}
	if strings.TrimSpace(current) != "" {
		lines = append(lines, Truncate(current, width))
	}
	return lines
}
