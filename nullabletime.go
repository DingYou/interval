package interval

import (
	"bytes"
	"strings"
	"time"
)

const (
	NullFlag = "NULL"
)

// NullableTimeInterval nullable time interval
type NullableTimeInterval struct {
	left           *time.Time
	right          *time.Time
	openClosedType OpenClosedType
}

// NewNullableTimeInterval return a new NullableTimeInterval
func NewNullableTimeInterval(left, right *time.Time, openCloseType ...OpenClosedType) *NullableTimeInterval {
	t := Default
	if len(openCloseType) > 0 {
		t = openCloseType[0]
	}
	return &NullableTimeInterval{
		left:           left,
		right:          right,
		openClosedType: t,
	}
}

// ParseNullableTimeInterval parse str to interval
func ParseNullableTimeInterval(intervalStr string, layout ...string) (ti *NullableTimeInterval, err error) {
	var lf, lv, rv, rf string
	if lf, lv, rv, rf, err = blowUp(intervalStr); err != nil {
		return
	}
	var openClosedType OpenClosedType
	if openClosedType, err = getOpenClosedType(lf, rf); err != nil {
		return nil, err
	}
	l := defaultTimeLayout
	if len(layout) > 0 {
		l = layout[0]
	}
	var lt, rt *time.Time
	if lt, err = parseNullableTimeStr(l, lv); err != nil {
		return nil, err
	}
	if rt, err = parseNullableTimeStr(l, rv); err != nil {
		return nil, err
	}
	return NewNullableTimeInterval(lt, rt, openClosedType), nil
}

// Left returns the left value of this interval
func (ti *NullableTimeInterval) Left() *time.Time {
	return ti.left
}

// Right returns the right value of this interval
func (ti *NullableTimeInterval) Right() *time.Time {
	return ti.right
}

// OpenClosedType returns the OpenClosedType of this interval
func (ti *NullableTimeInterval) OpenClosedType() OpenClosedType {
	return ti.openClosedType
}

// LeftClosed returns true if interval is a left-closed interval
func (ti *NullableTimeInterval) LeftClosed() bool {
	return ti.openClosedType&ClosedOpen == ClosedOpen
}

// RightClosed returns true if interval is a right-closed interval
func (ti *NullableTimeInterval) RightClosed() bool {
	return ti.openClosedType&OpenClosed == OpenClosed
}

// Contains return ture if this interval contains the given element
// null element is out of any interval
func (ti *NullableTimeInterval) Contains(e *time.Time) bool {
	if e == nil {
		return false
	}
	switch {
	case ti.left != nil && ti.right != nil:
		return (e.After(*ti.left) && e.Before(*ti.right)) ||
			(e.Equal(*ti.left) && ti.openClosedType&ClosedOpen == ClosedOpen) ||
			(e.Equal(*ti.right) && ti.openClosedType&OpenClosed == OpenClosed)
	case ti.left == nil && ti.right == nil:
		return false
	case ti.left == nil:
		return (e.Before(*ti.right)) || (e.Equal(*ti.right) && ti.openClosedType&OpenClosed == OpenClosed)
	case ti.right == nil:
		return (e.After(*ti.left)) || (e.Equal(*ti.left) && ti.openClosedType&ClosedOpen == ClosedOpen)
	}
	return false
}

// String returns a readable string of this interval
func (ti *NullableTimeInterval) String(layout ...string) string {
	l := defaultTimeLayout
	if len(layout) > 0 {
		l = layout[0]
	}
	bs := &bytes.Buffer{}
	if ti.LeftClosed() {
		bs.WriteString(LeftClosed)
	} else {
		bs.WriteString(LeftOpen)
	}
	if ti.left != nil {
		bs.WriteString(ti.left.Format(l))
	} else {
		bs.WriteString(NullFlag)
	}
	bs.WriteString(Spacer)
	if ti.right != nil {
		bs.WriteString(ti.right.Format(l))
	} else {
		bs.WriteString(NullFlag)
	}
	if ti.RightClosed() {
		bs.WriteString(RightClosed)
	} else {
		bs.WriteString(RightOpen)
	}
	return bs.String()
}

func parseNullableTimeStr(layout, value string) (*time.Time, error) {
	if strings.ToUpper(value) == NullFlag {
		return nil, nil
	}
	if t, err := time.Parse(layout, value); err != nil {
		return nil, err
	} else {
		return &t, nil
	}
}
