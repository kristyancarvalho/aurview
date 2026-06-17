# Release Notes

## v0.4.1 - 2026-06-17

### Summary

Release maintenance update focused on release readiness, AUR packaging and simpler user-facing documentation.

### Added

- `aurview --version` output with build-time version, commit and date metadata
- Release notes in `docs/release-notes.md`

### Changed

- Simplified the README around installation, usage, keybindings, configuration and release flow
- Release workflow now validates release tags against `main`
- Release artifact actions were updated
- GitHub Release publication now checks out repository history before verifying tags
- AUR package metadata now targets `v0.4.1`

### Fixed

- Source archive generation no longer depends on broad development docs being present

### Validation Notes

- `gofmt -w .`
- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `makepkg --printsrcinfo`
- `makepkg --verifysource`

## v0.4.0 - 2026-06-16

### Summary

Initial read-only AUR browser release.

### Added

- AUR RPC package search and metadata inspection
- Local ranking, search history and clipboard copy support
- Mouse interaction, configurable metadata sources and built-in themes
- Release and AUR packaging metadata
