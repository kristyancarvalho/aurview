package keymap

import "testing"

func TestResolverVimMotions(t *testing.T) {
	var resolver Resolver
	tests := []struct {
		key     string
		editing bool
		want    Action
	}{
		{key: "j", want: ActionDown},
		{key: "k", want: ActionUp},
		{key: "/", want: ActionSearch},
		{key: "ctrl+d", want: ActionHalfDown},
		{key: "G", want: ActionBottom},
		{key: "?", want: ActionHelp},
		{key: "f", want: ActionFilter},
		{key: "tab", want: ActionNextFilter},
		{key: " ", want: ActionToggleFilter},
		{key: "r", want: ActionResetFilters},
		{key: "q", want: ActionQuit},
		{key: "esc", editing: true, want: ActionBlur},
		{key: "tab", editing: true, want: ActionFilter},
		{key: "j", editing: true, want: ActionNone},
	}

	for _, tt := range tests {
		if got := resolver.Resolve(tt.key, tt.editing); got != tt.want {
			t.Fatalf("Resolve(%q, %v) = %v, want %v", tt.key, tt.editing, got, tt.want)
		}
	}
}

func TestResolverGG(t *testing.T) {
	var resolver Resolver
	if got := resolver.Resolve("g", false); got != ActionNone {
		t.Fatalf("first g = %v, want none", got)
	}
	if got := resolver.Resolve("g", false); got != ActionTop {
		t.Fatalf("second g = %v, want top", got)
	}
}
