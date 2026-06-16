package app

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/aurview/internal/clipboard"
	"github.com/kristyancarvalho/aurview/internal/config"
	"github.com/kristyancarvalho/aurview/internal/history"
	"github.com/kristyancarvalho/aurview/internal/sources"
	"github.com/kristyancarvalho/aurview/internal/tui"
	"github.com/kristyancarvalho/aurview/internal/tui/theme"
)

type Options struct {
	InitialQuery string
}

func Run(ctx context.Context, opts Options) error {
	historyStore := history.New(history.DefaultLimit)
	historyPath := config.HistoryPath()
	if err := historyStore.Load(historyPath); err != nil {
		return fmt.Errorf("load history: %w", err)
	}
	defer func() {
		_ = historyStore.Save(historyPath)
	}()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	sourceClient, err := sources.FromConfig(cfg)
	if err != nil {
		return err
	}
	selectedTheme, err := theme.Detect(cfg.UI.Theme)
	if err != nil {
		return err
	}

	model := tui.New(tui.Options{
		Client:       sourceClient,
		Copier:       clipboard.NewLinuxCopier(),
		History:      historyStore,
		InitialQuery: strings.TrimSpace(opts.InitialQuery),
		Theme:        selectedTheme,
	})

	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithContext(ctx))
	_, err = program.Run()
	return err
}
