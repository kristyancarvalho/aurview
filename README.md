# aurview

`aurview` is a lightweight, read-only terminal UI for browsing, searching, ranking and inspecting AUR packages.

The project never installs, downloads, clones, builds, updates or removes packages. It only queries the official AUR RPC interface, previews metadata and copies selected package names to the clipboard.

## Features

- Fast AUR RPC search by package name and description.
- Debounced non-blocking search with in-memory query cache.
- Local relevance ranking using name match, description match, votes, popularity, recency, out-of-date status and maintainer status.
- Compact custom TUI with result list, detail panel, status line and help overlay.
- Rich package detail preview via AUR RPC `info`.
- In-session and persistent search history using XDG state paths.
- Linux clipboard support through `wl-copy`, `xclip` or `xsel`.
- Vim-style movement keys.

## Install / Build

```sh
git clone https://github.com/kristyancarvalho/aurview.git
cd aurview
make build
./bin/aurview
```

You can also run without building:

```sh
make run
go run ./cmd/aurview paru
```

The optional command-line argument is used as the initial search query.

## Usage

Type a package name or keyword in the search prompt. Results are ranked locally and shown in a compact table. Move through the list, inspect the detail panel and press `Enter` to copy the selected package name.

No package-management actions exist in this tool.

## Keybindings

| Key | Action |
| --- | --- |
| `/` | Focus search |
| `j` / `k` | Move selection down/up, or scroll detail when detail is focused |
| `h` / `l` | Move focus between search, list and detail |
| `gg` / `G` | Top/bottom |
| `Ctrl+d` / `Ctrl+u` | Half-page down/up |
| `Ctrl+f` / `Ctrl+b` | Page down/up |
| `n` / `N` | Navigate search history newer/older |
| `Enter` | Copy selected package name |
| `Esc` | Blur search or close overlay |
| `?` | Toggle help overlay |
| `q` | Quit |

## Data Source

`aurview` uses the official AUR RPC v5 interface:

- <https://wiki.archlinux.org/title/Aurweb_RPC_interface>
- <https://aur.archlinux.org/rpc>

Search uses `type=search` with `by=name-desc` by default. Detail lookup uses `type=info`.

## Development

```sh
make fmt
make test
make lint
make build
```

Useful direct commands:

```sh
go test ./...
go vet ./...
go run ./cmd/aurview
```

## Architecture

- `cmd/aurview`: dependency wiring and app startup.
- `internal/app`: application assembly.
- `internal/aur`: read-only AUR RPC client and package models.
- `internal/ranking`: isolated relevance scoring.
- `internal/history`: in-memory and persistent search history.
- `internal/clipboard`: testable clipboard abstraction and Linux providers.
- `internal/tui`: Bubble Tea event loop, custom renderer and key handling.
- `internal/tui/components`: formatting helpers.
- `internal/tui/keymap`: vim-style action resolver.
- `internal/tui/theme`: ANSI visual theme with monochrome fallback.
- `internal/platform`: date conversion helpers.

## Screenshots

Screenshots will be added after the first tagged release.

## Clipboard

On Linux, `aurview` tries providers in this order:

1. `wl-copy`
2. `xclip -selection clipboard`
3. `xsel --clipboard --input`

If none are installed, the TUI stays usable and shows a clipboard warning when `Enter` is pressed.

## History

Search history is kept in memory during the session and saved to:

```text
$XDG_STATE_HOME/aurview/history
```

If `XDG_STATE_HOME` is unset, it falls back to:

```text
~/.local/state/aurview/history
```

## Testing

The test suite covers:

- AUR RPC response parsing and client error handling.
- Relevance ranking behavior.
- History navigation and persistence.
- Clipboard fallback behavior.
- Keymap action resolution.
- Formatting and date helpers.

## Limitations

- Linux clipboard support is provider based and depends on `wl-copy`, `xclip` or `xsel`.
- Persistent history uses XDG state paths on Linux.
- The app is intentionally read-only and does not integrate with AUR helpers or package managers.
- The TUI depends on terminal size; very small terminals show a size warning.
