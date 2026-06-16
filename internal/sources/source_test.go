package sources

import (
	"context"
	"reflect"
	"testing"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/config"
)

type fakeSource struct {
	name string
	pkgs []aur.Package
	err  error
}

func (f fakeSource) Name() string { return f.name }
func (f fakeSource) Type() string { return "fake" }
func (f fakeSource) Search(context.Context, string) ([]aur.Package, error) {
	out := make([]aur.Package, len(f.pkgs))
	copy(out, f.pkgs)
	for i := range out {
		out[i].Source = f.name
	}
	return out, f.err
}
func (f fakeSource) Info(context.Context, string) (aur.Package, error) {
	if len(f.pkgs) == 0 {
		return aur.Package{}, aur.ErrEmptyQuery
	}
	pkg := f.pkgs[0]
	pkg.Source = f.name
	return pkg, nil
}

func TestMultiClientKeepsDuplicateNamesAcrossSources(t *testing.T) {
	client := NewMultiClient([]Source{
		fakeSource{name: "aur", pkgs: []aur.Package{{Name: "paru"}}},
		fakeSource{name: "custom", pkgs: []aur.Package{{Name: "paru"}}},
	})

	pkgs, err := client.Search(context.Background(), "paru")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	got := []string{pkgs[0].DisplaySource(), pkgs[1].DisplaySource()}
	if want := []string{"AUR", "custom"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("sources = %#v, want %#v", got, want)
	}
}

func TestFromConfigUsesEnabledSources(t *testing.T) {
	disabled := false
	enabled := true
	client, err := FromConfig(config.Config{
		DefaultSources: []string{"aur", "custom"},
		Sources: []config.SourceConfig{
			{Name: "aur", Type: "aur-rpc", Enabled: &enabled, URL: aur.DefaultBaseURL},
			{Name: "custom", Type: "aur-rpc", Enabled: &disabled, URL: "https://example.com/rpc"},
		},
	})
	if err != nil {
		t.Fatalf("FromConfig() error = %v", err)
	}
	if got := client.SourceCount(); got != 1 {
		t.Fatalf("SourceCount() = %d, want 1", got)
	}
}
