package interval

import "testing"

func TestNewBaseInterval(t *testing.T) {
	i1 := NewBaseInterval[int](1, 2)
	t.Log(i1)
	t.Log(i1.Contains(2))
	t.Log(i1)
}
