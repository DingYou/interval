package interval

import (
	"bytes"
	"fmt"
)

type BaseInterval[T baseSortable] struct {
	left           T
	right          T
	openClosedType OpenClosedType
}

func NewBaseInterval[T baseSortable](left, right T, openCloseType ...OpenClosedType) *BaseInterval[T] {
	t := Default
	if len(openCloseType) > 0 {
		t = openCloseType[0]
	}
	return &BaseInterval[T]{
		left:           left,
		right:          right,
		openClosedType: t,
	}
}

// Left returns the left value of this interval
func (bi *BaseInterval[T]) Left() T {
	return bi.left
}

// Right returns the right value of this interval
func (bi *BaseInterval[T]) Right() T {
	return bi.right
}

// OpenClosedType returns the OpenClosedType of this interval
func (bi *BaseInterval[T]) OpenClosedType() OpenClosedType {
	return bi.openClosedType
}

// LeftClosed returns true if this interval is a left-closed interval
func (bi *BaseInterval[T]) LeftClosed() bool {
	return bi.openClosedType&ClosedOpen == ClosedOpen
}

// RightClosed returns true if this interval is a right-closed interval
func (bi *BaseInterval[T]) RightClosed() bool {
	return bi.openClosedType&OpenClosed == OpenClosed
}

// Contains returns true if the given element is in this interval
func (bi *BaseInterval[T]) Contains(e T) bool {
	if e > bi.left && e < bi.right {
		return true
	}
	if e == bi.right && bi.openClosedType&OpenClosed == OpenClosed {
		return true
	}
	if e == bi.left && bi.openClosedType&ClosedOpen == ClosedOpen {
		return true
	}
	return false
}

// String returns a readable string of this interval
func (bi *BaseInterval[T]) String() string {
	bs := bytes.Buffer{}
	if bi.LeftClosed() {
		bs.WriteString(LeftClosed)
	} else {
		bs.WriteString(LeftOpen)
	}
	bs.WriteString(fmt.Sprint(bi.left))
	bs.WriteString(Spacer)
	bs.WriteString(fmt.Sprint(bi.right))
	if bi.RightClosed() {
		bs.WriteString(RightClosed)
	} else {
		bs.WriteString(RightOpen)
	}
	return bs.String()
}
