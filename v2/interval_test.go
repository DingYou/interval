package interval

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func mustInt(t *testing.T, a, b Bound[int]) Interval[int] {
	t.Helper()
	i, err := Ordered(a, b)
	if err != nil {
		t.Fatal(err)
	}
	return i
}
func TestMembership(t *testing.T) {
	for _, lc := range []bool{false, true} {
		for _, rc := range []bool{false, true} {
			l, r := Exclusive(1), Exclusive(3)
			if lc {
				l = Inclusive(1)
			}
			if rc {
				r = Inclusive(3)
			}
			i := mustInt(t, l, r)
			for _, v := range []int{0, 1, 2, 3, 4} {
				want := v == 2 || (v == 1 && lc) || (v == 3 && rc)
				if i.Contains(v) != want {
					t.Errorf("%s contains %d", i, v)
				}
			}
			equal := mustInt(t, l, func() Bound[int] {
				if rc {
					return Inclusive(1)
				}
				return Exclusive(1)
			}())
			if equal.Contains(1) != (lc && rc) || equal.IsCanonicalEmpty() == (lc && rc) {
				t.Errorf("equal endpoints %v %v", lc, rc)
			}
		}
	}
}
func TestEmptyAndUnbounded(t *testing.T) {
	var zero Interval[int]
	if zero.Contains(0) || !zero.IsCanonicalEmpty() || zero.String() != "empty" {
		t.Fatal(zero)
	}
	if _, _, ok := zero.Bounds(); ok {
		t.Fatal("empty bounds reported valid")
	}
	all := mustInt(t, Unbounded[int](), Unbounded[int]())
	for _, v := range []int{-100, 0, 100} {
		if !all.Contains(v) {
			t.Fatal(v)
		}
	}
	if all.String() != "(,)" {
		t.Fatal(all)
	}
	left := mustInt(t, Unbounded[int](), Inclusive(2))
	right := mustInt(t, Exclusive(2), Unbounded[int]())
	if !left.Contains(-100) || left.Contains(3) || right.Contains(2) || !right.Contains(100) {
		t.Fatal("one-sided membership")
	}
	a, b, ok := left.Bounds()
	if !ok || !a.IsUnbounded() || !b.IsInclusive() {
		t.Fatal("bounds")
	}
	if v, finite := b.Value(); !finite || v != 2 {
		t.Fatal("value")
	}
	discrete := mustInt(t, Exclusive(1), Exclusive(2))
	if discrete.IsCanonicalEmpty() || discrete.Contains(1) || discrete.Contains(2) {
		t.Fatal("discrete semantics")
	}
}
func TestErrorsAndOrdering(t *testing.T) {
	if _, err := Ordered(Inclusive(2), Inclusive(1)); !errors.Is(err, ErrReversed) {
		t.Fatal(err)
	}
	if _, err := New(Unbounded[int](), Unbounded[int](), nil); !errors.Is(err, ErrNilCompare) {
		t.Fatal(err)
	}
	type record struct {
		Key     int
		Payload []byte
	}
	compare := func(a, b record) int { return cmp.Compare(a.Key, b.Key) }
	i, err := New(Inclusive(record{Key: 1}), Exclusive(record{Key: 3}), compare)
	if err != nil {
		t.Fatal(err)
	}
	if !i.Contains(record{Key: 2, Payload: []byte("non-comparable")}) || i.Contains(record{Key: 3}) {
		t.Fatal("custom ordering")
	}
	reverse, err := New(Inclusive(3), Inclusive(1), func(a, b int) int { return cmp.Compare(b, a) })
	if err != nil || !reverse.Contains(2) {
		t.Fatal(reverse, err)
	}
	nan, err := Ordered(Inclusive(math.NaN()), Inclusive(math.NaN()))
	if err != nil || !nan.Contains(math.NaN()) || nan.Contains(0) {
		t.Fatal(nan, err)
	}
	zeros, err := Ordered(Inclusive(math.Copysign(0, -1)), Inclusive(0.0))
	if err != nil || !zeros.Contains(0) {
		t.Fatal(zeros, err)
	}
}
func TestParse(t *testing.T) {
	valid := map[string]string{" [1, 3) ": "[1,3)", "(,)": "(,)", "(,2]": "(,2]", "[2,)": "[2,)", "[2,2)": "empty", "empty": "empty"}
	for input, want := range valid {
		i, err := Parse(input, strconv.Atoi, cmp.Compare[int])
		if err != nil || i.String() != want {
			t.Errorf("%q: %v %v", input, i, err)
		}
	}
	for _, input := range []string{"", "[]", "]1,2[", "[1,2,3]", "[,2]", "[1,]", "[\"1,2]", "[1,(2)]", "[\"1\"junk,2]", "[1,2\\]"} {
		if _, err := Parse(input, strconv.Atoi, cmp.Compare[int]); !errors.Is(err, ErrSyntax) {
			t.Errorf("%q: %v", input, err)
		}
	}
	if _, err := Parse("[3,1]", strconv.Atoi, cmp.Compare[int]); !errors.Is(err, ErrReversed) {
		t.Fatal(err)
	}
	if _, err := Parse("empty", strconv.Atoi, nil); !errors.Is(err, ErrNilCompare) {
		t.Fatal(err)
	}
	if _, err := Parse[int]("empty", nil, cmp.Compare[int]); !errors.Is(err, ErrNilDecode) {
		t.Fatal(err)
	}
	sentinel := errors.New("decode failed")
	_, err := Parse("[x,y]", func(string) (int, error) { return 0, sentinel }, cmp.Compare[int])
	if !errors.Is(err, sentinel) || !strings.Contains(err.Error(), "lower") {
		t.Fatal(err)
	}
	_, err = Parse("[1,x]", strconv.Atoi, cmp.Compare[int])
	var numberErr *strconv.NumError
	if !errors.As(err, &numberErr) || !strings.Contains(err.Error(), "upper") {
		t.Fatal(err)
	}
}
func TestTextRoundTrip(t *testing.T) {
	decode := func(s string) (string, error) { return s, nil }
	for _, s := range []string{"", "plain", "a,b", "[x]", "(x)", " a ", "a\nb", "a\"b", "a\\b", "你好"} {
		i, err := Ordered(Inclusive(s), Inclusive(s))
		if err != nil {
			t.Fatal(err)
		}
		j, err := Parse(i.String(), decode, cmp.Compare[string])
		if err != nil || !j.Contains(s) {
			t.Fatalf("%q => %s: %v", s, i, err)
		}
		if i.Format(nil) != i.String() {
			t.Fatal("nil formatter")
		}
	}
}
func TestTime(t *testing.T) {
	start, err := time.Parse(time.RFC3339Nano, "2026-11-13T12:00:00.123456789+08:00")
	if err != nil {
		t.Fatal(err)
	}
	end := start.Add(time.Second)
	i, err := Time(Inclusive(start), Exclusive(end))
	if err != nil {
		t.Fatal(err)
	}
	var _ fmt.Stringer = i
	parsed, err := ParseTime(FormatTime(i))
	if err != nil {
		t.Fatal(err)
	}
	lower, upper, ok := parsed.Bounds()
	a, _ := lower.Value()
	b, _ := upper.Value()
	if !ok || !a.Equal(start) || !b.Equal(end) || !parsed.Contains(start) || parsed.Contains(end) {
		t.Fatal(FormatTime(parsed))
	}
	if _, err := ParseTime("[bad,bad]"); err == nil {
		t.Fatal("bad time accepted")
	}
}
func FuzzParse(f *testing.F) {
	for _, s := range []string{"[1,2)", "empty", "(,)", "[\"1\",2]", "]1,2["} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		i, err := Parse(s, strconv.Atoi, cmp.Compare[int])
		if err != nil {
			return
		}
		j, err := Parse(i.String(), strconv.Atoi, cmp.Compare[int])
		if err != nil {
			t.Fatal(err)
		}
		if i.String() != j.String() {
			t.Fatal("unstable round trip")
		}
		for _, v := range []int{-1, 0, 1, 2, 100} {
			if i.Contains(v) != j.Contains(v) {
				t.Fatal("membership changed")
			}
		}
	})
}
func FuzzMembership(f *testing.F) {
	f.Add(1, 3, 2, true, false)
	f.Fuzz(func(t *testing.T, a, b, v int, lc, rc bool) {
		l, r := Exclusive(a), Exclusive(b)
		if lc {
			l = Inclusive(a)
		}
		if rc {
			r = Inclusive(b)
		}
		i, err := Ordered(l, r)
		if a > b {
			if !errors.Is(err, ErrReversed) {
				t.Fatal(err)
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		want := (v > a || (v == a && lc)) && (v < b || (v == b && rc))
		if i.Contains(v) != want {
			t.Fatal("membership mismatch")
		}
	})
}
func BenchmarkContains(b *testing.B) {
	i, err := Ordered(Inclusive(1), Exclusive(10))
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		_ = i.Contains(n % 12)
	}
}
func ExampleOrdered() {
	i, err := Ordered(Inclusive(1), Exclusive(5))
	if err != nil {
		panic(err)
	}
	fmt.Println(i, i.Contains(1), i.Contains(5))
	// Output: [1,5) true false
}

func FuzzStringRoundTrip(f *testing.F) {
	for _, s := range []string{"", "a,b", "a\"b", "\x00\xff", " \u2003", "你好", "a\\b"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		i, err := Ordered(Inclusive(s), Inclusive(s))
		if err != nil {
			t.Fatal(err)
		}
		j, err := Parse(i.String(), func(s string) (string, error) { return s, nil }, cmp.Compare[string])
		if err != nil {
			t.Fatal(err)
		}
		if !j.Contains(s) || j.String() != i.String() {
			t.Fatal("string round trip changed endpoint")
		}
	})
}
