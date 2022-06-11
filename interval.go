package interval

import (
	"bytes"
	"fmt"
	"strings"
)

type OpenClosedType uint8

type IInterval[T any] interface {
	// Left returns the left value of this interval
	Left() T
	// Right returns the right value of this interval
	Right() T
	// OpenClosedType returns the OpenClosedType of interval
	OpenClosedType() OpenClosedType
	// LeftClosed returns true if interval is a left-closed interval
	LeftClosed() bool
	// RightClosed returns true if interval is a right-closed interval
	RightClosed() bool
	// Contains returns true if given element is in interval
	Contains(e T) bool
	// String returns a readable string
	String() string
}

type Interval[T SortComparable[T]] struct {
	left           T
	right          T
	openClosedType OpenClosedType
}

// Left returns the left value of this interval
func (i *Interval[T]) Left() T {
	return i.left
}

func (i *Interval[T]) Right() T {
	return i.right
}

// OpenClosedType returns the OpenClosedType of this interval
func (i *Interval[T]) OpenClosedType() OpenClosedType {
	return i.openClosedType
}

// LeftClosed returns true if this interval is a left-closed interval
func (i *Interval[T]) LeftClosed() bool {
	return i.openClosedType&ClosedOpen == ClosedOpen
}

// RightClosed returns true if this interval is a right-closed interval
func (i *Interval[T]) RightClosed() bool {
	return i.openClosedType&OpenClosed == OpenClosed
}

func (i *Interval[T]) Contains(e T) bool {
	if Compare[T](i.left, e) < 0 && Compare[T](i.right, e) > 0 {
		return true
	}
	if Compare[T](i.left, e) == 0 && i.openClosedType&ClosedOpen == ClosedOpen {
		return true
	}
	if Compare[T](i.right, e) == 0 && i.openClosedType&OpenClosed == OpenClosed {
		return true
	}
	return false
}

// String returns a readable string of this interval
func (i Interval[T]) String() string {
	bs := bytes.Buffer{}
	if i.LeftClosed() {
		bs.WriteString(LeftClosed)
	} else {
		bs.WriteString(LeftOpen)
	}
	bs.WriteString(fmt.Sprintf("%v", i.left))
	bs.WriteString(Spacer)
	bs.WriteString(fmt.Sprintf("%v", i.right))
	if i.RightClosed() {
		bs.WriteString(RightClosed)
	} else {
		bs.WriteString(RightOpen)
	}
	return bs.String()
}

func blowUp(str string) (string, string, string, string, error) {
	strLen := len(str)
	if strLen < 5 {
		return "", "", "", "", ParseTooShortErr
	}
	leftFlag, rightFlag := string(str[0]), string(str[strLen-1])
	str = str[1 : strLen-1]
	values := strings.Split(str, Spacer)
	if len(values) != 2 {
		return "", "", "", "", ValueStrErr
	}
	return leftFlag, strings.Trim(values[0], Space), strings.Trim(values[1], Space), rightFlag, nil
}

func isClosedFlag(in string) (bool, error) {
	if _, exist := OpenFlags[in]; exist {
		return false, nil
	}
	if _, exist := ClosedFlags[in]; exist {
		return true, nil
	}
	return false, OpenClosedFlagErr
}

func getOpenClosedType(lf, rf string) (OpenClosedType, error) {
	openClosedType := Open
	if closed, err := isClosedFlag(lf); err != nil {
		return openClosedType, err
	} else if closed {
		openClosedType |= ClosedOpen
	}
	if closed, err := isClosedFlag(rf); err != nil {
		return openClosedType, err
	} else if closed {
		openClosedType |= OpenClosed
	}
	return openClosedType, nil
}
