# interval v2

One generic interval type, explicit endpoint kinds, and caller-selected ordering. No dependencies. Requires Go 1.21 or newer; works with Go 1.27. This checkout is not yet a published v2 release.

```go
import interval "github.com/dingyou/interval/v2"

window, err := interval.Ordered(interval.Inclusive(1), interval.Exclusive(5))
if err != nil {
    return err
}
fmt.Println(window.Contains(1)) // true
fmt.Println(window.Contains(5)) // false
fmt.Println(window)             // [1,5)
```

## Why this shape?

`Interval[T any]` owns its comparator; `Bound[T]` stores a value and whether it is inclusive, exclusive, or unbounded. There are no separate numeric/time/custom-type membership implementations and no requirement to modify a type to add `CompareTo` methods. There is no reflection or retained endpoint pointer in the time helpers.

Use `Ordered` for `cmp.Ordered` types, `Time` for `time.Time`, or `New(lower, upper, compare)` for any other ordering. Comparators return a negative/zero/positive result, not necessarily -1/0/1, and must define a deterministic total order. Reverse order and equivalence classes are supported. The caller must not mutate values or comparator state in ways that change ordering. Copies are shallow for types containing slices, maps, pointers, or closures. Concurrent reads require a concurrency-safe comparator.

Go 1.27's generic methods are not required by these operations: their endpoint type is already the receiver's type parameter. The design uses Go 1.21 `cmp` rather than forcing a newer compiler simply for syntax. Generic methods could become useful for future transformations, but no unsafe arbitrary `Map` operation is exposed: it could reverse order or destroy interval structure.

## Empty, unbounded, and discrete domains

- `var i Interval[T]` and `Empty[T]()` are the canonical empty interval.
- `New(Unbounded[T](), Unbounded[T](), compare)` contains every value in the comparator's domain.
- `[x,x]` contains the equivalence class of `x`; `(x,x)`, `[x,x)`, and `(x,x]` normalize to canonical empty.
- Reversed finite endpoints return `ErrReversed`, not an implicit swap.
- Nil comparison functions return `ErrNilCompare`, even for fully unbounded intervals.
- Membership must satisfy **both** bounds.
- `IsCanonicalEmpty` only recognizes the explicit/normalized empty representation. There is deliberately no `IsEmpty` claiming domain-aware emptiness: `(1,2)` has no integer members but does have real-number members. Such discrete reasoning needs successor/predecessor knowledge, outside this minimal API.
- `Bounds()` returns `ok == false` for canonical empty. Always check `ok`: its zero returned bounds must not be mistaken for a full-domain interval.
- `Bound.Value()` similarly returns a boolean distinguishing finite endpoints from unbounded endpoints.

For floats, `Ordered` uses **`cmp.Compare` semantics**: NaN is before every non-NaN, all NaNs compare equal, and signed zeros compare equal. For example `[NaN,NaN]` contains NaN but not zero. This is a documented total order, not IEEE `<`/`==` membership. Applications that disallow NaNs should validate input before construction/query.

## Text and time

`Parse(text, decode, compare)` separates endpoint decoding from interval algebra. It accepts strict bracket directions and exactly two endpoints, with a comma separator. `empty` is the empty interval; `(,)` is the full domain; `(,2]` and `[2,)` are one-sided. An omitted endpoint must have an open bracket. Outer whitespace is ignored. Reversed endpoints fail validation after decoding.

Finite empty strings, delimiters, outer whitespace, quotes, and backslashes are quoted as Go string literals. For example `["",""]` is the singleton empty string, while `(,)` is unbounded. `Format` emits canonical quoting; `Parse` unquotes before calling your decoder. Endpoint errors are wrapped with lower/upper context and preserved for `errors.Is` and `errors.As`. Nil decoders return `ErrNilDecode`; malformed interval structure returns `ErrSyntax`.

`String()` satisfies `fmt.Stringer` and formats endpoints using `fmt.Sprint`. It is a display representation, not a universal serialization guarantee. For round trips use `i.Format(encode)` with a matching `decode`; nil `encode` falls back to `fmt.Sprint`.

`Time` uses `time.Time.Compare`; `ParseTime` and `FormatTime` use RFC3339Nano to preserve fractional seconds. Text encoding does not preserve a monotonic clock reading or named Location identity. Normalize input with `t.Round(0)` when consistent wall-clock ordering across serialization is important; mixing values with and without monotonic readings is not a general total-order guarantee. For durable serialization use ordinary RFC3339-representable civil timestamps; Go's time formatter is not a validator for arbitrary years/timezone offsets. All time constructors share the generic membership implementation.

For nullable time text, `ParseTime` also recognizes these case-insensitive aliases after trimming whitespace:

- `NULL` and `nil` are equivalent and mean an unbounded endpoint, regardless of the bracket on their own side.
- `-inf` is accepted only on the lower end, `inf` or `+inf` only on the upper end. Like an omitted endpoint, these require an open bracket.
- Any two valid unbounded endpoints describe the full time domain, even with mixed aliases: for example, `(,)`, `(-inf,+inf)`, `[null,nil]`, and `(nil,+inf)`. Use `empty` for an explicitly empty interval. If migrating v1 text that used two `NULL` endpoints to mean empty, convert that text to `empty` before parsing.

Quoted endpoint tokens are decoded as timestamps, not aliases. `FormatTime` emits canonical text (`empty` for empty, omitted text for unbounded ends). The generic `Parse` does not recognize these aliases.

When `start` and `end` are `*time.Time`, construct a half-open window like this:

```go
window, err := interval.Time(
    interval.TimeInclusive(start),
    interval.TimeExclusive(end),
)
if err != nil {
    return err
}
fmt.Println(window.Contains(time.Now()))
```

Each helper copies the time value when non-nil and returns an unbounded endpoint for nil. Like text `(nil,nil)`, two Go nil pointers passed through these helpers produce the full time domain in v2; when migrating code that relied on v1's two-nil empty interval, use `Empty[time.Time]()` explicitly in that case.

## Migration from v1

1. Keep old imports pinned to an existing v1 revision until each caller is migrated. The v1 implementation is preserved in Git history, but this branch only builds and tests v2.
2. Import `github.com/dingyou/interval/v2` for new callers. During local development use a consumer module's `replace github.com/dingyou/interval/v2 => /path/to/interval/v2` until a release is published.
3. Replace `BaseInterval` numeric constructors with `Ordered(Inclusive(...), Exclusive(...))` (choose endpoint kinds explicitly).
4. Replace `CompareTo` constraints with `New(..., compare)`; existing types can adapt their method to a comparison function.
5. Replace time wrappers with `Time`, `ParseTime`, and `FormatTime`. Ordinary `String()` now has a fixed signature; custom formatting is separate.
6. Translate nullable time endpoints explicitly. In v1, two nil endpoints meant empty; preserve that intent with `Empty[time.Time]()`. In v2 two `Unbounded` endpoints mean the entire domain. A single nil translates to `Unbounded`, and a non-nil value to `Inclusive(*p)` or `Exclusive(*p)`.
7. Handle constructor/parse errors. Do not assume reversed endpoints are repaired or that malformed bracket directions are accepted.

There are no intersection/union or comparator-identity rules in this release. A reusable Domain abstraction can be added if those operations justify it, without duplicating the core today.

## Development

```sh
# From repository root:
cd v2
go test -race -cover ./...
go vet ./...
go test -run '^$' -fuzz '^FuzzParse$' -fuzztime=10s
go test -run '^$' -fuzz '^FuzzMembership$' -fuzztime=10s
go test -run '^$' -bench . -benchmem
```

The repository root is not a Go module; run Go commands from `v2/`. Before publication, use a repository tag like `v2.0.0` (not `v2/v2.0.0`) and verify module resolution from a clean consumer.
