package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/kristyancarvalho/aurview/internal/aur"
)

type Config struct {
	DefaultSources []string       `toml:"default_sources"`
	UI             UIConfig       `toml:"ui"`
	Sources        []SourceConfig `toml:"sources"`
	Path           string         `toml:"-"`
}

type UIConfig struct {
	Theme string `toml:"theme"`
}

type SourceConfig struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Enabled *bool  `toml:"enabled"`
	URL     string `toml:"url"`
}

func Default() Config {
	enabled := true
	return Config{
		DefaultSources: []string{"aur"},
		UI:             UIConfig{Theme: "arch"},
		Sources: []SourceConfig{{
			Name:    "aur",
			Type:    "aur-rpc",
			Enabled: &enabled,
			URL:     aur.DefaultBaseURL,
		}},
	}
}

func Load() (Config, error) {
	return LoadFromPaths(LocalConfigPath(), ConfigPath())
}

func LoadFromPaths(paths ...string) (Config, error) {
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		cfg, err := LoadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
	}
	return Default(), nil
}

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	cfg := Default()
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	cfg.Path = path
	if err := cfg.Normalize(); err != nil {
		return Config{}, fmt.Errorf("invalid config %s: %w", path, err)
	}
	return cfg, nil
}

func (c *Config) Normalize() error {
	if len(c.Sources) == 0 {
		c.Sources = Default().Sources
	}
	seen := map[string]bool{}
	for i := range c.Sources {
		source := &c.Sources[i]
		source.Name = strings.TrimSpace(source.Name)
		source.Type = strings.TrimSpace(source.Type)
		source.URL = strings.TrimSpace(source.URL)
		if source.Name == "" {
			return errors.New("source name is required")
		}
		key := strings.ToLower(source.Name)
		if seen[key] {
			return fmt.Errorf("duplicate source %q", source.Name)
		}
		seen[key] = true
		if source.Type == "" {
			source.Type = "aur-rpc"
		}
		if source.Type != "aur-rpc" {
			return fmt.Errorf("source %q uses unsupported type %q", source.Name, source.Type)
		}
		if source.URL == "" {
			if strings.EqualFold(source.Name, "aur") {
				source.URL = aur.DefaultBaseURL
			} else {
				return fmt.Errorf("source %q requires url", source.Name)
			}
		}
	}
	if len(c.DefaultSources) == 0 {
		for _, source := range c.Sources {
			if source.IsEnabled() {
				c.DefaultSources = append(c.DefaultSources, source.Name)
			}
		}
	}
	if c.UI.Theme == "" {
		c.UI.Theme = Default().UI.Theme
	}
	return nil
}

func (s SourceConfig) IsEnabled() bool {
	return s.Enabled == nil || *s.Enabled
}

func (c Config) EnabledSources() []SourceConfig {
	allowed := map[string]bool{}
	for _, name := range c.DefaultSources {
		allowed[strings.ToLower(strings.TrimSpace(name))] = true
	}
	out := make([]SourceConfig, 0, len(c.Sources))
	for _, source := range c.Sources {
		if !source.IsEnabled() {
			continue
		}
		if len(allowed) > 0 && !allowed[strings.ToLower(source.Name)] {
			continue
		}
		out = append(out, source)
	}
	return out
}
