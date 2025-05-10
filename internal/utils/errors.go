package utils

import "errors"

// Errors maps a string key to a list of values.
type Errors map[string][]string

// Different error types
var (
	ErrServerNotFound = errors.New("server not found")
)
