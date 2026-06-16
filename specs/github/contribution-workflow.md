# Contribution Workflow

aurview uses `dev` as the integration branch and keeps `main` stable.

## Branching

- Start feature work from `dev`.
- Use a focused branch name such as `staging/tui-interaction-polish`.
- Do not commit directly to `main`.
- Merge back into `dev` only after validation passes.

## Commit Style

Use Conventional Commits:

- `feat:`
- `fix:`
- `refactor:`
- `test:`
- `docs:`
- `chore:`
- `build:`

## Validation

Run:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/aurview
```

## Read-only Safety

aurview is metadata-only. Contributions must not add package installation, downloads, cloning, builds, updates, removals or script execution.
