package interval

import "testing"

func TestSortComparable(t *testing.T) {
	s1, s2 := &testCompareStruct{"Ha", 1}, &testCompareStruct{"So", 2}
	t.Log(Compare(s1, s2))

	s3 := (*testCompareStruct)(nil)
	t.Log(Compare(nil, s3))
	t.Log(s3.CompareTo(s1))
}

type testCompareStruct struct {
	Name  string
	Score int
}

func (s *testCompareStruct) CompareTo(other *testCompareStruct) int {
	if s == other {
		return 0
	}
	if s == nil {
		return -1
	}
	if other == nil {
		return 1
	}
	return s.Score - other.Score
}
