package cadence

import "errors"

var (
	ErrInvalidPattern = errors.New("cadence: invalid pattern")
	ErrNotImplemented = errors.New("cadence: not implemented")
	ErrEOF                       = errors.New("cadence: eof")
)
