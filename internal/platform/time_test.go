package platform

import "testing"

func TestUnixDateHelpers(t *testing.T) {
	if got := UnixDate(0); got != "-" {
		t.Fatalf("UnixDate(0) = %q", got)
	}
	if got := UnixDate(1_700_000_000); got == "-" || len(got) != len("2006-01-02") {
		t.Fatalf("UnixDate(valid) = %q", got)
	}
	v := int64(1_700_000_000)
	if got := OptionalUnixDate(&v); got == "-" {
		t.Fatalf("OptionalUnixDate(valid) = %q", got)
	}
}
