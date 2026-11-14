# interval

A generic Go interval library.

- **v1**: the original module at `github.com/dingyou/interval` (Go 1.18). Existing API and behavior are preserved.
- **v2**: a redesigned, comparator-based API in [`v2`](v2/README.md), module `github.com/dingyou/interval/v2` (Go 1.21+). This is an unreleased implementation, not a published version.

Run tests in both modules independently:

```sh
go test ./...
(cd v2 && go test ./...)
```

The v2 module intentionally does not silently alter v1 endpoint, parser, or nullable-time semantics. See the [migration guide](v2/README.md#migration-from-v1).
