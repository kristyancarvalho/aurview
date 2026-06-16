# Release Checklist

Use this checklist before creating a release tag.

- [ ] Confirm `dev` has been validated and merged through the expected workflow.
- [ ] Confirm `main` is stable and contains only release-ready changes.
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
- [ ] Check README instructions and badges.
- [ ] Check `CONTRIBUTING.md` and tracked `docs/` content for stale release guidance.
- [ ] Create a signed or annotated tag if the project policy requires it.
- [ ] Publish GitHub release notes.
