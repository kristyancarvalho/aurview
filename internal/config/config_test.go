package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadFromPathsDefaultsWhenMissing(t *testing.T) {
	cfg, err := LoadFromPathsWithDetector(noRepos, filepath.Join(t.TempDir(), "missing.toml"))
	if err != nil {
		t.Fatalf("LoadFromPaths() error = %v", err)
	}
	if got := cfg.EnabledSources(); len(got) != 1 || got[0].Name != "aur" {
		t.Fatalf("EnabledSources() = %#v, want default aur", got)
	}
}

func TestLoadFromPathsAddsDetectedLocalRepositories(t *testing.T) {
	cfg, err := LoadFromPathsWithDetector(staticRepos("core", "extra", "multilib", "chaotic-aur"), filepath.Join(t.TempDir(), "missing.toml"))
	if err != nil {
		t.Fatalf("LoadFromPathsWithDetector() error = %v", err)
	}

	if got, want := cfg.DefaultSources, []string{"aur", "core", "extra", "multilib", "chaotic-aur"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultSources = %#v, want %#v", got, want)
	}
	for _, repo := range []string{"core", "extra", "multilib", "chaotic-aur"} {
		source, ok := sourceByName(cfg.Sources, repo)
		if !ok {
			t.Fatalf("detected repo %q missing from sources: %#v", repo, cfg.Sources)
		}
		if source.Type != SourceTypePacmanSyncDB || source.Repo != repo || source.DBPath != PacmanSyncDBPath(repo) {
			t.Fatalf("source %q = %#v, want pacman sync db defaults", repo, source)
		}
	}
}

func TestLoadFileParsesSourcesAndDisabledSources(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	writeConfig(t, path, `
default_sources = ["aur", "custom"]

[ui]
theme = "mono"

[theme]
accent = "#112233"
selected_bg = "#445566"

[[sources]]
name = "aur"
type = "aur-rpc"
enabled = true
url = "https://aur.archlinux.org/rpc"

[[sources]]
name = "custom"
type = "aur-rpc"
enabled = false
url = "https://example.com/rpc"
`)

	cfg, err := LoadFileWithDetector(path, noRepos)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if cfg.UI.Theme != "mono" {
		t.Fatalf("theme = %q, want mono", cfg.UI.Theme)
	}
	if cfg.Theme.Accent != "#112233" || cfg.Theme.SelectedBG != "#445566" {
		t.Fatalf("theme colors = %#v, want parsed custom colors", cfg.Theme)
	}
	if got := cfg.EnabledSources(); len(got) != 1 || got[0].Name != "aur" {
		t.Fatalf("EnabledSources() = %#v, want only aur", got)
	}
}

func TestLoadFileRespectsDisabledDetectedSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	writeConfig(t, path, `
default_sources = ["aur", "core", "extra"]

[[sources]]
name = "aur"
type = "aur-rpc"
enabled = true

[[sources]]
name = "core"
type = "pacman-syncdb"
enabled = false
repo = "core"
`)

	cfg, err := LoadFileWithDetector(path, staticRepos("core", "extra"))
	if err != nil {
		t.Fatalf("LoadFileWithDetector() error = %v", err)
	}

	got := names(cfg.EnabledSources())
	want := []string{"aur", "extra"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EnabledSources() names = %#v, want %#v", got, want)
	}
}

func TestLoadFileRejectsInvalidSources(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	writeConfig(t, path, `
[[sources]]
name = "custom"
type = "unknown"
url = "https://example.com/rpc"
`)

	if _, err := LoadFileWithDetector(path, noRepos); err == nil {
		t.Fatal("LoadFile() error = nil, want invalid source error")
	}
}

func noRepos() ([]string, error) {
	return nil, nil
}

func staticRepos(repos ...string) RepoDetector {
	return func() ([]string, error) {
		return repos, nil
	}
}

func sourceByName(sources []SourceConfig, name string) (SourceConfig, bool) {
	for _, source := range sources {
		if source.Name == name {
			return source, true
		}
	}
	return SourceConfig{}, false
}

func names(sources []SourceConfig) []string {
	out := make([]string, 0, len(sources))
	for _, source := range sources {
		out = append(out, source.Name)
	}
	return out
}

func writeConfig(t *testing.T, path, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
