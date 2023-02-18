package interval

import (
	"bytes"
	"fmt"
	"strconv"
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

// ParseStrInterval parse str to interval
func ParseStrInterval(intervalStr string) (i *BaseInterval[string], err error) {
	var lf, lv, rv, rf string
	if lf, lv, rv, rf, err = blowUp(intervalStr); err != nil {
		return
	}
	var openClosedType OpenClosedType
	if openClosedType, err = getOpenClosedType(lf, rf); err != nil {
		return nil, err
	}
	return NewBaseInterval[string](lv, rv, openClosedType), nil
}

// ParseIntInterval parse str to interval
func ParseIntInterval(intervalStr string) (i *BaseInterval[int64], err error) {
	var lf, lv, rv, rf string
	if lf, lv, rv, rf, err = blowUp(intervalStr); err != nil {
		return
	}
	var openClosedType OpenClosedType
	if openClosedType, err = getOpenClosedType(lf, rf); err != nil {
		return nil, err
	}
	var lfv, rfv int64
	if lfv, err = strconv.ParseInt(lv, 10, 64); err != nil {
		return
	}
	if rfv, err = strconv.ParseInt(rv, 10, 64); err != nil {
		return
	}
	return NewBaseInterval[int64](lfv, rfv, openClosedType), nil
}

// ParseFloatInterval parse str to interval
func ParseFloatInterval(intervalStr string) (i *BaseInterval[float64], err error) {
	var lf, lv, rv, rf string
	if lf, lv, rv, rf, err = blowUp(intervalStr); err != nil {
		return
	}
	var openClosedType OpenClosedType
	if openClosedType, err = getOpenClosedType(lf, rf); err != nil {
		return nil, err
	}
	var lfv, rfv float64
	if lfv, err = strconv.ParseFloat(lv, 64); err != nil {
		return
	}
	if rfv, err = strconv.ParseFloat(rv, 64); err != nil {
		return
	}
	return NewBaseInterval[float64](lfv, rfv, openClosedType), nil
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
	return (e > bi.left && e < bi.right) ||
		(e == bi.right && bi.openClosedType&OpenClosed == OpenClosed) ||
		(e == bi.left && bi.openClosedType&ClosedOpen == ClosedOpen)
}

// String returns a readable string of this interval
func (bi *BaseInterval[T]) String() string {
	bs := &bytes.Buffer{}
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
