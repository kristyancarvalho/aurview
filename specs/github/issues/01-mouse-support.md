# Add mouse support to result list and detail panel

## Context

aurview is keyboard-first, but many terminal users expect mouse selection and wheel scrolling to work when their terminal supports it.

## Goal

Add optional mouse support without weakening the existing Vim-oriented workflow.

## Tasks

- [ ] Enable Bubble Tea mouse events.
- [ ] Click package rows to select them.
- [ ] Double-click a selected package row to copy the package name, if practical.
- [ ] Scroll the package list with the mouse wheel.
- [ ] Scroll the detail panel with the mouse wheel.
- [ ] Click the search area to focus search.
- [ ] Preserve existing keyboard behavior.

## Acceptance Criteria

- [ ] Mouse row clicks update the selected package.
- [ ] Mouse wheel scrolling works in result and detail areas.
- [ ] Keyboard shortcuts still work.
- [ ] Mouse support is optional and does not block non-mouse terminals.

## Test / Validation Notes

- Add unit tests for coordinate-to-action behavior where practical.
- Manually test click, double-click and wheel events in a terminal.
