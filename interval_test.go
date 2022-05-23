package interval

import "testing"

func TestNewInterval(t *testing.T) {
	i1 := NewBaseInterval[OpenClosedType](1, 2)
	t.Log(i1)
}
