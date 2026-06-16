# AUR RPC

aurview uses the official AUR RPC interface for metadata.

References:

- <https://wiki.archlinux.org/title/Aurweb_RPC_interface>
- <https://aur.archlinux.org/rpc>

Default search uses:

```text
type=search
by=name-desc
v=5
```

Detail lookup uses:

```text
type=info
v=5
```

The client uses context timeouts, keeps a small in-memory cache and decodes JSON into `internal/aur.Package`.
