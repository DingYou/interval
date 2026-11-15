package interval

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrSyntax indicates invalid interval text.
	ErrSyntax = errors.New("interval: invalid syntax")
	// ErrNilDecode indicates a missing endpoint decoder.
	ErrNilDecode = errors.New("interval: nil endpoint decoder")
)

// Format writes empty as "empty", and unbounded endpoints as omitted text.
// Finite text that is empty or contains delimiters, quotes, backslashes or outer
// whitespace is Go-quoted, so Parse can distinguish it from an unbounded end.
// A nil formatter uses fmt.Sprint.
func (i Interval[T]) Format(format func(T) string) string {
	if i.IsCanonicalEmpty() {
		return "empty"
	}
	if format == nil {
		format = func(v T) string { return fmt.Sprint(v) }
	}
	endpoint := func(b Bound[T]) string {
		if b.IsUnbounded() {
			return ""
		}
		s := format(b.value)
		if s == "" || strings.TrimSpace(s) != s || strings.ContainsAny(s, ",[]()\"\\\n\r\t") {
			return strconv.Quote(s)
		}
		return s
	}
	left, right := "(", ")"
	if i.lower.IsInclusive() {
		left = "["
	}
	if i.upper.IsInclusive() {
		right = "]"
	}
	return left + endpoint(i.lower) + "," + endpoint(i.upper) + right
}

// Parse reads "empty" or an interval such as [1,2), (,2], or (,).
// An omitted endpoint must use an open bracket. Quoted endpoints use Go string
// literal escaping; surrounding unquoted whitespace is ignored. Endpoint decoder
// errors are wrapped and remain accessible through errors.Is/errors.As.
func Parse[T any](text string, decode func(string) (T, error), compare func(T, T) int) (Interval[T], error) {
	return parse(text, decode, compare, nil)
}

// aliases handles optional unquoted endpoint syntax before finite decoding.
func parse[T any](text string, decode func(string) (T, error), compare func(T, T) int, aliases func(raw string, lower, closed bool) (matched bool, err error)) (Interval[T], error) {
	if compare == nil {
		return Empty[T](), ErrNilCompare
	}
	if decode == nil {
		return Empty[T](), ErrNilDecode
	}
	text = strings.TrimSpace(text)
	if text == "empty" {
		return Empty[T](), nil
	}
	if len(text) < 3 || (text[0] != '[' && text[0] != '(') || (text[len(text)-1] != ']' && text[len(text)-1] != ')') {
		return Empty[T](), ErrSyntax
	}
	body := text[1 : len(text)-1]
	split, quoted, escaped := -1, false, false
	for j, r := range body {
		if escaped {
			escaped = false
			continue
		}
		if quoted && r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			quoted = !quoted
			continue
		}
		if r == ',' && !quoted {
			if split >= 0 {
				return Empty[T](), ErrSyntax
			}
			split = j
		}
	}
	if quoted || split < 0 {
		return Empty[T](), ErrSyntax
	}
	parseBound := func(raw string, lower, closed bool) (Bound[T], error) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			if closed {
				return Bound[T]{}, ErrSyntax
			}
			return Unbounded[T](), nil
		}
		if aliases != nil && raw[0] != '"' {
			matched, err := aliases(raw, lower, closed)
			if matched || err != nil {
				return Unbounded[T](), err
			}
		}
		if raw[0] == '"' {
			s, err := strconv.Unquote(raw)
			if err != nil {
				return Bound[T]{}, fmt.Errorf("%w: %v", ErrSyntax, err)
			}
			raw = s
		} else if strings.ContainsAny(raw, "[]()\"\\") {
			return Bound[T]{}, ErrSyntax
		}
		value, err := decode(raw)
		if err != nil {
			return Bound[T]{}, err
		}
		if closed {
			return Inclusive(value), nil
		}
		return Exclusive(value), nil
	}
	lower, err := parseBound(body[:split], true, text[0] == '[')
	if err != nil {
		return Empty[T](), fmt.Errorf("interval: lower endpoint: %w", err)
	}
	upper, err := parseBound(body[split+1:], false, text[len(text)-1] == ']')
	if err != nil {
		return Empty[T](), fmt.Errorf("interval: upper endpoint: %w", err)
	}
	return New(lower, upper, compare)
}

// TimeInclusive copies a non-nil time into a closed finite endpoint.
// A nil time represents an unbounded endpoint.
func TimeInclusive(t *time.Time) Bound[time.Time] {
	if t == nil {
		return Unbounded[time.Time]()
	}
	return Inclusive(*t)
}

// TimeExclusive copies a non-nil time into an open finite endpoint.
// A nil time represents an unbounded endpoint.
func TimeExclusive(t *time.Time) Bound[time.Time] {
	if t == nil {
		return Unbounded[time.Time]()
	}
	return Exclusive(*t)
}

// Time constructs a time interval using time.Time.Compare.
func Time(lower, upper Bound[time.Time]) (Interval[time.Time], error) {
	return New(lower, upper, time.Time.Compare)
}

// ParseTime parses RFC3339Nano endpoints. Unquoted NULL and nil are
// case-insensitive aliases that ignore their brackets. Any two unbounded
// endpoints represent the full time domain. Directional infinity (-inf on the
// lower side, inf or +inf on the upper side) requires an open bracket. Serialized times have
// no monotonic clock reading.
func ParseTime(text string) (Interval[time.Time], error) {
	aliases := func(raw string, lower, closed bool) (bool, error) {
		switch strings.ToLower(raw) {
		case "null", "nil":
			return true, nil
		case "-inf":
			if lower && !closed {
				return true, nil
			}
		case "inf", "+inf":
			if !lower && !closed {
				return true, nil
			}
		default:
			return false, nil
		}
		return true, ErrSyntax
	}
	decode := func(s string) (time.Time, error) { return time.Parse(time.RFC3339Nano, s) }
	return parse(text, decode, time.Time.Compare, aliases)
}

// FormatTime formats time endpoints as RFC3339Nano, preserving fractional seconds.
// Like other textual time encodings it does not preserve monotonic clock readings
// or named Location identity; wall-clock instants and numeric offsets are retained.
func FormatTime(i Interval[time.Time]) string {
	return i.Format(func(t time.Time) string { return t.Format(time.RFC3339Nano) })
}
