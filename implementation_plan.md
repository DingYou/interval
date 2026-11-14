# Interval v2 implementation plan

1. Preserve the root v1 module unchanged; add a separate `/v2` module.
2. Implement one generic, comparator-backed interval with explicit finite and unbounded bounds. Zero value is empty. Reject reversed endpoints and nil comparators. Canonicalize equal excluded bounds to empty.
3. Separate parsing/formatting from ordering, support escaped endpoint text, and provide RFC3339Nano time helpers.
4. Add boundary, error, float, custom-type, parser, property/fuzz, example, and regression tests. Document migration and ordering/domain limitations.
5. Format, test both modules, run race/vet/fuzz checks using an official Go toolchain when available. Initial scope excluded commits and pushes; the user subsequently authorized a `2.0.0` branch and local commit. Remote push remains out of scope.

## Design decisions

- Minimum Go 1.21 for `cmp.Ordered` and `cmp.Compare`; compatible with Go 1.27 without gratuitous generic methods.
- No public Domain abstraction or binary set operations yet; each interval owns a comparison function.
- `IsEmpty` intentionally omitted: ordering alone cannot determine whether a discrete domain has elements strictly between distinct endpoints. `IsCanonicalEmpty` reports only the normalized empty representation.
- Ordered floats use `cmp.Compare` total ordering: NaN sorts before every non-NaN and equals NaN; signed zeros compare equal.
- v1 source and its nullable-time semantics are not silently changed.

## Completed verification

- Implemented separate v2 module, explicit bounds, unified membership, strict quoted-endpoint parsing, time helpers, API/migration documentation, and two-module CI matrix.
- v1 tracked source files were left unchanged.
- Official Go 1.27.4 linux/amd64 archive downloaded from go.dev and verified against its published SHA256 (`7226736c2339719cda28fbf3a54d71aa63bcefe376e98cef00034268dab24840`). Toolchain/cache installed outside the repository; no system installation.
- Go 1.27.4: v1 tests and race tests passed; v2 race tests and vet passed; v2 statement coverage 100.0%.
- Five-second fuzz campaigns passed: parser 277,629 executions; membership 332,095; arbitrary string round trip 56,702. These are finite checks, not correctness proofs.
- Contains microbenchmark on this cloud machine: 7.777 ns/op, 0 B/op, 0 allocations/op. Not a cross-machine performance guarantee.
- Go 1.21.13 minimum-version tests passed for both modules (toolchain obtained via Go's verified toolchain download mechanism).
- Independent read-only implementation review found no blocking defects; its string-fuzz and time-documentation suggestions were incorporated.
- `git diff --check` passed. No commit, push, tag, or published release was created. CI configuration is written but remote CI has not run.

## Authorized local commit follow-up

- The original task's read-only session prevented Git metadata writes. A later writable execution session verified `.git` writes and resumed the user-authorized `2.0.0` branch/local commit workflow.
- Re-ran v1 and v2 race tests and vet with verified Go 1.27.4; all passed, v2 statement coverage remained 100.0%. `gofmt` and whitespace checks passed.
- Remote push, tag, and release publication remain out of scope.
