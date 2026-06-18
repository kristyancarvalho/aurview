# Release Notes

## v0.5.2 - 2026-06-17

### Summary

Patch release improving TUI scanability with better theme color distribution, distinct repository source badges and developer/maintainer query filters.

### Added

- Query filters for package developers and maintainers:
  - `dev:<name>`
  - `developer:<name>`
  - `maint:<name>`
  - `maintainer:<name>`
- AUR maintainer search through the AUR RPC maintainer search mode
- Local pacman repository packager search using `%PACKAGER%` metadata from sync databases
- Repository source badges with stable labels for `AUR`, `CORE`, `EXT`, `MULTI`, `CHAOTIC` and custom repositories

### Changed

- Matugen theme mapping now keeps primary colors focused on app/header identity while using secondary, tertiary and surface roles for selected rows, filters and badges
- Built-in themes now use more distinct filled styles for headers, selected rows, badges and filter states
- Source badge styles are repository-specific and keep selected-row contrast readable
- The maintained/orphaned filter chip is now labeled `maint-state` to avoid confusion with `maint:<name>` query filters

### Known Limitations

- AUR maintainer-only queries depend on the read-only AUR RPC maintainer search endpoint.
- aurview still does not install, remove, upgrade, clone, build or mutate packages.

### Validation Notes

- `gofmt -w .`
- `go test ./...`
- `go vet ./...`
- `go build ./cmd/aurview`
- Search validation for normal package names and `dev:`, `developer:`, `maint:` and `maintainer:` query aliases
- Badge validation for `AUR`, `core`, `extra`, `multilib`, `chaotic-aur` and custom repositories
- Matugen and built-in theme visual validation

## v0.5.0 - 2026-06-17

### Summary

Feature release adding read-only local Arch repository search, Matugen-generated color themes and complete practical config documentation.

### Added

- Local pacman sync database sources for repositories detected by `pacman-conf --repo-list`
- Default search coverage for AUR plus detected local repositories such as `core`, `extra`, `multilib` and `chaotic-aur`
- Read-only parsing for `/var/lib/pacman/sync/<repo>.db`, including gzip and zstd databases
- Matugen color-only theme support with safe fallback for invalid or missing hex colors
- README configuration docs for paths, precedence, `[ui]`, `[theme]`, `[[sources]]`, source types and examples

### Changed

- Mixed AUR and local repository search results clearly preserve source labels in results, filters and details
- Missing local sync databases no longer prevent other configured sources from returning results

### Known Limitations

- aurview still does not install, remove, upgrade, clone, build or mutate packages.
- Local repository metadata depends on existing pacman sync databases; aurview does not refresh them.

### Validation Notes

- `gofmt -w .`
- `go test ./...`
- `go vet ./...`
- `go build ./cmd/aurview`
- Local search validation against AUR, `core`, `extra`, `multilib` and `chaotic-aur`
- Missing local sync database validation
- Matugen template validation

## v0.4.2 - 2026-06-17

### Summary

Patch release focused on trustworthy release metadata, interactive result filters and clearer TUI visual hierarchy.

### Added

- Search filters for source, maintained/orphaned packages, out-of-date flags, minimum votes, minimum popularity, recently updated packages and match mode
- Compact keyboard-friendly filter bar with reset support

### Changed

- Top headers, table headers and filter chips now use compact htop-like filled backgrounds
- Build metadata now has a single `internal/version` source and release builds override it with `-ldflags`
- AUR package metadata now targets `v0.4.2`

### Fixed

- Release artifacts no longer show `commit: none` or `date: unknown`

### Known Limitations

- Installed/local package status is not available because aurview does not call package managers at runtime.

### Validation Notes

- `gofmt -w .`
- `go test ./...`
- `go vet ./...`
- `go build ./...`
- Release-style `go build -ldflags ...` metadata check
- `makepkg --printsrcinfo`
- `makepkg --verifysource`

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
