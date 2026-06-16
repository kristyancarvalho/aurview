# CI/CD

aurview uses two GitHub Actions workflows:

- `ci`: runs on pushes to `dev`, `main` and `staging/**`, plus pull requests
- `release`: runs only on tags that start with `v`

## CI Workflow

The regular CI workflow validates repository health with:

- `gofmt -w .` followed by `git diff --exit-code`
- `go vet ./...`
- `go test ./...`
- `go build ./...`

This keeps staging branches and `dev` ready to merge without publishing anything.

## Release Workflow

The release workflow is intentionally narrow:

- it accepts only semver tags in the form `vX.Y.Z`
- it verifies that the tagged commit is reachable from `origin/main`
- it reruns formatting, vet, tests and full build checks
- it builds Linux release archives for `amd64` and `arm64`
- it generates a `SHA256SUMS` file for the uploaded artifacts
- it publishes a GitHub Release using the repository `GH_TOKEN`

The workflow does not publish on ordinary pushes and does not attempt to update
the AUR automatically.

## AUR Publishing

AUR publishing remains a maintainer action outside GitHub Actions. That keeps
SSH credentials and AUR history management out of the default release path.

The expected maintainer flow is:

1. Validate the repository release checklist in [`docs/release-checklist.md`](release-checklist.md).
2. Update `packaging/aur/PKGBUILD` and `packaging/aur/.SRCINFO` for the release tag.
3. Validate sources with `makepkg --verifysource`.
4. Push the packaging update to the AUR repository manually.
