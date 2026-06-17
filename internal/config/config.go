package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/kristyancarvalho/aurview/internal/aur"
)

type Config struct {
	DefaultSources []string       `toml:"default_sources"`
	UI             UIConfig       `toml:"ui"`
	Theme          ThemeConfig    `toml:"theme"`
	Sources        []SourceConfig `toml:"sources"`
	Path           string         `toml:"-"`
}

type UIConfig struct {
	Theme string `toml:"theme"`
}

type ThemeConfig struct {
	Accent      string `toml:"accent"`
	Good        string `toml:"good"`
	Warn        string `toml:"warn"`
	Danger      string `toml:"danger"`
	Muted       string `toml:"muted"`
	Dim         string `toml:"dim"`
	Focus       string `toml:"focus"`
	SelectedFG  string `toml:"selected_fg"`
	SelectedBG  string `toml:"selected_bg"`
	BadgeFG     string `toml:"badge_fg"`
	BadgeBG     string `toml:"badge_bg"`
	HeaderFG    string `toml:"header_fg"`
	HeaderBG    string `toml:"header_bg"`
	FilterFG    string `toml:"filter_fg"`
	FilterBG    string `toml:"filter_bg"`
	FilterOnFG  string `toml:"filter_on_fg"`
	FilterOnBG  string `toml:"filter_on_bg"`
	FilterHotFG string `toml:"filter_hot_fg"`
	FilterHotBG string `toml:"filter_hot_bg"`
}

type SourceConfig struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Enabled *bool  `toml:"enabled"`
	URL     string `toml:"url"`
	Repo    string `toml:"repo"`
	DBPath  string `toml:"db_path"`
}

type RepoDetector func() ([]string, error)

const (
	SourceTypeAURRPC       = "aur-rpc"
	SourceTypePacmanSyncDB = "pacman-syncdb"
)

func Default() Config {
	return defaultWithRepos(nil)
}

func defaultWithRepos(repos []string) Config {
	enabled := true
	cfg := Config{
		DefaultSources: []string{"aur"},
		UI:             UIConfig{Theme: "arch"},
		Sources: []SourceConfig{{
			Name:    "aur",
			Type:    SourceTypeAURRPC,
			Enabled: &enabled,
			URL:     aur.DefaultBaseURL,
		}},
	}
	cfg.addDetectedRepos(repos)
	cfg.DefaultSources = enabledSourceNames(cfg.Sources)
	return cfg
}

func Load() (Config, error) {
	return LoadFromPathsWithDetector(DetectPacmanRepositories, LocalConfigPath(), ConfigPath())
}

func LoadFromPaths(paths ...string) (Config, error) {
	return LoadFromPathsWithDetector(DetectPacmanRepositories, paths...)
}

func LoadFromPathsWithDetector(detector RepoDetector, paths ...string) (Config, error) {
	repos := detectRepos(detector)
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		cfg, err := LoadFileWithRepos(path, repos)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
	}
	return defaultWithRepos(repos), nil
}

func LoadFile(path string) (Config, error) {
	return LoadFileWithDetector(path, DetectPacmanRepositories)
}

func LoadFileWithDetector(path string, detector RepoDetector) (Config, error) {
	return LoadFileWithRepos(path, detectRepos(detector))
}

func LoadFileWithRepos(path string, repos []string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	cfg.Path = path
	if err := cfg.NormalizeWithRepos(repos); err != nil {
		return Config{}, fmt.Errorf("invalid config %s: %w", path, err)
	}
	return cfg, nil
}

func (c *Config) Normalize() error {
	return c.NormalizeWithRepos(nil)
}

func (c *Config) NormalizeWithRepos(repos []string) error {
	if len(c.Sources) == 0 {
		c.Sources = defaultWithRepos(repos).Sources
	} else {
		c.addDetectedRepos(repos)
	}
	seen := map[string]bool{}
	for i := range c.Sources {
		source := &c.Sources[i]
		source.Name = strings.TrimSpace(source.Name)
		source.Type = strings.TrimSpace(source.Type)
		source.URL = strings.TrimSpace(source.URL)
		source.Repo = strings.TrimSpace(source.Repo)
		source.DBPath = strings.TrimSpace(source.DBPath)
		if source.Name == "" {
			return errors.New("source name is required")
		}
		key := strings.ToLower(source.Name)
		if seen[key] {
			return fmt.Errorf("duplicate source %q", source.Name)
		}
		seen[key] = true
		if source.Type == "" {
			source.Type = SourceTypeAURRPC
		}
		switch source.Type {
		case SourceTypeAURRPC:
			if source.URL == "" {
				if strings.EqualFold(source.Name, "aur") {
					source.URL = aur.DefaultBaseURL
				} else {
					return fmt.Errorf("source %q requires url", source.Name)
				}
			}
		case SourceTypePacmanSyncDB:
			if source.Repo == "" {
				source.Repo = source.Name
			}
			if source.DBPath == "" {
				source.DBPath = PacmanSyncDBPath(source.Repo)
			}
		default:
			return fmt.Errorf("source %q uses unsupported type %q", source.Name, source.Type)
		}
	}
	if len(c.DefaultSources) == 0 {
		c.DefaultSources = enabledSourceNames(c.Sources)
	}
	if c.UI.Theme == "" {
		c.UI.Theme = Default().UI.Theme
	}
	return nil
}

func (c *Config) addDetectedRepos(repos []string) {
	if len(repos) == 0 {
		return
	}
	existing := map[string]bool{}
	for _, source := range c.Sources {
		existing[strings.ToLower(strings.TrimSpace(source.Name))] = true
	}
	enabled := true
	for _, repo := range repos {
		repo = strings.TrimSpace(repo)
		if repo == "" {
			continue
		}
		key := strings.ToLower(repo)
		if existing[key] {
			continue
		}
		c.Sources = append(c.Sources, SourceConfig{
			Name:    repo,
			Type:    SourceTypePacmanSyncDB,
			Enabled: &enabled,
			Repo:    repo,
			DBPath:  PacmanSyncDBPath(repo),
		})
		existing[key] = true
	}
}

func enabledSourceNames(sources []SourceConfig) []string {
	names := make([]string, 0, len(sources))
	for _, source := range sources {
		if source.IsEnabled() {
			names = append(names, source.Name)
		}
	}
	return names
}

func detectRepos(detector RepoDetector) []string {
	if detector == nil {
		return nil
	}
	repos, err := detector()
	if err != nil {
		return nil
	}
	return uniqueTrimmed(repos)
}

func DetectPacmanRepositories() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pacman-conf", "--repo-list")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}
	return parseRepoList(string(out)), nil
}

func PacmanSyncDBPath(repo string) string {
	return filepath.Join("/var/lib/pacman/sync", repo+".db")
}

func parseRepoList(out string) []string {
	return uniqueTrimmed(strings.Split(out, "\n"))
}

func uniqueTrimmed(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
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
