# Configuration

aurview runs without configuration. The default config searches AUR only.

Search order:

1. `./aurview.toml`
2. `~/.config/aurview/config.toml`
3. built-in defaults

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
```

Only read-only metadata sources are supported.
