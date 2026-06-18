package sources

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestAURRPCSourceSearchesMaintainerQueries(t *testing.T) {
	alice := "Alice Developer"
	bob := "Bob Developer"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("by") != string(aur.SearchByMaintainer) || query.Get("arg") != "ali" {
			t.Fatalf("unexpected AUR RPC query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"version":5,"type":"search","resultcount":2,"results":[` +
			`{"Name":"alice-tool","PackageBase":"alice-tool","Maintainer":"` + alice + `"},` +
			`{"Name":"bob-tool","PackageBase":"bob-tool","Maintainer":"` + bob + `"}` +
			`]}`))
	}))
	defer server.Close()

	source := NewAURRPCSource("aur", server.URL)
	results, err := source.Search(context.Background(), "dev:ali")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 || results[0].Name != "alice-tool" || results[0].MaintainerName() != alice {
		t.Fatalf("Search() = %#v, want only Alice maintainer match", results)
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
