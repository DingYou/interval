package interval

import "time"

const (
	Open       OpenClosedType = 0
	OpenClosed OpenClosedType = 1
	ClosedOpen OpenClosedType = OpenClosed << 1
	Closed     OpenClosedType = ClosedOpen | OpenClosed

	Default = ClosedOpen
)

const defaultTimeLayout = time.RFC3339

const (
	LeftClosed  = "["
	LeftOpen    = "("
	RightClosed = "]"
	RightOpen   = ")"

	Spacer = ","
	Space  = " "
)

var (
	OpenFlags   = map[string]struct{}{LeftOpen: {}, RightOpen: {}}
	ClosedFlags = map[string]struct{}{LeftClosed: {}, RightClosed: {}}
)
