# Add configurable repository sources

## Context

aurview currently searches AUR only. Users should be able to declare additional AUR RPC-compatible metadata sources while keeping AUR as the default.

## Goal

Load package sources from config and search enabled sources.

## Tasks

- [ ] Add XDG config path support for `~/.config/aurview/config.toml`.
- [ ] Support optional local `./aurview.toml`.
- [ ] Add default config behavior with AUR enabled.
- [ ] Add a `Source` interface with `Name`, `Type` and `Search`.
- [ ] Implement AUR RPC source adapter.
- [ ] Ignore disabled sources.

## Acceptance Criteria

- [ ] App starts without a config file and searches AUR.
- [ ] Enabled sources are searched.
- [ ] Disabled sources are ignored.
- [ ] Invalid config returns a helpful error.

## Test / Validation Notes

- Cover config loading, defaults, disabled sources and source parsing.
