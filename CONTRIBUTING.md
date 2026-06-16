# Contributing

aurview uses `dev` as the integration branch and keeps `main` stable.

## Branching

- Start feature work from `dev`.
- Use focused staging branches, for example `staging/tui-layout-fixes`.
- Do not commit directly to `main`.
- Merge back into `dev` only after validation passes.

## Commit Style

Use Conventional Commits with a scope when it helps clarity:

```text
feat(tui): add keyboard navigation
fix(config): reject duplicate source names
docs(readme): clarify read-only behavior
test(ranking): cover score ordering
chore(repo): update repository metadata
```

## Validation

Run the standard checks before opening a pull request:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/aurview
```

The `Makefile` provides equivalent shortcuts:

```sh
make fmt
make test
make lint
make build
```

## Read-Only Safety

aurview is metadata-only. Contributions must not add package installation,
downloads, repository cloning, package builds, updates, removals or script
execution.

Do not call `pacman`, `yay`, `paru`, `makepkg`, AUR helpers or `git clone` from
the application.
