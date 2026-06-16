# aurview

`aurview` is a read-only terminal UI for browsing, searching, ranking and inspecting AUR packages.

The project never installs, downloads, clones, builds, updates or removes packages. It only queries the official AUR RPC interface, previews metadata and copies selected package names to the clipboard.

## Status

Initial implementation is in progress. The core RPC, ranking, history and clipboard packages are implemented first; the custom TUI lands in the second staging branch.

## Data Source

`aurview` uses the official AUR RPC v5 interface:

- <https://wiki.archlinux.org/title/Aurweb_RPC_interface>
- <https://aur.archlinux.org/rpc>

Search uses `type=search` with `by=name-desc` by default. Detail lookup uses `type=info`.

## Development

```sh
go test ./...
go vet ./...
go run ./cmd/aurview
```

## Planned Keybindings

| Key | Action |
| --- | --- |
| `/` | Focus search |
| `j` / `k` | Move selection down/up |
| `h` / `l` | Move focus between panes |
| `gg` / `G` | Top/bottom |
| `Ctrl+d` / `Ctrl+u` | Half-page down/up |
| `Ctrl+f` / `Ctrl+b` | Page down/up |
| `n` / `N` | Navigate history or matches |
| `Enter` | Copy selected package name |
| `Esc` | Blur input or close overlay |
| `?` | Toggle help |
| `q` | Quit |

## Limitations

- Linux clipboard support is provider based and depends on `wl-copy`, `xclip` or `xsel`.
- Persistent history uses XDG state paths on Linux.
- The app is intentionally read-only and does not integrate with AUR helpers or package managers.
