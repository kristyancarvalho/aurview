# Show repository/source badge in package results

## Context

When multiple sources are enabled, identical package names may appear from different repositories. The TUI must make source identity explicit.

## Goal

Show source/repository labels in result rows and package details.

## Tasks

- [ ] Add source metadata to normalized package results.
- [ ] Show `AUR` for the default AUR source.
- [ ] Show configured source names for custom sources.
- [ ] Include source/repository in the detail panel.
- [ ] Keep duplicate package names visible as separate rows.

## Acceptance Criteria

- [ ] Search results show a source label when appropriate.
- [ ] Package details show source/repository.
- [ ] Duplicate names from multiple sources remain distinguishable.

## Test / Validation Notes

- Add tests for source label formatting and duplicate package names.
