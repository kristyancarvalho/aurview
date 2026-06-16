package config

import (
	"os"
	"path/filepath"
)

const appName = "aurview"

func StateDir() string {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".local", "state", appName)
	}
	return filepath.Join(os.TempDir(), appName)
}

func ConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", appName)
	}
	return filepath.Join(os.TempDir(), appName)
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

func LocalConfigPath() string {
	return "aurview.toml"
}

func HistoryPath() string {
	return filepath.Join(StateDir(), "history")
}
