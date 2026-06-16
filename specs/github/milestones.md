# Planned Milestones

## v0.2.0 - TUI interaction polish

Focus: keyboard-first interaction polish with optional mouse support.

- Add package row mouse selection.
- Add result list and detail panel mouse wheel scrolling.
- Preserve Vim motions and terminal-first UX.
- Add practical event handling tests.

## v0.3.0 - Multi-source package search

Focus: configurable package metadata sources with AUR as the default.

- Add XDG config loading.
- Add source abstraction.
- Add AUR RPC source implementation.
- Show source/repository in result rows and detail panel.
- Rank results across enabled sources.

## v0.4.0 - Themes and documentation

Focus: theme support and beginner-friendly documentation.

- Add built-in themes.
- Select theme from config.
- Move development docs into `/specs`.
- Rewrite README for beginners.
- Add open-source badges and project license.

## v1.0.0 - Stable read-only AUR browser

Focus: stable read-only package discovery workflow.

- Harden errors and edge cases.
- Finalize source/config behavior.
- Finalize docs and contribution workflow.
- Confirm read-only behavior throughout the codebase.
