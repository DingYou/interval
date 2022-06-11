package interval

import "errors"

var (
	ParseTooShortErr  = errors.New("parse interval string err: str too short")
	OpenClosedFlagErr = errors.New("parse interval string err: open closed flag err")
	ValueStrErr       = errors.New("parse interval string err: value err")
)
