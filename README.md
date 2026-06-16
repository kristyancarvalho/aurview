# aurview

[![CI](https://img.shields.io/github/actions/workflow/status/kristyancarvalho/aurview/test.yml?branch=dev&label=ci)](https://github.com/kristyancarvalho/aurview/actions/workflows/test.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/kristyancarvalho/aurview)](go.mod)
[![License: MIT](https://img.shields.io/github/license/kristyancarvalho/aurview)](LICENSE)
[![Release](https://img.shields.io/github/v/release/kristyancarvalho/aurview?include_prereleases)](https://github.com/kristyancarvalho/aurview/releases)
[![Issues](https://img.shields.io/github/issues/kristyancarvalho/aurview)](https://github.com/kristyancarvalho/aurview/issues)
[![Pull requests](https://img.shields.io/github/issues-pr/kristyancarvalho/aurview)](https://github.com/kristyancarvalho/aurview/pulls)
[![Last commit](https://img.shields.io/github/last-commit/kristyancarvalho/aurview/dev)](https://github.com/kristyancarvalho/aurview/commits/dev)
[![Stars](https://img.shields.io/github/stars/kristyancarvalho/aurview?style=social)](https://github.com/kristyancarvalho/aurview/stargazers)
[![Forks](https://img.shields.io/github/forks/kristyancarvalho/aurview?style=social)](https://github.com/kristyancarvalho/aurview/forks)
[![Go Report Card](https://goreportcard.com/badge/github.com/kristyancarvalho/aurview)](https://goreportcard.com/report/github.com/kristyancarvalho/aurview)

`aurview` is a terminal package browser for Arch Linux and AUR users. It helps you search package metadata, compare results, inspect package details and copy package names without leaving your terminal.

## Read-Only By Design

aurview does **not** install packages.

It also does **not** download, clone, build, update, remove or execute package scripts. It never calls `pacman`, `yay`, `paru`, `makepkg` or `git clone`.

This matters because AUR package management can execute arbitrary build scripts. aurview keeps discovery separate from installation so you can inspect metadata first and decide what to do somewhere else.

## What It Solves

Searching AUR in a browser is useful, but terminal users often want something faster and denser. aurview gives you:

- Fast AUR RPC search.
- Local relevance ranking.
- Compact package rows with votes, popularity, maintainer, update age and source.
- A detail panel with dependencies, licenses, URLs and status.
- Search history.
- Vim-style keyboard movement.
- Optional mouse selection and scrolling.
- Clipboard copying for package names.
- Configurable metadata sources.
- Built-in themes.

## Screenshot

Screenshots will be added after the first tagged release.

## Installation

Build from source:

```sh
git clone https://github.com/kristyancarvalho/aurview.git
cd aurview
make build
./bin/aurview
```

Run without building:

```sh
go run ./cmd/aurview
```

Start with an initial search:

```sh
go run ./cmd/aurview paru
make run ARGS=paru
```

## Basic Usage

Type a package name or keyword. aurview searches enabled read-only package metadata sources and ranks the results locally.

Use the result list to select a package. The detail panel shows package metadata. Press `Enter` to copy the selected package name to your clipboard.

Examples:

```text
paru
linux kernel
wayland screenshot
nerd font
```

## Keyboard Shortcuts

| Key | Action |
| --- | --- |
| `/` | Focus search |
| `j` / `k` | Move selection down/up, or scroll detail when detail is focused |
| `h` / `l` | Move focus between search, list and detail |
| `gg` / `G` | Jump to top/bottom |
| `Ctrl+d` / `Ctrl+u` | Half-page down/up |
| `Ctrl+f` / `Ctrl+b` | Page down/up |
| `n` / `N` | Navigate search history newer/older |
| `Enter` | Copy selected package name |
| `Esc` | Blur search or close overlay |
| `?` | Toggle help overlay |
| `q` | Quit |

## Mouse Support

Mouse support is optional and terminal-dependent.

- Click the search line to focus search.
- Click a package row to select it.
- Double-click a selected package row to copy the package name.
- Use the mouse wheel over the result list to move selection.
- Use the mouse wheel over the detail panel to scroll details.

The keyboard workflow remains the primary interface.

## Configuration

aurview works without a config file. If no config exists, it searches AUR only.

Preferred config path:

```text
~/.config/aurview/config.toml
```

Local project config is also supported:

```text
./aurview.toml
```

Example:

```toml
default_sources = ["aur"]

[ui]
theme = "arch"

[[sources]]
name = "aur"
type = "aur-rpc"
enabled = true
url = "https://aur.archlinux.org/rpc"

[[sources]]
name = "custom"
type = "aur-rpc"
enabled = false
url = "https://example.com/rpc"
```

## Package Sources

The default source is `AUR`, backed by the official AUR RPC interface.

Every result shows a source label. If two enabled sources return the same package name, aurview shows both rows so you can tell where each result came from.

Supported source types:

| Type | Description |
| --- | --- |
| `aur-rpc` | AUR RPC-compatible read-only metadata endpoint |

## Themes

Built-in themes:

- `arch`
- `mono`
- `dark`
- `light`
- `high-contrast`

Select a theme in config:

```toml
[ui]
theme = "high-contrast"
```

aurview respects `NO_COLOR` and falls back to readable plain text when color is unavailable.

## Clipboard

On Linux, aurview tries these clipboard providers:

1. `wl-copy`
2. `xclip -selection clipboard`
3. `xsel --clipboard --input`

If none are installed, aurview keeps running and shows a warning when copy is requested.

## Development

Common commands:

```sh
make fmt
make test
make lint
make build
```

Direct equivalents:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/aurview
```

Development notes live in [`/specs`](specs/).

## Project Status

aurview is an early MVP. The current focus is a polished read-only discovery workflow, clean source abstractions, and a custom terminal interface that does not look like a stock Bubble Tea demo.

## Roadmap

- `v0.2.0`: TUI interaction polish.
- `v0.3.0`: Multi-source package search.
- `v0.4.0`: Themes and documentation.
- `v1.0.0`: Stable read-only AUR browser.

More detail is tracked in [`specs/github/milestones.md`](specs/github/milestones.md).

## Contributing

Contributions should start from `dev`, use focused branches, and keep aurview read-only.

Before opening a pull request, run:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/aurview
```

See [`specs/github/contribution-workflow.md`](specs/github/contribution-workflow.md).

## License

aurview is released under the MIT License.

Copyright (c) 2026 aurview contributors.
