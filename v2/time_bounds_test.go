package interval

import (
	"testing"
	"time"
)

func TestTimePointerBounds(t *testing.T) {
	var start time.Time
	end := start.Add(time.Hour)
	before := start.Add(-time.Second)
	after := end.Add(time.Second)

	for _, tc := range []struct {
		name         string
		lower, upper Bound[time.Time]
		want         [4]bool
	}{
		{"finite", TimeInclusive(&start), TimeExclusive(&end), [4]bool{false, true, false, false}},
		{"unbounded lower", TimeInclusive(nil), TimeExclusive(&end), [4]bool{true, true, false, false}},
		{"unbounded upper", TimeInclusive(&start), TimeExclusive(nil), [4]bool{false, true, true, true}},
		{"full domain", TimeInclusive(nil), TimeExclusive(nil), [4]bool{true, true, true, true}},
		{"open lower and closed upper", TimeExclusive(&start), TimeInclusive(&end), [4]bool{false, false, true, false}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			window, err := Time(tc.lower, tc.upper)
			if err != nil {
				t.Fatal(err)
			}
			for idx, value := range []time.Time{before, start, end, after} {
				if got := window.Contains(value); got != tc.want[idx] {
					t.Errorf("Contains(%v) = %v, want %v", value, got, tc.want[idx])
				}
			}
		})
	}
}

func TestTimePointerBoundsCopyValues(t *testing.T) {
	start := time.Unix(0, 0)
	end := start.Add(time.Hour)
	originalStart, originalEnd := start, end
	window, err := Time(TimeInclusive(&start), TimeExclusive(&end))
	if err != nil {
		t.Fatal(err)
	}
	start, end = start.Add(time.Minute), end.Add(time.Minute)
	if !window.Contains(originalStart) || window.Contains(originalEnd) {
		t.Fatal("pointer mutation changed the constructed interval")
	}
}
