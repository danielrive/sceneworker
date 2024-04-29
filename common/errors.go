package common

import "errors"

// Common errors
var (
	ErrorTimeout = errors.New("Timeout")
	ErrorRefused = errors.New("Refused")
)
