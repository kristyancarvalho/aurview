# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and the project uses Semantic Versioning
for tagged releases.

## [Unreleased]

## [0.4.1] - 2026-06-17

### Added

- Release notes under `docs/release-notes.md`
- Build-time version output through `aurview --version`

### Changed

- Simplified `README.md` around practical installation, usage, configuration and release information
- Validated release tags against `main`, the stable release branch
- Updated release artifact upload/download actions
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
