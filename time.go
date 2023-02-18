package interval

import (
	"bytes"
	"time"
)

// TimeInterval time interval
type TimeInterval struct {
	left           time.Time
	right          time.Time
	openClosedType OpenClosedType
}

// NewTimeInterval return a new TimeInterval
func NewTimeInterval(left, right time.Time, openCloseType ...OpenClosedType) *TimeInterval {
	t := Default
	if len(openCloseType) > 0 {
		t = openCloseType[0]
	}
	return &TimeInterval{
		left:           left,
		right:          right,
		openClosedType: t,
	}
}

// ParseTimeInterval parse str to interval
func ParseTimeInterval(intervalStr string, layout ...string) (ti *TimeInterval, err error) {
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
	var lt, rt time.Time
	if lt, err = time.Parse(l, lv); err != nil {
		return nil, err
	}
	if rt, err = time.Parse(l, rv); err != nil {
		return nil, err
	}
	return NewTimeInterval(lt, rt, openClosedType), nil
}

// Left returns the left value of this interval
func (ti *TimeInterval) Left() time.Time {
	return ti.left
}

// Right returns the right value of this interval
func (ti *TimeInterval) Right() time.Time {
	return ti.right
}

// OpenClosedType returns the OpenClosedType of this interval
func (ti *TimeInterval) OpenClosedType() OpenClosedType {
	return ti.openClosedType
}

// LeftClosed returns true if interval is a left-closed interval
func (ti *TimeInterval) LeftClosed() bool {
	return ti.openClosedType&ClosedOpen == ClosedOpen
}

// RightClosed returns true if interval is a right-closed interval
func (ti *TimeInterval) RightClosed() bool {
	return ti.openClosedType&OpenClosed == OpenClosed
}

// Contains return ture if this interval contains the given element
func (ti *TimeInterval) Contains(e time.Time) bool {
	return (e.After(ti.left) && e.Before(ti.right)) ||
		(e.Equal(ti.left) && ti.openClosedType&ClosedOpen == ClosedOpen) ||
		(e.Equal(ti.right) && ti.openClosedType&OpenClosed == OpenClosed)
}

// String returns a readable string of this interval
func (ti *TimeInterval) String(layout ...string) string {
	l := defaultTimeLayout
	if len(layout) > 0 {
		l = layout[0]
	}
	bs := bytes.Buffer{}
	if ti.LeftClosed() {
		bs.WriteString(LeftClosed)
	} else {
		bs.WriteString(LeftOpen)
	}
	bs.WriteString(ti.left.Format(l))
	bs.WriteString(Spacer)
	bs.WriteString(ti.right.Format(l))
	if ti.RightClosed() {
		bs.WriteString(RightClosed)
	} else {
		bs.WriteString(RightOpen)
	}
	return bs.String()
}
