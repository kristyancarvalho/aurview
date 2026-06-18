package sources

import (
	"context"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/filter"
)

type AURRPCSource struct {
	name   string
	url    string
	client *aur.Client
}

func NewAURRPCSource(name, url string) *AURRPCSource {
	client := aur.NewClient(nil)
	if url != "" {
		client = client.WithBaseURL(url)
	}
	return &AURRPCSource{name: name, url: url, client: client}
}

func (s *AURRPCSource) Name() string {
	return s.name
}

func (s *AURRPCSource) Type() string {
	return "aur-rpc"
}

func (s *AURRPCSource) Search(ctx context.Context, query string) ([]aur.Package, error) {
	parsed := filter.ParseQuery(query)
	searchBy := aur.SearchByNameDesc
	searchText := parsed.SourceSearchText()
	if parsed.Text == "" && parsed.HasDeveloper() {
		searchBy = aur.SearchByMaintainer
	}
	pkgs, err := s.client.Search(ctx, searchText, searchBy)
	if err != nil {
		return nil, err
	}
	out := make([]aur.Package, 0, len(pkgs))
	for i := range pkgs {
		s.stamp(&pkgs[i])
		if parsed.MatchDeveloper(pkgs[i]) {
			out = append(out, pkgs[i])
		}
	}
	return out, nil
}

func (s *AURRPCSource) Info(ctx context.Context, name string) (aur.Package, error) {
	pkgs, err := s.client.Info(ctx, name)
	if err != nil {
		return aur.Package{}, err
	}
	if len(pkgs) == 0 {
		return aur.Package{}, aur.ErrEmptyQuery
	}
	pkg := pkgs[0].Clone()
	s.stamp(&pkg)
	return pkg, nil
}

func (s *AURRPCSource) stamp(pkg *aur.Package) {
	pkg.Source = s.name
	pkg.SourceType = s.Type()
	pkg.SourceURL = s.url
}
