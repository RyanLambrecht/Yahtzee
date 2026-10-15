package cli

import "errors"

var (
	ErrRerollTooLong    = errors.New("reroll input too long")
	ErrRerollAtoi       = errors.New("input isnt a integer")
	ErrSkipReroll       = errors.New("skipped reroll")
	ErrInputOutOfBounds = errors.New("input out of bounds")
)
