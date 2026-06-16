# Release Checklist

Use this checklist before creating a release tag.

## Repository State

- [ ] Confirm the working tree is clean.
- [ ] Confirm the current branch is `dev` while preparing release changes.
- [ ] Confirm `dev` has been validated and merged through the expected workflow.
- [ ] Confirm `main` is stable and contains only release-ready changes.
- [ ] Confirm `/specs/` is still untracked and ignored.

## Versioning And Docs

- [ ] Confirm the intended release version is consistent across `CHANGELOG.md`, release notes and packaging metadata.
- [ ] Check `README.md` instructions and badges.
- [ ] Check `LICENSE` and confirm it credits `aurview contributors`.
- [ ] Check `CHANGELOG.md` and add the release entry if needed.
- [ ] Check `CONTRIBUTING.md` and tracked `docs/` content for stale release guidance.

## Local Validation

- [ ] Run `gofmt -w .`.
- [ ] Run `go test ./...`.
- [ ] Run `go vet ./...`.
- [ ] Run `go build ./...`.
- [ ] Start the app without a config file.
- [ ] Confirm AUR search works.
- [ ] Confirm package details render.
- [ ] Confirm `Enter` copies the selected package name or reports a missing clipboard provider.
- [ ] Confirm mouse row selection and scrolling still work in a terminal that supports mouse events.
- [ ] Confirm configured themes remain readable, including `mono`.
- [ ] Confirm no package install, download, clone, build, update, remove or script execution behavior exists.

## GitHub Release

- [ ] Confirm the release tag uses `vX.Y.Z`.
- [ ] Create an annotated tag if that is the intended release action.
- [ ] Push the release tag to GitHub.
- [ ] Confirm the tag points to a commit reachable from `dev`.
- [ ] Confirm the GitHub release workflow uploads release artifacts.
- [ ] Confirm the GitHub release workflow uploads the source tarball.
- [ ] Confirm the GitHub release workflow uploads checksums.
- [ ] Confirm GitHub release notes and attached artifacts are correct.

## AUR Packaging

- [ ] Confirm `packaging/aur/PKGBUILD` exists and matches the release version.
- [ ] Confirm `packaging/aur/.SRCINFO` matches `PKGBUILD`.
- [ ] Run `makepkg --printsrcinfo > .SRCINFO` from `packaging/aur`.
- [ ] Run `makepkg --verifysource` from `packaging/aur`.
- [ ] Optionally run `makepkg -f` in a clean packaging environment.
- [ ] Confirm the source URL and checksum are valid for the release tag.

## AUR Publication

- [ ] Confirm the AUR remote or temporary AUR clone is configured.
- [ ] Update the AUR repository checkout with `PKGBUILD` and `.SRCINFO` only.
- [ ] Commit the AUR repository update with the release version.
- [ ] Push to the AUR `master` branch.
- [ ] Verify the AUR package page after publication.
