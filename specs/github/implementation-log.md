# GitHub Implementation Log

## 2026-06-16

- Audited existing `.github` assets.
- Found issue templates and workflows copied from another project.
- Updated issue template areas to match aurview.
- Added read-only safety confirmations to issue templates and the pull request template.
- Replaced the CI workflow with a Go-focused workflow for formatting, vetting, tests and build.
- Simplified release workflow so it builds `./cmd/aurview` instead of an unrelated project path.
- Added Dependabot configuration for Go modules and GitHub Actions.
- Added project planning specs for milestones, issues, labels, release checklist and contribution workflow.
- Checked GitHub CLI availability with `gh auth status`; authentication was present for `kristyancarvalho/aurview`.
- Attempted to create labels, milestones and issues with `gh`, but subsequent `gh issue`, `gh label` and `gh api` commands hung in this environment and had to be stopped.
- Preserved ready-to-copy milestone and issue bodies under `/specs/github/milestones/` and `/specs/github/issues/` as the required fallback path.

## 2026-06-16 - TUI interaction polish

- Added Bubble Tea mouse cell motion support in the app runner.
- Added custom coordinate hit-testing for search, result list rows and detail panel.
- Added mouse click support for focusing search and selecting package rows.
- Added double-click copy behavior for selected package rows.
- Added mouse wheel support for result list selection and detail scrolling.
- Preserved existing keyboard and Vim motion handling.
- Added TUI mouse interaction tests.

## 2026-06-16 - Configurable sources

- Added TOML config loading with local `./aurview.toml` and XDG config path support.
- Kept AUR as the default source when no config file exists.
- Added an isolated `Source` interface and multi-source client.
- Added AUR RPC source adapter that stamps packages with source metadata.
- Updated TUI search/detail flow to use package source plus package name for detail lookups.
- Added source labels in package rows and source metadata in the detail panel.
- Added tests for config defaults, disabled sources, source parsing and duplicate package names across sources.
