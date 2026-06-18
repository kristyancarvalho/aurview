package sources

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/filter"
)

type PacmanSyncDBSource struct {
	name   string
	repo   string
	dbPath string

	mu       sync.Mutex
	loaded   bool
	packages []aur.Package
	loadErr  error
}

func NewPacmanSyncDBSource(name, repo, dbPath string) *PacmanSyncDBSource {
	if repo == "" {
		repo = name
	}
	return &PacmanSyncDBSource{name: name, repo: repo, dbPath: dbPath}
}

func (s *PacmanSyncDBSource) Name() string {
	return s.name
}

func (s *PacmanSyncDBSource) Type() string {
	return "pacman-syncdb"
}

func (s *PacmanSyncDBSource) Search(_ context.Context, query string) ([]aur.Package, error) {
	parsed := filter.ParseQuery(query)
	searchText := strings.TrimSpace(parsed.Text)
	if searchText == "" && !parsed.HasDeveloper() {
		return nil, aur.ErrEmptyQuery
	}
	pkgs, err := s.loadPackages()
	if err != nil {
		return nil, err
	}
	terms := strings.Fields(strings.ToLower(searchText))
	out := make([]aur.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		if syncDBPackageMatches(pkg, terms) && parsed.MatchDeveloper(pkg) {
			out = append(out, pkg.Clone())
		}
	}
	return out, nil
}

func (s *PacmanSyncDBSource) Info(_ context.Context, name string) (aur.Package, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return aur.Package{}, aur.ErrEmptyQuery
	}
	pkgs, err := s.loadPackages()
	if err != nil {
		return aur.Package{}, err
	}
	for _, pkg := range pkgs {
		if strings.EqualFold(pkg.Name, name) {
			return pkg.Clone(), nil
		}
	}
	return aur.Package{}, fmt.Errorf("%s: package %q not found", s.name, name)
}

func (s *PacmanSyncDBSource) loadPackages() ([]aur.Package, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.loaded {
		s.packages, s.loadErr = readPacmanSyncDB(s.dbPath)
		for i := range s.packages {
			s.stamp(&s.packages[i])
		}
		sort.SliceStable(s.packages, func(i, j int) bool {
			return s.packages[i].Name < s.packages[j].Name
		})
		s.loaded = true
	}
	if s.loadErr != nil {
		return nil, s.loadErr
	}
	return clonePackages(s.packages), nil
}

func (s *PacmanSyncDBSource) stamp(pkg *aur.Package) {
	pkg.Source = s.name
	pkg.SourceType = s.Type()
	pkg.SourceURL = s.dbPath
	if pkg.PackageBase == "" {
		pkg.PackageBase = pkg.Name
	}
	if pkg.Maintainer == nil {
		maintainer := s.repo
		pkg.Maintainer = &maintainer
	}
}

func readPacmanSyncDB(dbPath string) ([]aur.Package, error) {
	reader, err := openPacmanSyncDB(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	var pkgs []aur.Package
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read pacman sync db %s: %w", dbPath, err)
		}
		if header.Typeflag != tar.TypeReg || path.Base(header.Name) != "desc" {
			continue
		}
		data, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, fmt.Errorf("read pacman sync db desc %s: %w", header.Name, err)
		}
		pkg, ok := parsePacmanDesc(data)
		if ok {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}

func openPacmanSyncDB(dbPath string) (io.ReadCloser, error) {
	file, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	buffered := bufio.NewReader(file)
	magic, _ := buffered.Peek(4)
	switch {
	case len(magic) >= 2 && magic[0] == 0x1f && magic[1] == 0x8b:
		gz, err := gzip.NewReader(buffered)
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		return readCloser{
			Reader: gz,
			close: func() error {
				return closeAll(gz.Close, file.Close)
			},
		}, nil
	case len(magic) >= 4 && bytes.Equal(magic[:4], []byte{0x28, 0xb5, 0x2f, 0xfd}):
		decoder, err := zstd.NewReader(buffered)
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		return readCloser{
			Reader: decoder,
			close: func() error {
				decoder.Close()
				return file.Close()
			},
		}, nil
	default:
		return readCloser{Reader: buffered, close: file.Close}, nil
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (r readCloser) Close() error {
	if r.close == nil {
		return nil
	}
	return r.close()
}

func closeAll(closers ...func() error) error {
	var first error
	for _, closer := range closers {
		if err := closer(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func parsePacmanDesc(data []byte) (aur.Package, bool) {
	sections := map[string][]string{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	section := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "%") && strings.HasSuffix(line, "%") {
			section = strings.Trim(line, "%")
			continue
		}
		if section != "" {
			sections[section] = append(sections[section], line)
		}
	}

	name := firstSectionValue(sections, "NAME")
	if name == "" {
		return aur.Package{}, false
	}
	pkg := aur.Package{
		Name:         name,
		PackageBase:  firstSectionValue(sections, "BASE"),
		Version:      firstSectionValue(sections, "VERSION"),
		Description:  firstSectionValue(sections, "DESC"),
		URL:          firstSectionValue(sections, "URL"),
		License:      sectionValues(sections, "LICENSE"),
		Depends:      sectionValues(sections, "DEPENDS"),
		MakeDepends:  sectionValues(sections, "MAKEDEPENDS"),
		CheckDepends: sectionValues(sections, "CHECKDEPENDS"),
		OptDepends:   sectionValues(sections, "OPTDEPENDS"),
		Conflicts:    sectionValues(sections, "CONFLICTS"),
		Provides:     sectionValues(sections, "PROVIDES"),
		LastModified: parseUnixSection(sections, "BUILDDATE"),
	}
	if pkg.PackageBase == "" {
		pkg.PackageBase = pkg.Name
	}
	if packager := firstSectionValue(sections, "PACKAGER"); packager != "" {
		pkg.Maintainer = &packager
	}
	return pkg, true
}

func firstSectionValue(sections map[string][]string, key string) string {
	values := sections[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func sectionValues(sections map[string][]string, key string) []string {
	return append([]string(nil), sections[key]...)
}

func parseUnixSection(sections map[string][]string, key string) int64 {
	value := firstSectionValue(sections, key)
	if value == "" {
		return 0
	}
	unix, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return unix
}

func syncDBPackageMatches(pkg aur.Package, terms []string) bool {
	if len(terms) == 0 {
		return true
	}
	haystack := strings.ToLower(strings.Join([]string{
		pkg.Name,
		pkg.PackageBase,
		pkg.Description,
		strings.Join(pkg.Provides, " "),
	}, " "))
	for _, term := range terms {
		if !strings.Contains(haystack, term) {
			return false
		}
	}
	return true
}
