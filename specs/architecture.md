# Architecture

aurview keeps metadata logic separate from terminal rendering.

## Packages

- `cmd/aurview`: process entrypoint.
- `internal/app`: dependency wiring and program startup.
- `internal/aur`: AUR RPC client and package model.
- `internal/sources`: package source abstraction and multi-source search.
- `internal/config`: XDG/local config loading.
- `internal/ranking`: local relevance scoring.
- `internal/history`: in-session and persistent search history.
- `internal/clipboard`: Linux clipboard provider fallback.
- `internal/tui`: Bubble Tea model, update loop and custom renderer.
- `internal/tui/theme`: named theme definitions.
- `internal/tui/keymap`: keyboard action resolution.
- `internal/tui/components`: rendering helpers.
- `internal/platform`: platform/date helpers.

## Boundary Rule

Package-management behavior is out of scope. No package may add install, download, clone, build, update or remove actions.
