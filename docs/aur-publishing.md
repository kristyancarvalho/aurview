# AUR Publishing

aurview is published to the AUR as a source-built package. The application
itself remains read-only; these steps are maintainer actions only.

## Prerequisites

- A tagged GitHub release in the form `vX.Y.Z`
- A clean local working tree
- `makepkg` available locally
- AUR SSH access configured for `aur.archlinux.org`

## Repository Metadata

Before updating the AUR package:

- confirm `CHANGELOG.md`, `README.md` and `LICENSE` are current
- confirm `packaging/aur/PKGBUILD` uses the release tag source archive
- regenerate `packaging/aur/.SRCINFO`
- verify the source archive and checksum

## Validation Commands

From the repository root:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./...
```

From `packaging/aur`:

```sh
makepkg --printsrcinfo > .SRCINFO
makepkg --verifysource
makepkg -f
```

Use `makepkg -f` only in a clean packaging environment. It is not required for
every metadata-only update if `--verifysource` already proves the release
archive is available and correct.

## Publishing Steps

Clone or update the AUR repository in a separate working directory:

```sh
mkdir -p /tmp/aurview-aur
cd /tmp/aurview-aur
git clone ssh://aur@aur.archlinux.org/aurview.git .
```

If the package does not exist yet, initialize the empty clone with the release
packaging files:

```sh
cp /path/to/aurview/packaging/aur/PKGBUILD .
cp /path/to/aurview/packaging/aur/.SRCINFO .
git add PKGBUILD .SRCINFO
git commit -m "add aurview"
git push origin master
```

For updates:

```sh
cp /path/to/aurview/packaging/aur/PKGBUILD .
cp /path/to/aurview/packaging/aur/.SRCINFO .
makepkg --printsrcinfo > .SRCINFO
makepkg --verifysource
git add PKGBUILD .SRCINFO
git commit -m "update aurview to vX.Y.Z"
git push origin master
```

After pushing, verify the package page on the AUR and confirm that the source
archive and metadata render correctly.
