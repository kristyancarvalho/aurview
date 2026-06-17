# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and the project uses Semantic Versioning
for tagged releases.

## [Unreleased]

## [0.5.0] - 2026-06-17

### Added

- Read-only local pacman sync database sources for detected repositories such as `core`, `extra`, `multilib` and `chaotic-aur`
- Automatic local repository detection through `pacman-conf --repo-list`
- Matugen color theme support through `[ui].theme = "matugen"` and `[theme]` hex color fields
- Configuration documentation covering paths, precedence, source options, local repository detection and Matugen templates

### Changed

- Default sources now include AUR plus detected enabled local pacman repositories
- Search results and filters continue to label each package by source repository across mixed AUR and local results
- Missing local sync databases are skipped without failing other sources

## [0.4.2] - 2026-06-17

### Added

- Interactive package filters for source, maintainer state, out-of-date state, votes, popularity, recency and match mode
- Compact filter bar with keyboard navigation and reset support

### Changed

- Styled top headers, table headers and filter chips with compact htop-like filled backgrounds
- Moved build metadata to `internal/version` so release builds can embed version, commit and date consistently
- Updated AUR packaging metadata for the `v0.4.2` release

### Fixed

- Released binaries now embed real commit and build date metadata instead of `commit: none` and `date: unknown`

## [0.4.1] - 2026-06-17

### Added

- Release notes under `docs/release-notes.md`
- Build-time version output through `aurview --version`

### Changed

- Simplified `README.md` around practical installation, usage, configuration and release information
- Validated release tags against `main`, the stable release branch
- Updated release artifact upload/download actions
- Checked out repository history before GitHub Release publication
- Kept staging branches local-only and `docs/` focused on release notes
- Updated AUR package metadata for the `v0.4.1` release

### Fixed

- Made source archive generation tolerate repositories without a `docs/` directory

## [0.4.0] - 2026-06-16

### Added

- Read-only AUR RPC package search and metadata inspection
- Local ranking, search history and clipboard copy support
- Mouse interaction, configurable metadata sources and built-in themes
- Release checklist and maintainer-facing CI/CD and AUR publishing guides
