package laeufer

import "errors"

var (
	ErrTimeout = errors.New("timeout")

	ErrWrongDirection = errors.New("wrong direction")
)
