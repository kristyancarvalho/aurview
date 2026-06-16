package components

import "testing"

func TestTruncateAndPad(t *testing.T) {
	if got := Truncate("abcdef", 4); got != "abc…" {
		t.Fatalf("Truncate() = %q", got)
	}
	if got := PadRight("ab", 4); got != "ab  " {
		t.Fatalf("PadRight() = %q", got)
	}
	if got := PadLeft("ab", 4); got != "  ab" {
		t.Fatalf("PadLeft() = %q", got)
	}
}

func TestFormatPopularity(t *testing.T) {
	tests := map[float64]string{1.234: "1.23", 12.34: "12.3", 123.4: "123"}
	for value, want := range tests {
		if got := FormatPopularity(value); got != want {
			t.Fatalf("FormatPopularity(%v) = %q, want %q", value, got, want)
		}
	}
}

func TestWrapLine(t *testing.T) {
	lines := WrapLine("deps", "alpha beta gamma", 14)
	if len(lines) < 2 {
		t.Fatalf("WrapLine() = %#v, want wrapped lines", lines)
	}
}
