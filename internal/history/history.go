package history

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const DefaultLimit = 100

type Store struct {
	limit int
	items []string
	pos   int
}

func New(limit int) *Store {
	if limit <= 0 {
		limit = DefaultLimit
	}
	return &Store{limit: limit, pos: -1}
}

func (s *Store) Add(query string) bool {
	query = strings.TrimSpace(query)
	if query == "" {
		return false
	}
	for i, item := range s.items {
		if item == query {
			copy(s.items[i:], s.items[i+1:])
			s.items[len(s.items)-1] = query
			s.pos = len(s.items)
			return false
		}
	}
	s.items = append(s.items, query)
	if len(s.items) > s.limit {
		s.items = append([]string(nil), s.items[len(s.items)-s.limit:]...)
	}
	s.pos = len(s.items)
	return true
}

func (s *Store) Prev() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	if s.pos < 0 || s.pos > len(s.items) {
		s.pos = len(s.items)
	}
	if s.pos > 0 {
		s.pos--
	}
	return s.items[s.pos], true
}

func (s *Store) Next() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	if s.pos < len(s.items)-1 {
		s.pos++
		return s.items[s.pos], true
	}
	s.pos = len(s.items)
	return "", false
}

func (s *Store) Items() []string {
	return append([]string(nil), s.items...)
}

func (s *Store) Load(path string) error {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.Add(scanner.Text())
	}
	s.pos = len(s.items)
	return scanner.Err()
}

func (s *Store) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range s.items {
		if _, err := writer.WriteString(item + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
