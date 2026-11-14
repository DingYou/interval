// Package interval implements immutable interval bounds with caller-selected ordering.
// A zero Interval is empty; construct an all-values interval with two Unbounded bounds.
package interval

import (
	"cmp"
	"errors"
	"fmt"
)

var (
	// ErrNilCompare indicates that no ordering function was supplied.
	ErrNilCompare = errors.New("interval: nil comparison function")
	// ErrReversed indicates a lower endpoint greater than the upper endpoint.
	ErrReversed = errors.New("interval: reversed endpoints")
)

type boundKind uint8

const (
	unbounded boundKind = iota
	inclusive
	exclusive
)

// Bound is a finite inclusive/exclusive endpoint or an unbounded endpoint.
// Its zero value is unbounded. Values are copied, not deep-cloned.
type Bound[T any] struct {
	value T
	kind  boundKind
}

// Inclusive constructs a closed finite endpoint.
func Inclusive[T any](value T) Bound[T] { return Bound[T]{value, inclusive} }

// Exclusive constructs an open finite endpoint.
func Exclusive[T any](value T) Bound[T] { return Bound[T]{value, exclusive} }

// Unbounded constructs an endpoint without a limit.
func Unbounded[T any]() Bound[T] { return Bound[T]{} }

// Value returns the finite endpoint and true, or the zero value and false.
func (b Bound[T]) Value() (T, bool) { return b.value, b.kind != unbounded }

// IsInclusive reports whether this is a closed finite endpoint.
func (b Bound[T]) IsInclusive() bool { return b.kind == inclusive }

// IsUnbounded reports whether this endpoint has no limit.
func (b Bound[T]) IsUnbounded() bool { return b.kind == unbounded }

// Interval is a pair of bounds interpreted by a comparison function.
// The function must be a deterministic total ordering: negative for a < b,
// zero for equivalent values, positive for a > b. The library cannot verify this.
// Values and comparison closures must not be mutated while an interval is in use.
// Safe concurrent reads require the comparator itself to be concurrency-safe.
type Interval[T any] struct {
	lower, upper Bound[T]
	compare      func(T, T) int
}

// Empty constructs the canonical empty interval.
func Empty[T any]() Interval[T] { return Interval[T]{} }

// New constructs an interval under compare. Reversed finite endpoints return
// ErrReversed. Equal endpoints with either end excluded normalize to Empty.
// A nil comparator is rejected even if the bounds are both unbounded.
func New[T any](lower, upper Bound[T], compare func(T, T) int) (Interval[T], error) {
	if compare == nil {
		return Empty[T](), ErrNilCompare
	}
	if !lower.IsUnbounded() && !upper.IsUnbounded() {
		c := compare(lower.value, upper.value)
		if c > 0 {
			return Empty[T](), ErrReversed
		}
		if c == 0 && (!lower.IsInclusive() || !upper.IsInclusive()) {
			return Empty[T](), nil
		}
	}
	return Interval[T]{lower, upper, compare}, nil
}

// Ordered constructs an interval using cmp.Compare, including its documented
// total order for floating-point NaN and signed zero.
func Ordered[T cmp.Ordered](lower, upper Bound[T]) (Interval[T], error) {
	return New(lower, upper, cmp.Compare[T])
}

// Bounds returns the endpoints and true for a noncanonical-empty interval.
// For canonical empty intervals it returns zero bounds and false; callers must
// check that boolean, since two unbounded bounds otherwise mean all values.
func (i Interval[T]) Bounds() (lower, upper Bound[T], ok bool) {
	return i.lower, i.upper, i.compare != nil
}

// IsCanonicalEmpty reports whether this interval is the normalized empty value.
// False does not promise a member exists: for example (1,2) over integers has none.
func (i Interval[T]) IsCanonicalEmpty() bool { return i.compare == nil }

// Contains reports membership, requiring both endpoint conditions.
func (i Interval[T]) Contains(value T) bool {
	if i.compare == nil {
		return false
	}
	if !i.lower.IsUnbounded() {
		c := i.compare(value, i.lower.value)
		if c < 0 || (c == 0 && !i.lower.IsInclusive()) {
			return false
		}
	}
	if !i.upper.IsUnbounded() {
		c := i.compare(value, i.upper.value)
		if c > 0 || (c == 0 && !i.upper.IsInclusive()) {
			return false
		}
	}
	return true
}

// String formats endpoints with fmt.Sprint. It is for display; use Format with
// a matching decoder for lossless serialization of arbitrary endpoint types.
func (i Interval[T]) String() string { return i.Format(func(v T) string { return fmt.Sprint(v) }) }
