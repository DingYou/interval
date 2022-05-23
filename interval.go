package interval

type OpenClosedType uint8

const (
	Open                      = 0
	OpenClosed                = 1
	ClosedOpen OpenClosedType = OpenClosed << 1
	Closed                    = ClosedOpen | OpenClosed

	Default = ClosedOpen
)

type Interval[T any] interface {
	// Left returns the left value of interval
	Left() T
	// Right returns the right value of interval
	Right() T
	// OpenCloseType returns the OpenClosedType of interval
	OpenCloseType() OpenClosedType
	// LeftClosed returns true if interval is a left-closed interval
	LeftClosed() bool
	// RightClosed returns true if interval is a right-closed interval
	RightClosed() bool
	// InInterval returns true if given element is in interval
	InInterval(e T) bool
	// String returns a readable string
	String() string
}

type BaseInterval[T baseSortable] struct {
	left          T
	right         T
	openCloseType OpenClosedType
}

func NewBaseInterval[T baseSortable](left, right T) *BaseInterval[T] {
	return &BaseInterval[T]{
		left:          left,
		right:         right,
		openCloseType: Default,
	}
}

func (bi *BaseInterval[T]) Left() T {
	return bi.left
}

func (bi *BaseInterval[T]) Right() T {
	return bi.right
}

func (bi *BaseInterval[T]) OpenClosedType() OpenClosedType {
	return bi.openCloseType
}

func (bi *BaseInterval[T]) LeftClosed() bool {
	return bi.openCloseType&ClosedOpen == ClosedOpen
}

func (bi *BaseInterval[T]) RightClosed() bool {
	return bi.openCloseType&OpenClosed == OpenClosed
}

func (bi *BaseInterval[T]) InInterval(e T) bool {
	return false
}
