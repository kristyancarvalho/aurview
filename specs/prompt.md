You are Codex working on a new Go project.

Project goal:
Build a lightweight, fast and visually unique terminal UI for browsing, searching and inspecting AUR packages.

Working name:
aurview

Core principle:
This project is read-only. It must never install, download, clone, build, update or remove packages. It only searches, ranks, previews and copies package names to the user's clipboard.

Reference:
Use the official AUR RPC interface as the data source:
https://wiki.archlinux.org/title/Aurweb_RPC_interface
https://aur.archlinux.org/rpc

Tech stack:
- Go
- A lightweight TUI approach
- You may use Bubble Tea/Charm libraries only if useful, but the final result must not look like a default Charm/Bubble Tea application.
- Avoid obvious stock widgets, default borders, default list styles, default help bars and generic Charm examples.
- Build a custom visual identity with a compact, information-dense layout.
- Keep dependencies minimal and justified.
- Prefer clean internal abstractions so the rendering engine can be replaced later.

Important UX goal:
The TUI should feel like a custom tool made for Arch/AUR power users, not like another generic Bubble Tea demo.

Main features:
1. Search AUR packages
   - Search by package name and keywords.
   - Use the AUR RPC API.
   - Debounce user input.
   - Handle loading, empty, error and rate-limit states gracefully.
   - Never block the UI while searching.
   - Cache recent queries in memory during the session.

2. Package list
   - Show results in a compact list.
   - Each row should include:
     - package name
     - version
     - description preview
     - votes
     - popularity
     - maintainer
     - last updated date
     - out-of-date status when available
     - relevance score calculated locally
   - Sort by a balanced relevance score, not only the order returned by the API.

3. Relevance scoring
   Implement a local relevance score based on:
   - exact name match
   - prefix match
   - substring match
   - description match
   - votes
   - popularity
   - recency / last update
   - package not flagged out-of-date
   - maintained vs orphaned package

   Keep this scoring logic isolated and covered by tests.

4. Package detail panel
   When a package is selected, show a rich detail view with:
   - name
   - package base
   - version
   - description
   - URL
   - maintainer
   - votes
   - popularity
   - first submitted
   - last modified
   - out-of-date status
   - license
   - dependencies
   - make dependencies
   - check dependencies
   - optional dependencies
   - conflicts
   - provides
   - keywords when available
   - AUR URL

5. History
   - Keep an in-session history of searched terms.
   - Allow navigating previous searches.
   - Show a small history area or command palette-style overlay.
   - Store persistent history locally if it can be done cleanly.
   - Use XDG paths on Linux.

6. Clipboard behavior
   - When a package is selected and the user presses Enter, copy the package name to the user's clipboard.
   - Show a subtle confirmation message.
   - Do not download or install anything.
   - Make clipboard support work on Linux first.
   - Gracefully handle missing clipboard providers.
   - Keep clipboard logic behind an interface for testability.

7. Vim motions
   Must support:
   - j / k: move down/up
   - h / l: move focus between list and detail/search areas when applicable
   - g g: go to top
   - G: go to bottom
   - /: focus search
   - n / N: move through search history or matches where appropriate
   - Ctrl+d / Ctrl+u: half-page down/up
   - Ctrl+f / Ctrl+b: page down/up
   - Enter: copy selected package name
   - Esc: blur input / close overlay
   - q: quit
   - ?: toggle help overlay

8. Visual design
   Avoid a generic Bubble Tea look.
   Design a unique layout inspired by:
   - compact Arch Linux tooling
   - terminal dashboards
   - fuzzy finders
   - package browsers
   - minimal but information-rich TUIs

   Visual requirements:
   - custom borders or separators
   - subtle status bar
   - compact header
   - clear focus indicators
   - no excessive padding
   - no huge empty spaces
   - adaptive layout for small terminals
   - readable monochrome fallback
   - tasteful use of ANSI colors
   - avoid looking like default Charm examples

9. Architecture
   Use a clean, modular Go structure.

   Suggested structure:
   - cmd/aurview/main.go
   - internal/app
   - internal/aur
   - internal/clipboard
   - internal/config
   - internal/history
   - internal/ranking
   - internal/tui
   - internal/tui/theme
   - internal/tui/keymap
   - internal/tui/components
   - internal/platform

   Requirements:
   - AUR API client must be isolated.
   - Ranking must be isolated.
   - Clipboard must be isolated.
   - TUI state/update/rendering must be organized and testable.
   - No business logic directly inside main.go.
   - main.go should only wire dependencies and start the app.

10. Testing
   Add tests for:
   - AUR API response parsing
   - ranking/relevance score
   - history behavior
   - clipboard interface fallback behavior
   - keymap behavior where practical
   - formatting helpers
   - date conversion helpers

   Use table-driven tests where appropriate.

11. Error handling
   Handle:
   - no internet
   - AUR API unavailable
   - invalid JSON
   - empty search
   - no results
   - terminal too small
   - clipboard unavailable
   - request timeout

12. Performance
   - Keep startup fast.
   - Use context timeouts for HTTP calls.
   - Avoid unnecessary allocations in rendering hot paths.
   - Avoid excessive goroutines.
   - Keep the TUI responsive during searches.

13. Documentation
   Add:
   - README.md
   - usage examples
   - keybindings table
   - explanation that the tool is read-only
   - installation/build instructions
   - development instructions
   - limitations
   - screenshots placeholder section
   - AUR RPC reference

14. Local development
   - Add a Makefile with useful targets:
     - make build
     - make run
     - make test
     - make lint
     - make fmt
     - make clean
   - Add a minimal Dockerfile or devcontainer only if useful, but the app itself is a terminal app and should work natively.
   - Do not over-engineer the initial setup.

Versioning and workflow preferences:
Follow this exact workflow.

Base branch:
- Work from dev.
- Do not commit directly to main.
- Keep main stable.
- All implementation branches must merge back into dev only after tests pass.

Before coding:
- Inspect the repository.
- If no repository structure exists, initialize a clean Go module.
- Create a short implementation plan.
- Split the work into two separate implementations with separate local staging branches.

Implementation 1 branch:
- Create local branch: staging/aurview-core
- Scope:
  - Go module setup
  - AUR RPC client
  - package models
  - ranking/relevance logic
  - history service
  - clipboard abstraction
  - tests for core logic
  - README initial draft
- Commit using Conventional Commits.
- Run tests before merging.
- Merge into dev only after validation.

Implementation 2 branch:
- Create local branch: staging/aurview-tui
- Scope:
  - TUI layout
  - custom visual identity
  - search input
  - result list
  - detail panel
  - vim motions
  - clipboard action on selection
  - help overlay
  - status/error states
  - Makefile targets
  - final README updates
- Commit using Conventional Commits.
- Run tests before merging.
- Merge into dev only after validation.

Commit style:
Use Conventional Commits:
- feat:
- fix:
- refactor:
- test:
- docs:
- chore:
- build:

Examples:
- feat(aur): add RPC client and package models
- feat(ranking): implement local relevance scoring
- feat(tui): add custom package browser layout
- feat(keymap): add vim-style navigation
- feat(clipboard): copy selected package name
- test(ranking): cover relevance score cases
- docs(readme): document read-only AUR browser usage

Quality gate:
Before each merge into dev, run:
- gofmt
- go test ./...
- go vet ./...
- go run ./cmd/aurview when possible

Do not skip tests.
Do not leave broken code in dev.
Do not add install/build/package-management behavior.
Do not shell out to yay, paru, pacman or makepkg.
Do not clone AUR repositories.
Do not execute package scripts.
Do not implement package installation.

Expected result:
A polished initial version of a read-only AUR package browser TUI in Go with:
- unique custom visual design
- fast AUR search
- detailed package preview
- relevance scoring
- search history
- vim motions
- clipboard copy on selection
- clean architecture
- tests
- documentation
- versioned implementation through the two staging branches and final merge into dev.
