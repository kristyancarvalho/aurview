You are Codex working on an existing Go TUI project called `aurview`.

Context:
The project already has an MVP. Before implementing anything, inspect the current repository structure, existing branches, current code style, tests, README, `.github` directory and any existing documentation.

Project goal:
Evolve `aurview` into a polished read-only package discovery TUI for Arch/AUR users.

Important:
This tool must remain read-only. It must never install, download, clone, build, update or remove packages. It is only for search, inspection, history, metadata visualization and copying package names to the clipboard.

Current new requirements:

1. Audit and improve `.github`
2. Create milestones and issues to log planned implementations
3. Add mouse support to the TUI
4. Add support for configurable package repositories, using AUR by default
5. Show the package repository/source in search results
6. Add theme support
7. Move development/spec documentation to `/specs`
8. Rewrite the README in a more didactic way for beginners
9. Add README widgets/badges suitable for open-source projects
10. Add a license: MIT or GPL-3.0, credited to `aurview contributors`

Versioning and workflow preferences:
Follow this exact workflow.

Base branch:

* Work from `dev`.
* Do not commit directly to `main`.
* Keep `main` stable.
* All implementation branches must merge back into `dev` only after validation.

Before coding:

* Inspect the repository.
* Inspect current git branches.
* Inspect existing tests.
* Inspect the `.github` directory.
* Inspect README and docs.
* Create a short implementation plan.
* Split the work into separate implementations with separate local staging branches.
* Each implementation must have its own branch, commits, tests and merge back into `dev`.

Do not squash everything into one large change.

Implementation 1 branch:
Create local branch:

```bash
staging/github-project-audit
```

Scope:

* Audit the `.github` directory.
* Add or improve issue templates.
* Add or improve pull request template.
* Add or improve GitHub Actions workflows.
* Add basic CI for Go:

  * checkout
  * setup-go
  * gofmt check
  * go vet ./...
  * go test ./...
  * go build ./cmd/aurview or the correct main package path
* Add Dependabot config if useful.
* Add labels documentation if useful.
* Create a `/specs/github/` directory containing markdown files for:

  * planned milestones
  * planned issues
  * implementation log
  * release checklist
  * contribution workflow

Milestones and issues:
If the GitHub CLI is available and authenticated, create the milestones and issues directly on GitHub.

If GitHub CLI is not available or not authenticated, do not fail the task. Instead, generate ready-to-copy markdown files under `/specs/github/issues/` and `/specs/github/milestones/`.

Suggested milestones:

* `v0.2.0 - TUI interaction polish`
* `v0.3.0 - Multi-source package search`
* `v0.4.0 - Themes and documentation`
* `v1.0.0 - Stable read-only AUR browser`

Suggested issues:

* Add mouse support to result list and detail panel
* Add configurable repository sources
* Show repository/source badge in package results
* Add theme system and built-in themes
* Move development documentation to `/specs`
* Rewrite README for beginners
* Add open-source project badges/widgets to README
* Add project license
* Improve CI and repository metadata
* Add release checklist
* Add contribution guide

Each issue should include:

* Context
* Goal
* Tasks
* Acceptance criteria
* Test/validation notes

Commit using Conventional Commits.
Run validation before merging into `dev`.

Implementation 2 branch:
Create local branch:

```bash
staging/tui-interaction-polish
```

Scope:

* Add mouse support.
* Keep existing keyboard and Vim motions working.
* Mouse support must include:

  * click package row to select it
  * double-click or Enter behavior should copy the selected package name, depending on what fits the current architecture
  * scroll result list with mouse wheel
  * scroll detail panel with mouse wheel if content overflows
  * click search input to focus it
  * click help/status areas only if useful
* Preserve terminal-first UX.
* Mouse support must be optional and should not make the app feel less keyboard-driven.

UX requirements:

* Keep the TUI visually unique.
* Do not make it look like a default Charm/Bubble Tea demo.
* Avoid stock list styling, generic borders and generic help bars.
* Keep compact density.
* Keep Vim motions:

  * j/k
  * h/l
  * gg/G
  * /
  * n/N
  * Ctrl+d/Ctrl+u
  * Ctrl+f/Ctrl+b
  * Enter
  * Esc
  * q
  * ?
* Add tests where practical for event handling, selection behavior and state changes.

Commit using Conventional Commits.
Run validation before merging into `dev`.

Implementation 3 branch:
Create local branch:

```bash
staging/configurable-sources
```

Scope:
Add support for searching packages from repositories declared in a configuration file, while using AUR by default.

Default behavior:

* If no config exists, search only AUR.
* AUR must remain the default source.
* The user should be able to run the app without configuring anything.

Configuration:
Use XDG paths on Linux.

Preferred config path:

```text
~/.config/aurview/config.toml
```

Also support a project/local config if reasonable:

```text
./aurview.toml
```

Suggested config shape:

```toml
default_sources = ["aur"]

[ui]
theme = "arch"

[[sources]]
name = "aur"
type = "aur-rpc"
enabled = true
url = "https://aur.archlinux.org/rpc"

[[sources]]
name = "custom"
type = "aur-rpc"
enabled = false
url = "https://example.com/rpc"
```

Repository/source behavior:

* Every package result must include a source/repository field.
* AUR results should show `AUR`.
* Results from configured custom sources should show the configured source name.
* Search results must clearly display the source when there is more than one enabled source.
* The detail panel must show the source/repository.
* Ranking should work across sources.
* Avoid duplicate confusion:

  * if the same package name appears in multiple sources, show both entries with clear source labels
  * optionally group or visually distinguish them
* Keep the source abstraction clean:

  * `Source` interface
  * `Search(ctx, query) ([]Package, error)`
  * `Name() string`
  * `Type() string`
* Keep AUR implementation isolated.
* Do not hardcode all package logic into the TUI.

Important:
Do not implement installation, download or clone behavior.
Do not call pacman, yay, paru, makepkg or git clone.
This is metadata-only.

Add tests for:

* config loading
* default config behavior
* source parsing
* disabled sources
* multi-source result normalization
* package source display formatting
* duplicate package names across sources

Commit using Conventional Commits.
Run validation before merging into `dev`.

Implementation 4 branch:
Create local branch:

```bash
staging/themes-and-docs
```

Scope:

* Add theme support.
* Move development documentation to `/specs`.
* Rewrite README for beginners.
* Add open-source widgets/badges.
* Add license.

Theme support:
Add built-in themes, for example:

* `arch`
* `mono`
* `dark`
* `light`
* `high-contrast`

Theme requirements:

* Theme should control:

  * accent color
  * muted text
  * warning color
  * error color
  * success color
  * border/separator style
  * selected row style
  * source/repository badge style
  * status bar style
* Keep a monochrome fallback.
* Respect terminal limitations.
* Avoid unreadable colors.
* Allow selecting theme from config:

```toml
[ui]
theme = "arch"
```

* Keep theme definitions isolated from business logic.

Documentation:
Move development-oriented documentation to `/specs`.

Suggested `/specs` structure:

```text
/specs
  /architecture.md
  /aur-rpc.md
  /configuration.md
  /themes.md
  /keyboard-and-mouse.md
  /ranking.md
  /github
    /implementation-log.md
    /release-checklist.md
    /milestones.md
    /issues
```

README rewrite:
Rewrite README.md to be didactic and beginner-friendly.

README should include:

* What is aurview?
* What problem does it solve?
* What it does not do
* Clear warning that it is read-only
* Why read-only matters
* Installation/build instructions
* Basic usage
* Search examples
* Keyboard shortcuts
* Mouse support
* Config file example
* Theme configuration example
* Package source/repository explanation
* Screenshot placeholder
* Project status
* Roadmap
* Contributing section
* License section
* Credits to `aurview contributors`

README widgets/badges:
Add useful open-source badges/widgets, such as:

* Go version
* CI status
* License
* Release
* Issues
* Pull requests
* Last commit
* Stars
* Forks
* Go Report Card if applicable
* Codecov only if coverage is actually configured

Do not add fake badges for services that are not configured.
Do not add broken badge URLs.
Prefer shields.io badges when appropriate.

License:
Add either MIT or GPL-3.0.

Preferred default:

* Use MIT unless the repository already indicates GPL preference.

License attribution:
Use:

```text
aurview contributors
```

Examples:

```text
MIT License

Copyright (c) 2026 aurview contributors
```

or, for GPL:

```text
Copyright (C) 2026 aurview contributors
```

Commit using Conventional Commits.
Run validation before merging into `dev`.

Final validation:
After all branches are merged into `dev`, run:

```bash
gofmt -w .
go test ./...
go vet ./...
go build ./...
```

Also manually validate:

* App starts without config
* AUR search still works
* Search result shows source/repository
* Package detail shows source/repository
* Enter copies package name
* Mouse click selects package
* Mouse scroll works
* Vim motions still work
* Theme can be changed via config
* Missing clipboard provider does not crash the app
* Missing config does not crash the app
* Invalid config shows a helpful error
* README links and badges are valid
* `/specs` contains development docs
* `.github` has useful templates/workflows

Conventional Commit examples:

* chore(github): audit repository workflows and templates
* docs(specs): add project milestone and issue specs
* feat(tui): add mouse selection and scrolling
* feat(config): add configurable package sources
* feat(sources): show package repository in search results
* feat(theme): add configurable theme system
* docs(readme): rewrite beginner-friendly project guide
* chore(license): add MIT license for aurview contributors
* test(config): cover source loading behavior
* test(ranking): cover multi-source package scoring

Rules:

* Do not break existing MVP behavior.
* Do not remove keyboard-first navigation.
* Do not make the UI look generic.
* Do not introduce installation behavior.
* Do not call package managers.
* Do not clone AUR repositories.
* Do not execute package scripts.
* Do not hide errors.
* Keep architecture modular.
* Keep commits small and meaningful.
* Keep all implementation logs under `/specs/github/implementation-log.md`.

Expected result:
A more mature `aurview` MVP with:

* audited `.github`
* planned milestones and issues
* useful CI
* mouse support
* configurable package sources
* AUR as default source
* source/repository labels in search and detail views
* theme support
* beginner-friendly README
* development specs under `/specs`
* open-source badges/widgets
* MIT or GPL-3.0 license credited to `aurview contributors`
* all changes implemented through staging branches and merged into `dev` only after tests pass.

