package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromPathsDefaultsWhenMissing(t *testing.T) {
	cfg, err := LoadFromPaths(filepath.Join(t.TempDir(), "missing.toml"))
	if err != nil {
		t.Fatalf("LoadFromPaths() error = %v", err)
	}
	if got := cfg.EnabledSources(); len(got) != 1 || got[0].Name != "aur" {
		t.Fatalf("EnabledSources() = %#v, want default aur", got)
	}
}

func TestLoadFileParsesSourcesAndDisabledSources(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	writeConfig(t, path, `
default_sources = ["aur", "custom"]

[ui]
theme = "mono"

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

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if cfg.UI.Theme != "mono" {
		t.Fatalf("theme = %q, want mono", cfg.UI.Theme)
	}
	if got := cfg.EnabledSources(); len(got) != 1 || got[0].Name != "aur" {
		t.Fatalf("EnabledSources() = %#v, want only aur", got)
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

	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile() error = nil, want invalid source error")
	}
}

func writeConfig(t *testing.T, path, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
