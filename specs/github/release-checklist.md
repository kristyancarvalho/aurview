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
- [ ] Confirm no package install/download/clone/build/update/remove behavior exists.
- [ ] Check README instructions and badges.
- [ ] Check `/specs` for updated architecture and release notes.
- [ ] Create a signed or annotated tag if the project policy requires it.
- [ ] Publish GitHub release notes.
