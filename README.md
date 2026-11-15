# interval

A generic, comparator-based Go interval library with explicit open, closed, and unbounded endpoints.

The active module on this branch is [`github.com/dingyou/interval/v2`](v2/README.md). It requires Go 1.21 or later. This branch is not yet a published v2 release.

- [Usage and API examples](v2/README.md)
- [Migrating from v1](v2/README.md#migration-from-v1)

Run tests from the v2 module:

```sh
cd v2
go test -race ./...
go vet ./...
```

The previous v1 implementation remains available in Git history.
