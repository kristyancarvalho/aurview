# Add theme system and built-in themes

## Context

aurview has a custom visual identity, but colors and separators are currently hardcoded.

## Goal

Add isolated, configurable themes while preserving monochrome fallback behavior.

## Tasks

- [ ] Add built-in `arch`, `mono`, `dark`, `light` and `high-contrast` themes.
- [ ] Let config select `[ui].theme`.
- [ ] Theme accent, muted, warning, error, success, separators, selected row, source badges and status bar.
- [ ] Keep terminal capability checks.

## Acceptance Criteria

- [ ] Theme definitions are isolated from business logic.
- [ ] Unknown theme names return helpful errors or fall back clearly.
- [ ] Monochrome terminals remain readable.

## Test / Validation Notes

- Add tests for theme lookup and fallback behavior.
