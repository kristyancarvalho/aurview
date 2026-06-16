package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/clipboard"
	"github.com/kristyancarvalho/aurview/internal/history"
	"github.com/kristyancarvalho/aurview/internal/ranking"
	"github.com/kristyancarvalho/aurview/internal/tui/components"
	"github.com/kristyancarvalho/aurview/internal/tui/keymap"
	"github.com/kristyancarvalho/aurview/internal/tui/theme"
)

const debounceDelay = 260 * time.Millisecond

type PackageClient interface {
	Search(ctx context.Context, query string) ([]aur.Package, error)
	Info(ctx context.Context, sourceName, name string) ([]aur.Package, error)
}

type Options struct {
	Client       PackageClient
	Copier       clipboard.Copier
	History      *history.Store
	InitialQuery string
}

type focusArea int

const (
	focusSearch focusArea = iota
	focusList
	focusDetail
)

type Model struct {
	client  PackageClient
	copier  clipboard.Copier
	history *history.Store
	theme   theme.Theme
	keys    keymap.Resolver
	scorer  ranking.Scorer

	width  int
	height int

	focus focusArea
	input string

	token       int
	loading     bool
	searchError string
	lastQuery   string

	results  []ranking.RankedPackage
	selected int
	scroll   int

	detailCache   map[string]aur.Package
	detailLoading bool
	detailError   string
	detailScroll  int

	help       bool
	status     string
	statusKind string

	lastClickName string
	lastClickAt   time.Time
}

func New(opts Options) Model {
	if opts.History == nil {
		opts.History = history.New(history.DefaultLimit)
	}
	return Model{
		client:      opts.Client,
		copier:      opts.Copier,
		history:     opts.History,
		theme:       theme.Detect(),
		scorer:      ranking.NewScorer(time.Now()),
		focus:       focusSearch,
		input:       opts.InitialQuery,
		detailCache: make(map[string]aur.Package),
		status:      "read-only: Enter copies package name",
		statusKind:  "info",
	}
}

func (m Model) Init() tea.Cmd {
	if strings.TrimSpace(m.input) == "" {
		return nil
	}
	return m.scheduleSearch()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ensureSelectionVisible()
		return m, nil
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.MouseMsg:
		return m.updateMouse(msg)
	case debounceMsg:
		if msg.token != m.token {
			return m, nil
		}
		query := strings.TrimSpace(msg.query)
		if query == "" {
			m.loading = false
			m.results = nil
			m.searchError = ""
			m.lastQuery = ""
			return m, nil
		}
		m.loading = true
		m.searchError = ""
		m.lastQuery = query
		return m, searchCmd(m.client, msg.token, query)
	case searchResultMsg:
		if msg.token != m.token {
			return m, nil
		}
		m.loading = false
		if msg.err != nil {
			m.results = nil
			m.searchError = userSearchError(msg.err)
			m.status = m.searchError
			m.statusKind = "error"
			return m, nil
		}
		m.history.Add(msg.query)
		m.results = m.scorer.Rank(msg.query, msg.packages)
		m.selected = 0
		m.scroll = 0
		m.detailScroll = 0
		m.searchError = ""
		if len(m.results) == 0 {
			m.status = "no packages matched " + msg.query
			m.statusKind = "warn"
			return m, nil
		}
		m.status = fmt.Sprintf("%d packages ranked for %q", len(m.results), msg.query)
		m.statusKind = "ok"
		return m, m.fetchSelectedDetail()
	case detailResultMsg:
		if msg.err != nil {
			if m.selectedKey() == msg.key() {
				m.detailLoading = false
				m.detailError = msg.err.Error()
			}
			return m, nil
		}
		m.detailCache[msg.key()] = msg.pkg.Clone()
		if m.selectedKey() == msg.key() {
			m.detailLoading = false
			m.detailError = ""
		}
		return m, nil
	case copyMsg:
		if msg.err != nil {
			m.status = "clipboard unavailable: " + msg.err.Error()
			m.statusKind = "warn"
			return m, nil
		}
		m.status = "copied " + msg.name
		m.statusKind = "ok"
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) updateMouse(msg tea.MouseMsg) (Model, tea.Cmd) {
	if m.help {
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			m.help = false
		}
		return m, nil
	}

	area := m.hitArea(msg.X, msg.Y)
	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionPress {
			return m, nil
		}
		switch area.kind {
		case hitSearch:
			m.focus = focusSearch
			m.status = "search focused"
			m.statusKind = "info"
		case hitListRow:
			m.focus = focusList
			m.selected = area.index
			m.detailScroll = 0
			m.ensureSelectionVisible()
			name := m.selectedName()
			doubleClick := name != "" && name == m.lastClickName && time.Since(m.lastClickAt) <= 450*time.Millisecond
			m.lastClickName = name
			m.lastClickAt = time.Now()
			if doubleClick {
				cmd := m.copySelected()
				return m, cmd
			}
			cmd := m.fetchSelectedDetail()
			return m, cmd
		case hitDetail:
			m.focus = focusDetail
		}
	case tea.MouseButtonWheelDown:
		return m.updateMouseWheel(area, 3)
	case tea.MouseButtonWheelUp:
		return m.updateMouseWheel(area, -3)
	}
	return m, nil
}

func (m Model) updateMouseWheel(area hitResult, delta int) (Model, tea.Cmd) {
	switch area.kind {
	case hitDetail:
		m.focus = focusDetail
		m.moveDetail(delta)
	case hitList, hitListRow:
		m.focus = focusList
		m.moveSelection(delta)
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case hitSearch:
		m.focus = focusSearch
	}
	return m, nil
}

func (m Model) updateKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.help {
		if msg.String() == "?" || msg.String() == "esc" || msg.String() == "q" {
			m.help = false
			return m, nil
		}
		return m, nil
	}

	if m.focus == focusSearch {
		if next, ok := m.editInput(msg); ok {
			m = next
			return m, m.scheduleSearch()
		}
	}

	action := m.keys.Resolve(msg.String(), m.focus == focusSearch)
	switch action {
	case keymap.ActionQuit:
		return m, tea.Quit
	case keymap.ActionHelp:
		m.help = !m.help
	case keymap.ActionSearch:
		m.focus = focusSearch
	case keymap.ActionBlur:
		if m.focus == focusSearch {
			m.focus = focusList
		}
	case keymap.ActionCopy:
		cmd := m.copySelected()
		return m, cmd
	case keymap.ActionDown:
		if m.focus == focusDetail {
			m.moveDetail(1)
			return m, nil
		}
		m.moveSelection(1)
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionUp:
		if m.focus == focusDetail {
			m.moveDetail(-1)
			return m, nil
		}
		m.moveSelection(-1)
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionLeft:
		if m.focus == focusDetail {
			m.focus = focusList
		} else {
			m.focus = focusSearch
		}
	case keymap.ActionRight:
		if m.focus == focusSearch {
			m.focus = focusList
		} else {
			m.focus = focusDetail
		}
	case keymap.ActionTop:
		if m.focus == focusDetail {
			m.detailScroll = 0
			return m, nil
		}
		m.selected, m.scroll = 0, 0
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionBottom:
		if m.focus == focusDetail {
			m.detailScroll = 9999
			return m, nil
		}
		if len(m.results) > 0 {
			m.selected = len(m.results) - 1
			m.ensureSelectionVisible()
		}
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionHalfDown:
		if m.focus == focusDetail {
			m.moveDetail(max(1, m.detailHeight()/2))
			return m, nil
		}
		m.moveSelection(max(1, m.listHeight()/2))
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionHalfUp:
		if m.focus == focusDetail {
			m.moveDetail(-max(1, m.detailHeight()/2))
			return m, nil
		}
		m.moveSelection(-max(1, m.listHeight()/2))
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionPageDown:
		if m.focus == focusDetail {
			m.moveDetail(max(1, m.detailHeight()))
			return m, nil
		}
		m.moveSelection(max(1, m.listHeight()))
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionPageUp:
		if m.focus == focusDetail {
			m.moveDetail(-max(1, m.detailHeight()))
			return m, nil
		}
		m.moveSelection(-max(1, m.listHeight()))
		cmd := m.fetchSelectedDetail()
		return m, cmd
	case keymap.ActionHistoryPrev:
		cmd := m.setHistory(m.history.Prev())
		return m, cmd
	case keymap.ActionHistoryNext:
		cmd := m.setHistory(m.history.Next())
		return m, cmd
	}
	return m, nil
}

func (m Model) editInput(msg tea.KeyMsg) (Model, bool) {
	switch msg.Type {
	case tea.KeyBackspace, tea.KeyCtrlH:
		if len(m.input) == 0 {
			return m, false
		}
		runes := []rune(m.input)
		m.input = string(runes[:len(runes)-1])
		m.token++
		m.loading = strings.TrimSpace(m.input) != ""
		return m, true
	case tea.KeySpace:
		m.input += " "
		m.token++
		m.loading = strings.TrimSpace(m.input) != ""
		return m, true
	case tea.KeyRunes:
		value := msg.String()
		if value == "?" || value == "/" {
			return m, false
		}
		m.input += value
		m.token++
		m.loading = strings.TrimSpace(m.input) != ""
		return m, true
	default:
		return m, false
	}
}

func (m *Model) setHistory(value string, ok bool) tea.Cmd {
	if !ok {
		m.status = "history boundary"
		m.statusKind = "info"
		return nil
	}
	m.input = value
	m.focus = focusSearch
	m.token++
	m.loading = true
	return m.scheduleSearch()
}

func (m Model) scheduleSearch() tea.Cmd {
	token := m.token
	query := m.input
	return tea.Tick(debounceDelay, func(time.Time) tea.Msg {
		return debounceMsg{token: token, query: query}
	})
}

func (m *Model) moveSelection(delta int) {
	if len(m.results) == 0 {
		return
	}
	m.selected = components.Clamp(m.selected+delta, 0, len(m.results)-1)
	m.detailScroll = 0
	m.ensureSelectionVisible()
}

func (m *Model) moveDetail(delta int) {
	m.detailScroll = max(0, m.detailScroll+delta)
}

func (m *Model) ensureSelectionVisible() {
	visible := m.listHeight()
	if visible <= 0 {
		return
	}
	if m.selected < m.scroll {
		m.scroll = m.selected
	}
	if m.selected >= m.scroll+visible {
		m.scroll = m.selected - visible + 1
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
}

func (m Model) listHeight() int {
	if m.height <= 0 {
		return 10
	}
	if m.width >= 110 {
		return max(1, m.height-5)
	}
	return max(1, (m.height-6)*2/3)
}

func (m Model) detailHeight() int {
	if m.height <= 0 {
		return 10
	}
	if m.width >= 110 {
		return max(1, m.height-5)
	}
	return max(1, m.height-m.listHeight()-6)
}

func (m Model) selectedName() string {
	if len(m.results) == 0 || m.selected < 0 || m.selected >= len(m.results) {
		return ""
	}
	return m.results[m.selected].Package.Name
}

func (m Model) selectedSource() string {
	if len(m.results) == 0 || m.selected < 0 || m.selected >= len(m.results) {
		return ""
	}
	source := m.results[m.selected].Package.Source
	if source == "" {
		return "aur"
	}
	return source
}

func (m Model) selectedKey() string {
	name := m.selectedName()
	if name == "" {
		return ""
	}
	return packageKey(m.selectedSource(), name)
}

func (m Model) selectedPackage() (aur.Package, bool) {
	key := m.selectedKey()
	if key == "" {
		return aur.Package{}, false
	}
	if pkg, ok := m.detailCache[key]; ok {
		return pkg.Clone(), true
	}
	return m.results[m.selected].Package.Clone(), true
}

func (m *Model) fetchSelectedDetail() tea.Cmd {
	name := m.selectedName()
	if name == "" {
		return nil
	}
	source := m.selectedSource()
	if _, ok := m.detailCache[packageKey(source, name)]; ok {
		m.detailLoading = false
		return nil
	}
	m.detailLoading = true
	return infoCmd(m.client, source, name)
}

func (m *Model) copySelected() tea.Cmd {
	name := m.selectedName()
	if name == "" {
		m.status = "nothing selected"
		m.statusKind = "warn"
		return nil
	}
	return copyCmd(m.copier, name)
}

func searchCmd(client PackageClient, token int, query string) tea.Cmd {
	return func() tea.Msg {
		pkgs, err := client.Search(context.Background(), query)
		return searchResultMsg{token: token, query: query, packages: pkgs, err: err}
	}
}

func infoCmd(client PackageClient, source, name string) tea.Cmd {
	return func() tea.Msg {
		pkgs, err := client.Info(context.Background(), source, name)
		if err != nil {
			return detailResultMsg{source: source, name: name, err: err}
		}
		if len(pkgs) == 0 {
			return detailResultMsg{source: source, name: name, err: errors.New("no detail returned")}
		}
		return detailResultMsg{source: source, name: name, pkg: pkgs[0]}
	}
}

func copyCmd(copier clipboard.Copier, name string) tea.Cmd {
	return func() tea.Msg {
		if copier == nil {
			return copyMsg{name: name, err: clipboard.ErrUnavailable}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return copyMsg{name: name, err: copier.Copy(ctx, name)}
	}
}

func userSearchError(err error) string {
	if errors.Is(err, aur.ErrRateLimit) {
		return "AUR RPC rate limit reached; pause before searching again"
	}
	if errors.Is(err, aur.ErrEmptyQuery) {
		return "empty search"
	}
	return "search failed: " + err.Error()
}

type debounceMsg struct {
	token int
	query string
}

type searchResultMsg struct {
	token    int
	query    string
	packages []aur.Package
	err      error
}

type detailResultMsg struct {
	source string
	name   string
	pkg    aur.Package
	err    error
}

func (m detailResultMsg) key() string {
	return packageKey(m.source, m.name)
}

type copyMsg struct {
	name string
	err  error
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func packageKey(source, name string) string {
	if source == "" {
		source = "aur"
	}
	return strings.ToLower(source) + "\x00" + strings.ToLower(name)
}

type hitKind int

const (
	hitNone hitKind = iota
	hitSearch
	hitList
	hitListRow
	hitDetail
)

type hitResult struct {
	kind  hitKind
	index int
}

func (m Model) hitArea(x, y int) hitResult {
	if y == 1 {
		return hitResult{kind: hitSearch}
	}
	if m.width >= 110 {
		leftWidth := max(62, m.width*58/100)
		if y >= 2 {
			if x <= leftWidth {
				return m.hitList(y, 2)
			}
			return hitResult{kind: hitDetail}
		}
		return hitResult{kind: hitNone}
	}
	listHeight := max(3, (m.height-5)*2/3)
	if y >= 2 && y <= 2+listHeight {
		return m.hitList(y, 2)
	}
	if y > 2+listHeight {
		return hitResult{kind: hitDetail}
	}
	return hitResult{kind: hitNone}
}

func (m Model) hitList(y, headerY int) hitResult {
	if y == headerY {
		return hitResult{kind: hitList}
	}
	row := y - headerY - 1
	visible := m.listHeight()
	if row < 0 || row >= visible {
		return hitResult{kind: hitList}
	}
	index := m.scroll + row
	if index < 0 || index >= len(m.results) {
		return hitResult{kind: hitList}
	}
	return hitResult{kind: hitListRow, index: index}
}
