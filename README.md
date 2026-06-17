# aurview [![CI](https://github.com/kristyancarvalho/aurview/actions/workflows/ci.yml/badge.svg)](https://github.com/kristyancarvalho/aurview/actions/workflows/ci.yml) [![Release workflow](https://github.com/kristyancarvalho/aurview/actions/workflows/release.yml/badge.svg)](https://github.com/kristyancarvalho/aurview/actions/workflows/release.yml) [![License: MIT](https://img.shields.io/github/license/kristyancarvalho/aurview)](LICENSE) [![Go version](https://img.shields.io/github/go-mod/go-version/kristyancarvalho/aurview)](go.mod) [![Latest release](https://img.shields.io/github/v/release/kristyancarvalho/aurview)](https://github.com/kristyancarvalho/aurview/releases) [![Active milestones](https://img.shields.io/badge/milestones-active-2563eb)](https://github.com/kristyancarvalho/aurview/milestones)

`aurview` is a read-only terminal UI for searching AUR package metadata, ranking results, inspecting details and copying package names without leaving the command line.

It does not install, download, clone, build, update, remove or execute packages. It never calls `pacman`, `yay`, `paru`, `makepkg`, AUR helpers or `git clone` at runtime.

![aurview demo](.github/assets/demo.gif)

## Features

- Read-only AUR RPC search
- Local relevance ranking
- Package details for versions, votes, popularity, dependencies, licenses and URLs
- Keyboard-first navigation with optional mouse support
- Search history and clipboard copy
- Configurable AUR-compatible metadata sources
- Built-in `arch`, `mono`, `dark`, `light` and `high-contrast` themes
- Build-time version, commit and date metadata

## Installation

### Arch Linux

From AUR:

```sh
paru -S aurview
```

or

```sh
yay -S aurview
```

The AUR package metadata lives in `packaging/aur`.

```sh
cd packaging/aur
makepkg -si
```

### From Source

Requirements:

- Go 1.26 or newer
- Git

```sh
git clone https://github.com/kristyancarvalho/aurview.git
cd aurview
make build
install -Dm755 bin/aurview ~/.local/bin/aurview
```

### With Go

```sh
go install github.com/kristyancarvalho/aurview/cmd/aurview@latest
```

This installs the binary into `$GOBIN`, or `$GOPATH/bin` when `GOBIN` is not set.

## Usage

```sh
aurview
aurview paru
aurview "wayland screenshot"
aurview --version
```

Type a package name or keyword, select a result, then press `Enter` to copy the package name.

## Keybindings

| Key | Action |
|-----|--------|
| `/` | Focus search |
| `j` / `k` | Move selection or scroll details |
| `h` / `l` | Move focus between search, list and detail |
| `gg` / `G` | Go to top or bottom |
| `Ctrl+d` / `Ctrl+u` | Half page down or up |
| `Ctrl+f` / `Ctrl+b` | Page down or up |
| `n` / `N` | Next or previous search history entry |
| `Enter` | Copy selected package name |
| `Esc` | Blur search or close overlay |
| `?` | Toggle help |
| `q` | Quit |

## Configuration

aurview works without a config file. It loads `./aurview.toml` first, then `~/.config/aurview/config.toml`.

```toml
default_sources = ["aur"]

[ui]
theme = "arch"

[[sources]]
name = "aur"
type = "aur-rpc"
enabled = true
url = "https://aur.archlinux.org/rpc"
```

## Development

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
go build ./...
```

## Project Planning

Active work is tracked in [GitHub milestones](https://github.com/kristyancarvalho/aurview/milestones). `dev` is the integration branch, `main` is the stable release branch, and staging branches stay local-only.

## Releases

Releases are created from `v*` tags on `main` by the GitHub Actions workflow in `.github/workflows/release.yml`. Release notes live in [`docs/release-notes.md`](docs/release-notes.md).

## Contributing

Contributions are welcome. Start with [`CONTRIBUTING.md`](CONTRIBUTING.md), keep runtime behavior read-only, and open pull requests against `dev`.

## License

`aurview` is released under the [MIT License](LICENSE).
