package app

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/clipboard"
	"github.com/kristyancarvalho/aurview/internal/config"
	"github.com/kristyancarvalho/aurview/internal/history"
	"github.com/kristyancarvalho/aurview/internal/tui"
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

	model := tui.New(tui.Options{
		Client:       aur.NewClient(nil),
		Copier:       clipboard.NewLinuxCopier(),
		History:      historyStore,
		InitialQuery: strings.TrimSpace(opts.InitialQuery),
	})

	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithContext(ctx))
	_, err := program.Run()
	return err
}
