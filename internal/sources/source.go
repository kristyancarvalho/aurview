package sources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/config"
)

type Source interface {
	Name() string
	Type() string
	Search(ctx context.Context, query string) ([]aur.Package, error)
	Info(ctx context.Context, name string) (aur.Package, error)
}

type MultiClient struct {
	sources []Source
}

func NewMultiClient(sources []Source) *MultiClient {
	return &MultiClient{sources: append([]Source(nil), sources...)}
}

func FromConfig(cfg config.Config) (*MultiClient, error) {
	sourceConfigs := cfg.EnabledSources()
	if len(sourceConfigs) == 0 {
		return nil, errors.New("no enabled package sources")
	}
	out := make([]Source, 0, len(sourceConfigs))
	for _, source := range sourceConfigs {
		switch source.Type {
		case config.SourceTypeAURRPC:
			out = append(out, NewAURRPCSource(source.Name, source.URL))
		case config.SourceTypePacmanSyncDB:
			out = append(out, NewPacmanSyncDBSource(source.Name, source.Repo, source.DBPath))
		default:
			return nil, fmt.Errorf("unsupported source type %q", source.Type)
		}
	}
	return NewMultiClient(out), nil
}

func (m *MultiClient) Search(ctx context.Context, query string) ([]aur.Package, error) {
	var all []aur.Package
	var errs []string
	for _, source := range m.sources {
		pkgs, err := source.Search(ctx, query)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", source.Name(), err))
			continue
		}
		all = append(all, pkgs...)
	}
	if len(all) == 0 && len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}
	return all, nil
}

func (m *MultiClient) Info(ctx context.Context, sourceName, name string) ([]aur.Package, error) {
	for _, source := range m.sources {
		if strings.EqualFold(source.Name(), sourceName) {
			pkg, err := source.Info(ctx, name)
			if err != nil {
				return nil, err
			}
			return []aur.Package{pkg}, nil
		}
	}
	return nil, fmt.Errorf("source %q not found", sourceName)
}

func (m *MultiClient) SourceCount() int {
	return len(m.sources)
}

func clonePackages(pkgs []aur.Package) []aur.Package {
	out := make([]aur.Package, len(pkgs))
	for i, pkg := range pkgs {
		out[i] = pkg.Clone()
	}
	return out
}
