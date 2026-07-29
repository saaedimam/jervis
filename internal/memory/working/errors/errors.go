package errors

import "fmt"

var (
	// ErrInvalidCapacity indicates the provided capacity is not positive.
	ErrInvalidCapacity = fmt.Errorf("invalid capacity: must be greater than zero")
)
